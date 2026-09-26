# EduCore

Course-registration backend for CSX 4110 Backend Application Development.

## Team Members

| Name             | Student ID |
| ---------------- | ---------- |
| Sein Tun Tuck    | 6622152    |
| Set Kaung Lwin   | 6632017    |
| Okkar Kaung Myat | 6632104    |

## Overview

EduCore is a course-registration backend with two primary roles:

- **Professors** — create and maintain course listings
- **Students** — browse, enroll, and manage course registrations

### Core Features

- **Course Management** — professors create, update, and delete course listings
- **Enrollment System** — students browse available courses and register/unregister
- **Role-Based Access Control** — different permissions for students, professors, and admins
- **JWT Authentication** — secure token-based session management
- **Subpath Hosting** — deploy under a URL prefix (e.g. `/educore`) via `BASE_PATH`

### External Integrations

- **OpenLibrary API** — lookup course syllabus books
- **HelpDesk API** — bidirectional API relationship with a partner team's ticketing service (Option 2)

### Tech Stack

- **Language:** Go 1.26+
- **ORM:** GORM
- **Database:** MySQL
- **Auth:** JWT with Role-Based Access Control (RBAC)
- **UI:** Static HTML + htmx, plain CSS — no frontend build step

## Requirements

- **Go** 1.26.4 or later
- **MySQL** database server
- **DATABASE_URL** environment variable (e.g. `root:password@tcp(127.0.0.1:3306)/educore?parseTime=true`)
- **JWT_SECRET** environment variable
- **AD_CLIENT_ID**, **AD_TENANT_ID**, **AD_CLIENT_SECRET** — from the Entra ID app registration
- **AD_REDIRECT_URI** — must exactly match a redirect URI registered on the app (e.g. `http://localhost:8080/auth/callback`)
- **BASE_PATH** (optional) — URL subpath prefix when hosted behind a reverse proxy (e.g. `/educore`)

### Go Dependencies

| Package                               | Purpose                                |
| ------------------------------------- | -------------------------------------- |
| `github.com/gin-gonic/gin`            | HTTP framework (used via chi-compat)   |
| `gorm.io/gorm`                        | ORM for database operations            |
| `gorm.io/driver/mysql`                | MySQL driver for GORM                  |
| `golang.org/x/oauth2`                 | OAuth2 / OIDC client                   |
| `github.com/golang-jwt/jwt/v5`        | JWT parsing and signing                |
| `github.com/alecthomas/kingpin/v2`    | CLI flags (if used)                    |

Install dependencies:

```bash
go mod tidy
```

## Web UI

The UI is a set of plain static HTML files served from `web/public/` — no Go templating and no frontend build step. Pages are small ES modules that call the JSON API directly with `fetch` and render the data into tables/selects (HTML-escaped via a shared `esc` helper), keeping UI and business logic fully separated. Styles are hand-written plain CSS split into `base.css` (tokens + reset), `layout.css` (page shell), `components.css` (buttons, forms, tables, toasts) and a few width/spacing helpers in `helpers.css`.

### Pages

| Route                        | Description                                           | Access          |
| ---------------------------- | ----------------------------------------------------- | --------------- |
| `GET /{$}`                   | Role-based redirect (`/my-courses`, `/students`, `/login`) | public     |
| `GET /login`                 | University (Entra ID) sign-in page                    | guests only     |
| `GET /setup`                 | First-time account setup after AD sign-in             | public          |
| `GET /students`              | Student list with live search                         | professor/admin |
| `GET /my-courses`            | Courses the logged-in student is enrolled in          | logged in       |
| `GET /course-offerings`      | Offerings filtered by semester (students can enroll)  | logged in       |
| `GET /add-course`            | Create course listings                                | professor/admin |
| `GET /add-offering`          | Schedule a course offering                            | professor/admin |
| `GET /recommended-books`     | Recommended books for a course offering               | logged in       |
| `GET /offering-detail`       | Detail view for a course offering                     | logged in       |
| `GET /css/`, `GET /js/`      | Stylesheets and page modules                          | public          |

All page routes are served as static files. When `BASE_PATH` is set (e.g. `/educore`), the app:
- Mounts all routes under that prefix
- Redirects `GET /educore` → `/educore/`
- Injects `<base href="/educore/">` into every HTML page so relative URLs resolve correctly
- Scopes cookies to the base path so they don't leak to sibling services (`/content`, `/api`)

### API Endpoints

