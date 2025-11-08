# urinvitedto.my (MVP)

Minimal e-invitation stack:

- Backend: Go (chi), JSON-only API under `/api/v1`
- Frontend: Vue 3 + Vite + TypeScript
- DB/Auth/Storage: Supabase (Postgres + Auth + Storage)

## Env vars

Backend `.env` (at `backend/.env`):

```bash
PORT=8080
ENV=dev
DATABASE_URL=postgresql://postgres:[password]@[host]:5432/postgres
SUPABASE_URL=https://YOUR_PROJECT.supabase.co
SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE
CORS_ORIGIN=http://localhost:5173
```

Frontend `.env` (at `frontend/.env`):

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_SUPABASE_URL=https://YOUR_PROJECT.supabase.co
VITE_SUPABASE_ANON_KEY=YOUR_ANON_KEY
```

## Supabase setup

1) In Supabase SQL editor, run `backend/docs/supabase/schema.sql`.
2) Insert your admin email:

    ```bash
    insert into public.admins (email) values ('you@example.com');
    ```

3) Create at least one host user (Auth > Users) and link their `auth_user_id` in `public.hosts` for their event.

## Run locally

Backend:

```bash
cd backend
go mod tidy
go run .
```

Frontend:

```bash
cd frontend
npm i
npm run dev
```

## Frontend (MVP)

- SPA in `frontend/` (Vue 3 + Vite + TS)
- Routes:
  - `/:type(wedding|birthday|party)/:slug` — event landing
  - `/:type/:slug/guest` — guest view (requires `?invite=CODE` for private events)
  - `/host/login`, `/host/dashboard`
  - `/admin`
- Env in `frontend/.env` uses `VITE_API_BASE_URL`, `VITE_SUPABASE_URL`, `VITE_SUPABASE_ANON_KEY`.
- Start dev server at `http://localhost:5173`.

## API (MVP)

- GET `/api/v1/events/:type/:slug/summary`
- GET `/api/v1/events/:type/:slug/details` (private requires `?invite=CODE`)
- POST `/api/v1/events/:type/:slug/rsvp` body `{ inviteCode, guestId, status: "yes"|"no", message?: string }`

## Notes

- Public events: one page with details (no guest list/RSVP).
- Private events: invite code reveals guest list and per-guest RSVP with optional message.

### JSON casing

- All API JSON uses camelCase (e.g., `isPublic`, `coverImageUrl`, `startsAt`, `displayName`, `rsvpStatus`).
- Database columns remain snake_case.
