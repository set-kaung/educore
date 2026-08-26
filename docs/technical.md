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

- [htmx](https://htmx.org) + the `client-side-templates` extension + Mustache render JSON API responses into tables/selects declaratively in the HTML
- `js/app.js` — ES module glue: session check via `GET /api/session`, role-based page guards and nav rendering, login/logout helpers, error toasts, 401 → `/login` redirect for htmx requests
- Per-page `<script type="module">` blocks handle forms that need JSON payloads or redirects
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

- One JWT source (`auth.Login`); API clients use `Authorization: Bearer`, browsers get it in the `edu_token` cookie (`HttpOnly`, `SameSite=Lax`) set by `POST /api/login`
- Token extraction order: header first, cookie fallback (`internal/auth/jwt.go`)
- Unmatched paths return JSON 404 under `/api|/public|/admin`, otherwise the static `404.html`

## Cross-site protection (`internal/web/sameorigin.go`)

State-changing requests are validated with `Origin` / `Sec-Fetch-Site`: cross-site browser requests are rejected with 403, while non-browser API clients (no Origin header, bearer/API-key auth) pass through. This replaces the old signed double-submit CSRF tokens, which required server-rendered forms.

## Role routing

- `index.html` checks the session client-side and redirects: students → `/my-courses`, staff → `/students`, guests → `/login`
- Server-side enforcement stays in the API: `RoleRequired` gates staff endpoints; students only see their own enrollments via `GET /api/my/courses` (`claims.UserID`)
