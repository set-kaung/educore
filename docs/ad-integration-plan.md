# Plan: University AD (Microsoft Entra ID) Integration via OIDC

Requirement (from `docs/requirements.md`):

> integrating the University's Microsoft Active Directory (AD) for user authentication (using MSAL, OAuth2, or OIDC).

Chosen approach: **OIDC Authorization Code Flow** (browser redirect to Microsoft login). ROPC
(direct username/password) was rejected: deprecated by Microsoft, commonly disabled tenant-wide,
and incompatible with MFA/Conditional Access.

## 0. Blocker: App Registration (permission required)

Microsoft Entra ID rejects unregistered clients — we cannot call its endpoints until someone with
admin rights in the university tenant registers EduCore. Request from the instructor / university IT:

- **Tenant ID** of the university Entra directory
- **App Registration** for EduCore (or permission for us to self-register)
- **Client ID** + **Client Secret** (confidential/web client)
- **Redirect URIs** registered on the app (exact string match, HTTPS required except localhost):
  - `http://localhost:8080/auth/callback` (local dev)
  - `https://<your-domain>/<project-path>/auth/callback` (production)
- API permissions: `openid`, `profile`, `email`, `User.Read` (admin-consented if the tenant requires)

**Development can start before this arrives**: create a free personal Azure account → own Entra
tenant → own app registration, and build against that. Everything is config-driven, so switching to
the university tenant later means only changing env vars.

## 1. Configuration

New env vars (`cmd/server` `Config`, read in `main.go`):

| Var                | Purpose                                        |
| ------------------ | ---------------------------------------------- |
| `AD_CLIENT_ID`     | App registration client ID                     |
| `AD_TENANT_ID`     | University (or dev) tenant ID                  |
| `AD_CLIENT_SECRET` | Client secret — from Azure Key Vault in prod   |
| `AD_REDIRECT_URI`  | Must exactly match a registered redirect URI   |

All four `AD_*` vars are required; the server fails fast at startup if any is missing.

## 2. Authentication flows

### Returning user (account already linked)

```
Browser  → GET /auth/login                    (generate state+nonce cookies, 302 to Microsoft)
Browser  → login.microsoftonline.com/...      (user signs in, MFA included)
Browser  → GET /auth/callback?code=...&state=...
Server:  verify state → exchange code for tokens (server-to-server, uses client secret)
         → verify ID token (signature via JWKS, iss, aud, expiry, nonce)
         → lookup Student/Professor WHERE ad_object_id = token "oid" claim
         → FOUND: issue normal session JWT, set edu_token cookie, 302 → /
```

### First-time user (no local account yet) → setup page

```
... same up to the ad_object_id lookup ...
         → NOT FOUND: issue short-lived SETUP token (JWT, scope="setup", ~15 min,
           carries oid + email + suggested name), 302 → /setup
Browser  → GET /setup                         (static setup.html)
Browser  → GET /api/setup/context             (setup token) → {email, name} for prefill
Browser  → GET /api/departments               (public) → department dropdown options
User fills form: name (editable, prefilled), role (student|professor), department
         (+ student ID field shown when role=student)
Browser  → POST /api/setup                    (setup token + SameOrigin guard)
Server:  validate → create Student or Professor row (ad_object_id = oid)
         → issue normal session JWT cookie → frontend redirects to /
```

The existing JWT / RBAC / `edu_token` cookie stack is untouched — OIDC only changes how the user
is identified before a session is issued.

## 3. Backend changes

### New dependencies

- `github.com/coreos/go-oidc/v3` — discovery, JWKS verification, ID token validation
- `golang.org/x/oauth2` — authorization-code exchange

(MSAL Go is not used: `go-oidc` + `x/oauth2` is the idiomatic Go stack for server-side web flows
and satisfies the "MSAL, OAuth2, or OIDC" requirement.)

### New code

| File | Contents |
| ---- | -------- |
| `internal/auth/oidc_provider.go` | Provider setup: discovery URL from `AD_TENANT_ID`, `oauth2.Config` (client ID/secret/redirect), ID-token verifier, claims struct (`oid`, `email`, `preferred_username`, `name`) |
| `internal/auth/oidc_handler.go` | `HandleADLogin` (state+nonce cookies, redirect), `HandleADCallback` (verify, lookup-or-setup branch), `HandleSetupContext`, `HandleSetup` |
| `internal/auth/store.go` | Extract token issuance into helpers: `IssueSessionToken(userID, role)` (existing behavior) and `IssueSetupToken(oid, email, name)` (adds `scope:"setup"` claim, 15-min expiry) |
| `internal/auth/jwt.go` | `parseToken` rejects setup-scoped tokens on normal routes; new `SetupMiddleware()` that *requires* `scope == "setup"` |
| `internal/department/` (new, tiny) | `HandleListDepartments` — `GET /api/departments` returns `{id, name}[]` for the setup dropdown; also seeds departments if the table is empty |

