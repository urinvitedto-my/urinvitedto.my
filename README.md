# urinvitedto.my (MVP)

Minimal e-invitation stack:

- Backend: Go (chi), JSON-only API under `/api/v1`
- Frontend: Vue 3 + Vite + TypeScript
- DB/Auth/Storage: Supabase (Postgres + Auth + Storage)

## Design System

**Style:** Modern, minimalist, elegant

**Mobile-first:** All views, pages, and components must be fully responsive with mobile as the primary target. Optimize for mobile but ensure excellent experience across all devices.

**Tools:**

- Icons: [Lucide](https://lucide.dev/icons/)
- Styling: Tailwind CSS everywhere

**Color Palette:**

```css
--black: #000000 
--oxford-blue: #14213D
--orange-web: #fca311 
--platinum: #e5e5e5
--antiflash-white: #ececec
```

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
  - `/` — landing page (Hero, Navbar, Footer)
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

## Private Event Page Components

After entering 6-char alphanumeric invite code (all caps), guests see modular sections:

**Top sections (order configurable per event):**

- Event details (title, description, hosts)
- Location photo
- Map (embedded/interactive)
- FAQ
- Event schedule/timeline
- Photo/video gallery
- Dress code
- Countdown timer
- QR code for monetary gifts
- Gift guide (physical gifts)

**Middle section:**

- Invite details: display names of guests on this invite
- RSVP form: per-guest confirm/decline with optional message

**Bottom section:**

- Confirmed guests list (all who RSVP'd yes across all invites)

## Notes

- Public events: one page with details (no guest list/RSVP).
- Private events: 6-char alphanumeric invite code (random and all caps) reveals full event page with modular components.
- Invite codes are unique per guest (tied to event).

### JSON casing

- All API JSON uses camelCase (e.g., `isPublic`, `coverImageUrl`, `startsAt`, `displayName`, `rsvpStatus`).
- Database columns remain snake_case.
