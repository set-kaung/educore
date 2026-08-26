# Core Features

## 1 JWT Authentication

Login with username and password to receive a JWT token. All protected endpoints require a `Bearer` token in the `Authorization` header. Browser sessions use the same JWT in an `HttpOnly` cookie.

## 2 Course Management

Professors create and maintain course listings.

## 3 Enrollment System

Students browse available courses and register/unregister.

## 4 Textbook Search

Search for syllabus textbooks via the OpenLibrary API.

## 5 Web Dashboard

Static HTML + htmx frontend served from disk (no Go templating, no embedding, no frontend build step). Pages call the JSON API directly; hand-written plain CSS. Role-based landing pages:

- **Students** see `My Courses` — their enrolled courses (`enrollments` joined with semester courses)
- **Professors/Admins** get the student list with live search

# External Integrations

## 1 IT Support Ticket System

Server errors generate tickets through a partner team's ticketing API. Tickets are tracked locally for status monitoring.

## 2 API Documentation

OpenAPI specification auto-generated from source annotations.
