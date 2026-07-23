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

## Requirements

- **Go** 1.26.4 or later
- **MySQL** database server

### Go Dependencies

| Package                     | Purpose                     |
| --------------------------- | --------------------------- |
| `gorm.io/gorm`              | ORM for database operations |
| `gorm.io/driver/mysql`      | MySQL driver for GORM       |
| `github.com/golang-jwt/jwt` | JWT authentication          |

Install dependencies:

```bash
go mod tidy
```

## Running

### Local Development

Set the `PORT` environment variable and run:

```bash
PORT=8080 go run .
```

### Docker

```bash
docker build -t educore .
docker run -p 8080:8080 -e PORT=8080 educore
```