| Endpoint                                      | Purpose                                            | Access          |
| --------------------------------------------- | -------------------------------------------------- | --------------- |
| `GET /health`                                 | Health check                                       | public          |
| `POST /api/logout`                            | Clears the session cookie                          | public + same-origin |
| `GET /api/session`                            | Current user id + role                             | logged in       |
| `GET /auth/login`                             | Start university (Entra ID) sign-in                | public          |
| `GET /auth/callback`                          | OIDC callback, issues session or setup token       | public          |
| `GET /api/setup/context`                      | Setup prefill (email, name)                        | setup token     |
| `POST /api/setup`                             | Create account from directory identity             | setup token     |
| `GET /api/departments`                        | Department options (seeds defaults)                | public          |
| `GET /api/students?q=`                        | Student list with filter                           | professor/admin |
| `GET /api/my/courses`                         | Enrolled courses for current student               | logged in       |
| `GET /api/my/taught-courses`                  | Courses taught by current professor                | professor/admin |
| `GET /api/semesters`                          | Semester records (current one flagged)             | logged in       |
| `GET /api/semester-courses?semester_id=`      | Offerings for a semester                           | logged in       |
| `GET /api/semester-courses/{id}`              | Offering detail                                    | logged in       |
| `POST /api/semester-courses`                  | Create offering (checks conflicts)                 | professor/admin |
| `DELETE /api/semester-courses/{id}`           | Soft-delete an offering                            | professor/admin |
| `POST /api/semester-courses/{id}/enroll`      | Enroll in an offering (checks schedule conflicts)  | student         |
| `GET /api/semester-courses/{id}/books`        | List recommended books for offering                | logged in       |
| `POST /api/semester-courses/{id}/books`       | Add recommended book to offering                   | professor/admin |
| `DELETE /api/books/{id}`                      | Remove recommended book                            | professor/admin |
| `GET /api/courses`                            | Course listings                                    | logged in       |
| `POST /api/courses`                           | Create course listing                              | professor/admin |
| `DELETE /api/courses/{id}`                    | Soft-delete a course listing                       | professor/admin |
| `GET /api/professors`                         | Professor directory                                | logged in       |
| `GET /api/textbooks`                          | Search OpenLibrary for textbooks                   | public          |
| `GET /public/students/{id}/departments/{dept}`| Verify student department enrollment (HelpDesk)    | API key         |
| `POST /admin/api-keys`                        | Grant API key to partner service                   | admin           |
| `DELETE /admin/api-keys/{key}`                | Revoke API key                                     | admin           |

State-changing `/api/*` requests are protected by a same-origin check (`Origin` / `Sec-Fetch-Site`). Browser sessions use the JWT delivered in an `HttpOnly`, `SameSite=Lax` cookie; API clients keep using the `Authorization: Bearer` header.

### UI Project Layout

```
web/
  public/              # served from disk at runtime
    index.html         # role-based redirect
    login.html, students.html, my-courses.html,
    course-offerings.html, add-course.html, add-offering.html,
    recommended-books.html, offering-detail.html, 404.html,
    setup.html
    favicon.ico, logo.png, MS_Logo.jpg
    css/base.css       # design tokens, reset, element defaults
    css/layout.css     # topbar, nav, page shell
    css/components.css # cards, buttons, forms, tables, toasts
    css/helpers.css    # small width/spacing/text helpers
    js/app.js          # session/nav/toast + fetch/render helpers

internal/web/          # disk-backed static file server + same-origin guard
```

## Running

### Local Development

Set the environment variables and run:

```bash
export DATABASE_URL="root:password@tcp(127.0.0.1:3306)/educore?parseTime=true"
export JWT_SECRET="secret123"
export AD_CLIENT_ID="..."
export AD_TENANT_ID="..."
export AD_CLIENT_SECRET="..."
export AD_REDIRECT_URI="http://localhost:8080/auth/callback"
export PORT=8080
# Optional: export BASE_PATH="/educore"  # when testing subpath hosting
go run ./cmd/server
```

The server starts on `http://localhost:8080`. On startup, GORM will automatically migrate all models (creating tables if they don't exist).

### Production Build

Cross-compile a static Linux binary:

```bash
./deploy.sh
```

This produces `app_linux` at the repo root. Verify with:

```bash
go build ./... && go vet ./...
```

### Behind a Reverse Proxy (Nginx)

When hosting under a subpath (e.g. `/educore`), set `BASE_PATH=/educore` and configure Nginx:

```nginx
location = /educore {
    return 301 /educore/;
}
location /educore/ {
    proxy_pass http://127.0.0.1:4000;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_cache_bypass $http_upgrade;
}
```

This keeps `/content` (WordPress) and `/api` (Lab) intact on the host. The `AD_REDIRECT_URI` must be the full callback URL including the subpath (e.g. `https://host/educore/auth/callback`).