# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Git workflow

**Never run `git commit` or `git push` without explicit confirmation from the user.** After making changes, summarise what was done and ask before running any git commands.

## Repositories

This project spans two repos that are developed together:
- **`/Wivvus/api`** — Go/Gin backend (this repo)
- **`/Wivvus/web`** — Angular 17 frontend (`/home/phil/proyectos/src/github.com/Wivvus/web`)

Issues and tasks are tracked in the **[Run Wivvus](https://github.com/orgs/Wivvus/projects/2)** GitHub project.

## Commands

### Running the API
```bash
# Load env vars and run:
source wivvus.env && go run cmd/api/main.go

# Build then run:
go build -o /tmp/wivvus-api ./cmd/api/ && source wivvus.env && /tmp/wivvus-api

# Build check only:
go build ./...

# Start Postgres via podman (if not already running):
make start-postgres
```

### Git remotes
Both repos use `upstream` as the remote name → `git push upstream main`

## Architecture

**Entry point:** `cmd/api/main.go` — initialises auth, JWT, DB, metrics, reminder scheduler, then starts Gin on `:8080`.

**Package layout under `internal/`:**

| Package | Responsibility |
|---------|---------------|
| `app/` | Wires together sub-routers from `auth`, `events`, `ratings` |
| `app/auth/` | `/auth/*` and `/user/*` routes — registration, login, Google OAuth, avatar upload, password management |
| `app/events/` | `/event/*` and `/events` routes — CRUD, attend/drop, attendees list |
| `app/ratings/` | `/event/:id/rate`, `/event/:id/ratings`, `/user/:id/ratings` |
| `models/` | GORM models + repo structs. `ConnectDB` runs `AutoMigrate`. Single package-level `db *gorm.DB` shared by all repos. |
| `middleware/` | `AuthRequired()` — verifies local JWT or Google OIDC token; sets `"user"` in gin context. `InitAuth` must be called at startup with `GOOGLE_OAUTH_CLIENT_ID`. |
| `tokens/` | HS256 JWT sign/verify, 30-day expiry |
| `storage/` | DigitalOcean Spaces (S3-compatible) for avatar storage. Google avatars are lazily copied to Spaces on first authenticated request. |
| `reminders/` | Background goroutine sending email reminders 12 hours before events and rating-request emails 12 hours after |
| `email/` | Transactional email sending |
| `metrics/` | PostHog server-side event tracking |

**Models pattern:**
- Each domain type (`Event`, `User`, `Attendance`, `Rating`, `EventOption`, etc.) has a `*Repo` struct with methods that accept/return plain structs
- `ToAPI()` methods produce API-safe decorator structs (e.g. `EventAPIDecorator`) that join related data — creator name/avatar, attendee count, creator rating, options
- Never expose GORM model fields (like `DeletedAt`) directly in API responses

**Auth flow:**
- Two token types share `AuthRequired()`: local JWT (`X-Auth-Provider: local` header) and Google OIDC tokens
- Google tokens are verified via `go-oidc` against `accounts.google.com`; the `verifier` is initialised once at startup
- Google OAuth2 code exchange: `POST /auth/google/code` with `{code, redirect_uri}` — exchanges for id_token then verifies it
- On first Google login after using a profile picture URL, the avatar is asynchronously copied to DO Spaces

**Database:**
- PostgreSQL via GORM
- Schema managed by `AutoMigrate` in `models.ConnectDB` — add new models there
- Soft deletes on `Event` and `User` (GORM's `DeletedAt`)

**Required env vars** (all in `wivvus.env`):
`GOOGLE_OAUTH_CLIENT_ID`, `GOOGLE_OAUTH_CLIENT_SECRET`, `JWT_SECRET`, `ALLOWED_ORIGINS`, `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASSWORD`, `PG_DB`, `DO_SPACES_KEY`, `DO_SPACES_SECRET`, `DO_SPACES_BUCKET`, `DO_SPACES_REGION`, `DO_SPACES_ENDPOINT`, `APP_URL`
