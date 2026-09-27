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

# Dependencies

Install dependencies:

```bash
go mod tidy
```

## Web UI

The UI is a set of plain static HTML files served from `web/public/` — no Go templating and no frontend build step. Pages are small ES modules that call the JSON API directly with `fetch` and render the data into tables/selects (HTML-escaped via a shared `esc` helper), keeping UI and business logic fully separated. Styles are hand-written plain CSS split into `base.css` (tokens + reset), `layout.css` (page shell), `components.css` (buttons, forms, tables, toasts) and a few width/spacing helpers in `helpers.css`.
