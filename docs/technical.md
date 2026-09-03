# Document Generation

Run
`swag init -g cmd/server/main.go --output docs/swagger --parseInternal --pd`

# Web UI Architecture

Single binary serves the JSON API and a static HTML frontend from disk. UI and business logic are separated: pages never render data server-side; the browser calls `/api/*` endpoints.

## Static assets (`internal/web/static.go`, `web/public/`)

- Pages are plain HTML files under `web/public/` served at clean URLs (`/login` → `login.html`, …)
- `web.NewStaticFiles(dir)` resolves every request against the static root, following symlinks, and rejects anything that escapes the root (path traversal) or is not a regular file (no directory listings)
- `GET /css/` and `GET /js/` mount the matching directories under `web/public/` with short-lived cache headers
- The root directory is configured via `STATIC_DIR` (default `web/public`; the Docker image copies it to `/srv/web/public`)

## Frontend stack (`web/public/`)

- No framework and no vendored libraries: pages load data with `fetch` and render JSON API responses into tables/selects via the shared `apiGet` / `renderRows` / `renderOptions` helpers in `js/app.js` (all interpolated values are HTML-escaped with `esc`)
- `js/app.js` — ES module glue: session check via `GET /api/session`, role-based page guards and nav rendering, login/logout helpers, error toasts, 401 → `/login` redirect on API calls
- Per-page `<script type="module">` blocks wire up forms, debounced search and semester filtering
- Plain CSS split across `css/base.css` (tokens + reset), `css/layout.css` (topbar, nav, page shell), `css/components.css` (cards, buttons, forms, tables, toasts) and `css/helpers.css`; semantic class names, no build step

## Request chains (`cmd/server/app.go`)

```
chain      = RequestLogMiddleWare
api        = chain + SameOrigin              (same-origin guard on unsafe methods)
protected  = api + JWT Middleware            (JSON errors)
professorOnly/adminOnly = protected + RoleRequired(...)
```

| Middleware       | Behavior on failure                        |
| ---------------- | ------------------------------------------ |
| SameOrigin       | 403 for cross-site unsafe methods          |
| JWT Middleware   | JSON 401                                   |
| RoleRequired     | JSON 403                                   |

## Auth model

- OIDC Authorization Code Flow via Microsoft Entra ID (`internal/auth/oidc_*.go`):
  `GET /auth/login` redirects to Microsoft, `GET /auth/callback` verifies `state` + `nonce` + ID token (JWKS / `iss` / `aud` / `exp` via go-oidc) and links the directory `oid` claim to `ad_object_id` on a single `users` row
- Known users (`ad_object_id` already present in the DB) receive a session JWT immediately; first-time users receive a 15-minute setup-scoped token (`scope=setup`) valid only on `/api/setup*`; `POST /api/setup` creates the `users` row with `role='student'` — session middleware rejects setup tokens and vice versa
- Professors are pre-provisioned via SQL (their `ad_object_id` column must be populated); students self-register through the setup flow (`POST /api/setup` always creates a `role='student'` row)
- One JWT source (`auth.IssueSessionToken`); API clients use `Authorization: Bearer`, browsers get it in the `edu_token` cookie (`HttpOnly`, `SameSite=Lax`)
- Token extraction order: header first, cookie fallback (`internal/auth/jwt.go`)
- Required env vars: `AD_CLIENT_ID`, `AD_TENANT_ID`, `AD_CLIENT_SECRET`, `AD_REDIRECT_URI` (redirect URI must match the app registration exactly); server fails fast at startup if any is missing or Entra discovery is unreachable
- Unmatched paths return JSON 404 under `/api|/public|/admin`, otherwise the static `404.html`

## Semester model

- `semesters` table: `name` (unique), `is_current`, optional `start_date`/`end_date`. Seeded with
  Spring/Summer/Fall 2025–2026 if the table is empty. `GET /api/semesters` lists them (current first).
- `semester_courses.semester_id` → `semesters.id`. Offerings join `semesters` for the display name.
- **Enrollment gating**: students may only enroll in offerings whose semester is the single
  `is_current = true` row (`GetCurrent`). Creating offerings is not restricted to the current
  semester, so professors can schedule ahead.

## Cross-site protection (`internal/web/sameorigin.go`)

State-changing requests are validated with `Origin` / `Sec-Fetch-Site`: cross-site browser requests are rejected with 403, while non-browser API clients (no Origin header, bearer/API-key auth) pass through. This replaces the old signed double-submit CSRF tokens, which required server-rendered forms.

## Role routing

- `index.html` checks the session client-side and redirects: students → `/my-courses`, staff → `/students`, guests → `/login`
- Server-side enforcement stays in the API: `RoleRequired` gates staff endpoints; students only see their own enrollments via `GET /api/my/courses` (`claims.UserID`)
