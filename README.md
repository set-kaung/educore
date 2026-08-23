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

### External Integrations

- **OpenLibrary API** — lookup course syllabus books
- **HelpDesk API** — bidirectional API relationship with a partner team's ticketing service (Option 2)

### Tech Stack

- **Language:** Go
- **ORM:** GORM
- **Database:** MySQL
- **Auth:** JWT with Role-Based Access Control (RBAC)
- **UI:** Server-rendered `html/template` + htmx (monolith)

## Requirements

- **Go** 1.26.4 or later
- **MySQL** database server
- **DATABASE_URL** environment variable (e.g. `root:password@tcp(127.0.0.1:3306)/educore`)
- **JWT_SECRET** environment variable
- **AUTH_TYPE** environment variable (`mock` for testing, `ad` for Active Directory)

### Go Dependencies

| Package                 | Purpose                     |
| ----------------------- | --------------------------- |
| `gorm.io/gorm`          | ORM for database operations |
| `gorm.io/driver/mysql`  | MySQL driver for GORM       |

Install dependencies:

```bash
go mod tidy
```

## Web UI

Server-rendered monolith: the same binary serves the JSON API and HTML pages. Templates and static assets are embedded via `go:embed` — no external files needed at runtime.

### Pages

| Route                   | Description                                        | Access          |
| ----------------------- | -------------------------------------------------- | --------------- |
| `GET /{$}`              | Redirects per role (`/my-courses` or `/students`)   | public          |
| `GET /login`            | Login form                                         | guests only     |
| `POST /session`         | Form login (sets `HttpOnly` session cookie)        | CSRF-protected  |
| `POST /logout`          | Clears session, redirects to `/login`              | public          |
| `GET /students`         | Student list with live search (htmx)               | professor/admin |
| `GET /students/search`  | HTML fragment for htmx search                      | professor/admin |
| `GET /my-courses`       | Courses the logged-in student is enrolled in       | logged in       |
| `GET /static/`          | CSS + vendored `htmx.min.js`                       | public          |

The JSON API routes (`POST /login`, `GET /student`, …) are unchanged. Browser sessions use the same JWT delivered in an `HttpOnly`, `SameSite=Lax` cookie; API clients keep using the `Authorization: Bearer` header.

### Test Users (`AUTH_TYPE=mock`)

| Username | Role      |
| -------- | --------- |
| `prof1`  | professor |
| `stu1`   | student   |

Passwords are ignored by the mock authenticator.

### UI Project Layout

```
web/
  templates/
    layouts/        # base layout
    partials/       # reusable fragments (htmx targets)
    pages/          # one file per page
  static/           # css + js, served at /static/

internal/web/       # rendering, CSRF, cookies, flash, middleware
internal/web/pages/ # page handlers
```

Page handlers call stores/services directly (no self-HTTP-calls). Forms are protected by signed double-submit CSRF tokens.

## Running

### Local Development

Set the environment variables and run:

```bash
export DATABASE_URL="root:password@tcp(127.0.0.1:3306)/educore"
export JWT_SECRET="secret123"
export AUTH_TYPE="mock"
export PORT=8080
go run ./cmd/server
```

The server starts on `http://localhost:8080`. On startup, GORM will automatically migrate all models (creating tables if they don't exist).

### Docker

```bash
docker build -t educore .
docker run -p 8080:8080 -e DATABASE_URL="root:password@tcp(host:3306)/educore" educore
```

### Docker Compose

Starts both MySQL and the app with everything wired up:

```bash
docker compose up --build
```

To stop and remove volumes:

```bash
docker compose down -v
```
