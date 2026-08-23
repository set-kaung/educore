# Document Generation

Run
`swag init -g cmd/server/main.go --output docs/swagger --parseInternal --pd`

# Web UI Architecture

Single binary serves JSON API + server-rendered HTML. No frontend build step.

## Rendering (`internal/web/render.go`)

- `web.Files` — `embed.FS` over `web/templates` and `web/static`
- `NewRenderer` parses every page set once at startup (base layout + all partials + page), fails fast on bad templates
- `Renderer.Page` — full page with layout; `Renderer.Fragment` — bare partial for htmx responses
- `Renderer.H(h)` — adapts `func(w, r) *HTTPError` handlers to HTML: errors render `pages/error.html` instead of JSON

## Request chains (`cmd/server/app.go`)

```
chain      = RequestLogMiddleWare
guestUI    = chain + CSRF + JWT OptionalMiddleware
authedUI   = guestUI + JWT PageMiddleware("/login")
protected  = chain + JWT Middleware          (JSON API, unchanged)
```

| Middleware                | Behavior on failure            |
| ------------------------- | ------------------------------ |
| CSRF                      | 403                            |
| JWT `OptionalMiddleware`  | never rejects, sets claims     |
| JWT `PageMiddleware`      | 303 redirect to `/login`       |
| JWT `Middleware`          | JSON 401 (API)                 |

## Auth model

- One JWT source (`auth.Login`); API clients use `Authorization: Bearer`, browsers get it in the `edu_token` cookie (`HttpOnly`, `SameSite=Lax`)
- Token extraction order: header first, cookie fallback (`internal/auth/jwt.go`)
- Unmatched paths are caught by a mux wrapper rendering a styled 404

## CSRF (`internal/web/csrf.go`)

Signed double-submit token: random cookie value + HMAC(secret) embedded in forms as `_csrf`; validated with constant-time compare on unsafe methods. Secret reuses `JWT_SECRET`.

## Role routing (pages)

- `GET /{$}` and `GET /students` inspect claims (`auth.TryGetClaims`): students are redirected to `/my-courses`, professors/admins to `/students`
- `GET /my-courses` resolves the student via `claims.UserID` (= `students.id`) and queries `semestercourse.GetEnrolledCourses`