### New/changed routes (`cmd/server/app.go`)

| Route | Chain | Purpose |
| ----- | ----- | ------- |
| `GET /auth/login` | `chain` (public) | Start OIDC redirect |
| `GET /auth/callback` | `chain` (public) | OIDC callback — must be reachable by the browser; no SameOrigin guard (top-level GET navigation) |
| `GET /setup` | static page | Serves `setup.html` |
| `GET /api/departments` | `chain` (public GET) | Department options |
| `GET /api/setup/context` | `api` + `SetupMiddleware` | Prefill data (email, suggested name) |
| `POST /api/setup` | `api` + `SetupMiddleware` | Create account, issue session cookie |
| `GET /api/auth-config` | `chain` (public) | `{"auth_type":"ad"|"mock"}` so `login.html` knows which UI to show |

### Data rules

- **Name**: `users.name` populated from the ID token `displayName` (falls back to `name`).
- **Role**: `users.role` (`student` | `professor`). Professors are pre-provisioned by the IT team via
  SQL (their `ad_object_id` column must be populated); students self-register through the setup flow.
- **Student ID**: `users.student_id` is required for students. The setup form prefills it from the
  email local-part if it is numeric (university convention), editable otherwise.
- **Storage**: students and professors share a single `users` table (introduced to simplify OIDC
  login and prevent self-sign-up as professors). `enrollments.user_id` → `users.id`,
  `semester_courses.taught_by` → `users.id`.

## 4. Frontend changes (`web/public/`)

- **`login.html`** — fetch `/api/auth-config`; in `ad` mode show a "Sign in with University
  Account" button linking to `/auth/login` (password form hidden); in `mock` mode keep the
  current form.
- **`setup.html`** (new) — plain HTML + htmx page matching existing style: name input, role radio,
  department `<select>` (htmx-loaded from `/api/departments`), conditional student-ID input.
  On load: fetch `/api/setup/context`; 401 → redirect `/login`. Submit: `POST /api/setup`
  (JSON) → success → redirect `/` (existing index role-redirect takes over).
- **`js/app.js`** — treat `/setup` as a valid destination for setup-token holders; keep the
  existing 401 → `/login` behavior for all other pages.

## 5. Security notes

- `state` and `nonce`: cryptographically random, stored in short-lived `HttpOnly`,
  `SameSite=Lax` cookies; `state` compared on callback, `nonce` verified inside the ID token.
- Redirect URI must match the app registration **exactly** (scheme, path, trailing slash).
- Production callback requires HTTPS (Nginx + Let's Encrypt — already a course requirement);
  Nginx must pass `/auth/*` and `/setup` through under the project path.
- Setup tokens are useless against normal `/api/*` routes, and normal sessions cannot call
  `/api/setup`.
- `AD_CLIENT_SECRET` comes from the **class Azure Key Vault** at runtime in production
  (separate course requirement); `.env` is acceptable for local dev only.

## 6. Testing

1. **Dev tenant end-to-end** (personal Entra tenant): new user → redirected to `/setup` →
   account created → logged in with chosen role; second login → straight through.
2. **Negative cases**: tampered `state` → 400; replayed/expired setup token → 401; setup token
   used on `/api/my/courses` → 401; normal token on `/api/setup` → 401.
3. **Regression**: `go build ./...`, `go vet ./...`, `staticcheck ./...` pass; `gofmt -l .` clean.
4. Regenerate Swagger: `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/server/main.go --output docs/swagger --pd`.

## 7. Rollout

1. Implement + test against personal dev tenant.
2. Receive university app registration → swap env vars (no code change).
3. Fetch secrets from Azure Key Vault at startup (prod).
4. Register production redirect URI; deploy via Docker Compose behind Nginx.
5. Update `README.md` and `docs/technical.md` (auth model section) to describe the OIDC flow.

## 8. Open questions for the instructor

- Is there a pre-registered class app, or do we self-register in the university tenant?
- Which redirect URIs are allowed (does the class have a standard domain/path)?
- Which API permissions are pre-consented?
- (Optional/future) Can role be derived from AD groups instead of self-selected at setup?
