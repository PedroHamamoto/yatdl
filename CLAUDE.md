# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

YATDL (Yet Another Todo List) is a Go-based web application using PostgreSQL as its database. The project follows a standard Go project layout with cmd/ for main applications and internal/ for private application code.

## Development Commands

### Database Setup

Start the PostgreSQL database:
```bash
docker-compose up -d
```

The database runs on port 7777 (mapped from container port 5432) with:
- Database: yatdl_db
- User: postgres
- Password: postgres

### Database Migrations

The project uses [golang-migrate](https://github.com/golang-migrate/migrate) via Docker for database migrations.

**Common migration commands:**

```bash
# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check current migration version
make migrate-version

# Create a new migration
make migrate-create NAME=add_todos_table

# Rollback all migrations
make migrate-down-all

# Force set version (use when migration is in dirty state)
make migrate-force VERSION=1

# Drop all tables (DANGEROUS - prompts for confirmation)
make migrate-drop
```

**Migration file structure:**
- Location: `migrations/` directory
- Naming: `<timestamp>_<description>.up.sql` and `<timestamp>_<description>.down.sql`
- The `.up.sql` file contains the forward migration (e.g., CREATE TABLE)
- The `.down.sql` file contains the rollback (e.g., DROP TABLE)
- Timestamps are Unix timestamps in seconds

**Network configuration:**
- Docker network: `yatdl_default` (auto-created by docker-compose)
- Database URL: `postgres://postgres:postgres@db:5432/yatdl_db?sslmode=disable`

## Architecture

### Project Structure

```
yatdl/
├── cmd/           # Main application entry points
├── internal/      # Private application code
├── migrations/    # Database migration files
```

### Database Schema

The application uses PostgreSQL with the following core tables:

- **users** (`migrations/0001_create_users.up.sql`): User authentication with email/password
  - id (bigserial primary key)
  - email (unique)
  - password_hash (bytea)
  - created_at (timestamp)

### Go Module

Module path: `yatdl`
Go version: 1.25

## Development Workflow

This is an early-stage project with basic structure in place. The main.go file exists but is empty, indicating the application is still being scaffolded.
