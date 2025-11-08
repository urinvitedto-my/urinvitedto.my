# MVP plan: urinvitedto.my with Go (chi) + Vue (Vite TS) + Supabase

## Terminology

- Host: celebrant(s) of the event
- Guest: invitee (individual person)
- Admin: you only (admin dashboard access)

## Scope decisions

- Public events: single page with event details; no RSVP, no guest lists.
- Private events: access via `invite_code`; after entering, show full event details and per-guest RSVP at the bottom.
- Hosts can log in to a read-only dashboard to view guests who approved and their messages. Admin retains full CRUD; hosts cannot modify data.

## Stack and separation

- Frontend: Vue 3 + Vite + TypeScript (SPA). Backend serves JSON APIs only.
- Backend: Go (chi) JSON API under `/api/v1`.
- DB/Auth/Storage: Supabase (Postgres + Auth email/password + Storage). Images in `event-media` bucket (public read).

## Data model (Supabase)

- events: id (uuid pk), type (`wedding|birthday|party`), slug (text), title (text), is_public (bool), cover_image_url (text), starts_at (timestamptz), location (text), created_at.
- hosts: id, event_id, display_name, contact_email (nullable), auth_user_id (uuid). // celebrant(s) tied to Supabase user
- invites: id, event_id, invite_code (text unique per event), label (text nullable), created_at.
- guests: id, event_id, invite_id (fk), display_name, rsvp_status (`pending|yes|no`), rsvp_message (text nullable), rsvp_at (timestamptz).
- admins: email (text pk) — allowlist of admin accounts (you).
- Indexes: unique (type, slug) on events; unique (event_id, invite_code) on invites; index guests(event_id, invite_id).

## Access control & RLS

- Enable RLS on all tables.
- Admin: full read/write on events/hosts/invites/guests via policy `exists (select 1 from admins where email = auth.email())`.
- Host (read-only): allow select on events/invites/guests where a `hosts` row exists with `hosts.event_id = table.event_id` and `hosts.auth_user_id = auth.uid()`. Deny insert/update/delete for hosts.
- Public consumption: frontend calls backend APIs; backend uses service role and enforces public/private rules.

## Backend API (chi)

Base URL: `/api/v1`

- GET `/events/:type/:slug/summary` → title, cover_image_url, is_public, starts_at, location.
- GET `/events/:type/:slug/details`
- Public: returns full event details (no guest info).
- Private: requires `invite` query param or `X-Invite-Code` header; returns event details plus guests for that invite.
- POST `/events/:type/:slug/rsvp`
- Body: `{ inviteCode: string, guestId: string, status: "yes" | "no", message?: string }`
- Effect: verifies invite and guest belong to event; updates `guests.rsvp_status`, `guests.rsvp_message`, and `guests.rsvp_at`.

Admin CRUD

- Implemented in frontend via Supabase JS (RLS gated by `admins`). No separate backend admin endpoints for MVP.

Host dashboard data

- Implement via Supabase JS with RLS (read-only for hosts). Queries limited by policies to host’s events.

## Frontend routes (Vue)

- `/:type(wedding|birthday|party)/:slug`
- If event is public: render details directly.
- If private: show invite code input; on submit, navigate to `/:type/:slug/guest?invite=...`.
- `/:type/:slug/guest`
- Requires `invite` query param for private; loads details (including invite’s guest list); bottom shows per-guest RSVP (Yes/No) with optional message field.
- `/host/login` → Supabase Auth email/password.
- `/host/dashboard` → list of events for the logged-in host; click through to `/host/events/:eventId` to view approved guests and their messages (filter rsvp_status = yes; toggle to see all).
- `/admin` → admin login; admin CRUD pages for events/hosts/invites/guests via Supabase JS.

## Frontend structure

- `frontend/src/`
  - `pages/`
    - `EventLanding.vue` — `/:type/:slug` (public details or invite entry)
    - `GuestPage.vue` — `/:type/:slug/guest` (private details + guest list + RSVP)
    - `HostLogin.vue`, `HostDashboard.vue`
    - `AdminPage.vue`
  - `router/index.ts` — routes above
  - `services/`
    - `api.ts` — wraps backend API calls using `VITE_API_BASE_URL`
    - `supabase.ts` — Supabase client using `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY`
  - `App.vue` — `<router-view />`

Implementation notes

- Fetch flow:
  - Landing page: call `GET /events/:type/:slug/summary` to decide public/private.
  - Public: load details with `GET /events/:type/:slug/details`.
  - Private: prompt for invite code; navigate to guest page with `?invite=CODE`.
  - Guest page: load details (includes invite’s guests); POST RSVP to `/events/:type/:slug/rsvp`.
- Auth (hosts/admin): Supabase email/password; hosts are read-only via RLS.

## Minimal UI

- Mobile-first; cover hero, simple info blocks (date/time, location). Private guest page shows list of guests with RSVP controls and an optional message textarea per guest.

## Configuration & env

- Backend: `PORT`, `DATABASE_URL`, `SUPABASE_URL`, `SUPABASE_SERVICE_ROLE_KEY`, `CORS_ORIGIN`.
- Frontend: `VITE_API_BASE_URL`, `VITE_SUPABASE_URL`, `VITE_SUPABASE_ANON_KEY`.
- Seed `admins` with your email; create host user accounts and link `hosts.auth_user_id` to their `auth.users.id`.

## Conventions

- API JSON uses camelCase: `isPublic`, `coverImageUrl`, `startsAt`, `displayName`, `rsvpStatus`, etc.
- DB tables/columns use snake_case.

## Deliverables

- Backend in `backend/` with chi server and the endpoints above.
- Frontend in `frontend/` with routes and pages above (public/private event pages, host dashboard, admin dashboard).
- Supabase SQL in `backend/docs/supabase/schema.sql` with tables, indexes, and RLS policies for admin and host read-only.
- README with setup and env variables.

## Notes

- No public attendee lists or public RSVPs.
- Event types extend via `:type` route segment.
