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

## Design System

**Style:** Modern, minimalist, elegant

**Mobile-first responsive:** All UI components must be mobile-first. Majority of users will access via mobile, but ensure excellent experience on tablets and desktop browsers.

**Tools:**

- Icons: [Lucide](https://lucide.dev/icons/)
- Styling: Tailwind CSS

**Color Palette:**

- `--black: #000000` — primary text, headers
- `--oxford-blue: #14213d` — secondary dark, accents
- `--orange-web: #fca311` — primary CTA, highlights
- `--platinum: #e5e5e5` — borders, dividers
- `--antiflash-white: #ececec` — backgrounds, cards

## Data model (Supabase)

> **Note:** For detailed explanation of the modular content architecture, see [modular-content-architecture.md](./modular-content-architecture.md)

### Core tables

- **events:** id (uuid pk), type (`wedding|birthday|party`), slug (text), title (text), description (text nullable), is_public (bool), cover_image_url (text), starts_at (timestamptz), location (text), custom_content (jsonb nullable), enabled_components (jsonb nullable), created_at.

- **hosts:** id, event_id, display_name, contact_email (nullable), auth_user_id (uuid). // celebrant(s) tied to Supabase user

- **invites:** id, event_id, invite_code (text unique per event - 6 char random alphanumeric all caps), label (text nullable), created_at.

- **guests:** id, event_id, invite_id (fk), display_name, rsvp_status (`pending|yes|no`), rsvp_message (text nullable), rsvp_at (timestamptz).

- **admins:** email (text pk) — allowlist of admin accounts (you).

### Extended tables (for modular components)

- **event_schedule:** id, event_id, time (timestamptz), title (text), description (text nullable), order_index (int).

- **event_faqs:** id, event_id, question (text), answer (text), order_index (int).

- **event_gallery:** id, event_id, media_type (`photo|video`), media_url (text), caption (text nullable), order_index (int).

- **event_gifts:** id, event_id, gift_type (`physical|monetary`), title (text), description (text nullable), link (text nullable), order_index (int).

### Indexes

- unique (type, slug) on events
- unique (event_id, invite_code) on invites
- index guests(event_id, invite_id)
- index event_schedule(event_id, order_index)
- index event_faqs(event_id, order_index)
- index event_gallery(event_id, order_index)
- index event_gifts(event_id, order_index)

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
  - Private: requires `invite` query param (6-char alphanumeric all caps); returns:
    - Event details (title, description, hosts, dates, location, etc.)
    - Component data: schedule items, FAQs, gallery media, dress code, gift info
    - Guests for this specific invite (names, current rsvp_status)
    - Confirmed guests count (for bottom section)

- GET `/events/:type/:slug/confirmed-guests`
  - Returns list of all guest display names who RSVP'd "yes"
  - Used by ConfirmedGuests.vue component

- POST `/events/:type/:slug/rsvp`
  - Body: `{ inviteCode: string, guestId: string, status: "yes" | "no", message?: string }`
  - Effect: verifies invite and guest belong to event; updates `guests.rsvp_status`, `guests.rsvp_message`, and `guests.rsvp_at`.
  - Returns updated guest object

Admin CRUD

- Implemented in frontend via Supabase JS (RLS gated by `admins`). No separate backend admin endpoints for MVP.

Host dashboard data

- Implement via Supabase JS with RLS (read-only for hosts). Queries limited by policies to host’s events.

## Frontend routes (Vue)

- `/` → Landing page with Hero section, Navbar, Footer. Marketing/about page for urinvitedto.my.
- `/:type(wedding|birthday|party)/:slug`
- If event is public: render details directly.
- If private: show invite code input (6-char random alphanumeric all caps); on submit ("Open Invite" button), navigate to `/:type/:slug/guest?invite=...`.
- `/:type/:slug/guest`
- Requires `invite` query param for private; loads full event details with modular components in this order:
  1. **Top sections** (order configurable per event type):
     - EventDetails.vue — title, description, hosts
     - LocationPhoto.vue — location hero image
     - EventMap.vue — embedded map (Google Maps, etc.)
     - EventSchedule.vue — timeline of activities
     - EventGallery.vue — photo/video gallery
     - DressCode.vue — attire guidelines with icons
     - CountdownTimer.vue — real-time countdown to event start
     - EventFAQ.vue — collapsible Q&A
     - MonetaryGifts.vue — QR code for digital payment
     - GiftGuide.vue — physical gift suggestions with links
  2. **Middle section:**
     - InviteRSVP.vue — displays guest names on this invite; per-guest RSVP controls (Yes/No buttons + optional message textarea)
  3. **Bottom section:**
     - ConfirmedGuests.vue — list of all guests who RSVP'd "yes" across all invites (public within event)
- `/host/login` → Supabase Auth email/password.
- `/host/dashboard` → list of events for the logged-in host; click through to `/host/events/:eventId` to view approved guests and their messages (filter rsvp_status = yes; toggle to see all).
- `/admin` → admin login; admin CRUD pages for events/hosts/invites/guests via Supabase JS.

## Frontend structure

- `frontend/src/`
  - `pages/`
    - `HomePage.vue` — `/` (landing page with hero)
    - `EventLanding.vue` — `/:type/:slug` (public details or invite entry)
    - `GuestPage.vue` — `/:type/:slug/guest` (private details + guest list + RSVP)
    - `HostLogin.vue`, `HostDashboard.vue`
    - `AdminPage.vue`
  - `components/`
    - `Navbar.vue` — global navigation
    - `Footer.vue` — global footer
    - `event/` — modular event page components
      - `EventDetails.vue` — title, description, hosts (reads from events table)
      - `LocationPhoto.vue` — location hero image (reads from events.cover_image_url)
      - `EventMap.vue` — embedded map (reads from events.custom_content.mapEmbedUrl)
      - `EventFAQ.vue` — accordion-style FAQs (reads from event_faqs table)
      - `EventSchedule.vue` — timeline of event activities (reads from event_schedule table)
      - `EventGallery.vue` — photo/video gallery (reads from event_gallery table)
      - `DressCode.vue` — attire guidelines (reads from events.custom_content.dressCode)
      - `CountdownTimer.vue` — time remaining until event (reads from events.custom_content.countdownTimer)
      - `MonetaryGifts.vue` — QR code for digital gifts (reads from events.custom_content.monetaryGifts)
      - `GiftGuide.vue` — physical gift suggestions (reads from event_gifts table)
      - `CustomSection.vue` — flexible component for unique sections (reads from events.custom_content.customSections)
      - `InviteRSVP.vue` — guest names + RSVP form (per-guest confirm/decline)
      - `ConfirmedGuests.vue` — list of all confirmed attendees
  - `router/index.ts` — routes above
  - `services/`
    - `api.ts` — wraps backend API calls using `VITE_API_BASE_URL`
    - `supabase.ts` — Supabase client using `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY`
  - `App.vue` — `<router-view />` with Navbar and Footer

Implementation notes

- Fetch flow:
  - Landing page: call `GET /events/:type/:slug/summary` to decide public/private.
  - Public: load details with `GET /events/:type/:slug/details`.
  - Private: prompt for invite code; navigate to guest page with `?invite=CODE`.
  - Guest page: load details (includes invite’s guests); POST RSVP to `/events/:type/:slug/rsvp`.
- Auth (hosts/admin): Supabase email/password; hosts are read-only via RLS.

## UI Guidelines

- **Mobile-first responsive:** Design for mobile screens first (320px+), then tablet (768px+), then desktop (1024px+).
- **Modern, minimalist, elegant:** Clean layouts, ample whitespace, subtle animations.
- **Tailwind utilities:** Use Tailwind classes for all styling; avoid custom CSS unless necessary.
- **Lucide icons:** Use Lucide icon library for all iconography.
- **Modular components:** Each event section is a standalone component, easily reordered or toggled on/off.
- Event pages: cover hero, simple info blocks (date/time, location).
- Private guest page: modular sections stacked vertically, mobile-optimized.
- Navbar: site branding, links to `/host/login` and `/admin`.
- Footer: minimal footer with copyright, links.

## Private Event Component Details

### Top Section Components (modular, order configurable)

1. **EventDetails.vue**
   - Event title (h1), type badge, description
   - Host names with avatars/initials
   - Start date/time formatted elegantly

2. **LocationPhoto.vue**
   - Hero-style location image
   - Overlay with location name

3. **EventMap.vue**
   - Embedded interactive map (iframe or component)
   - Location address display
   - "Get Directions" CTA

4. **EventSchedule.vue**
   - Timeline layout (vertical on mobile)
   - Time + activity title + description
   - Icons for activity types (Lucide)

5. **EventGallery.vue**
   - Masonry or grid layout
   - Support photos and videos
   - Lightbox on click
   - Lazy loading

6. **DressCode.vue**
   - Dress code category (e.g., "Black Tie", "Casual")
   - Visual examples or description
   - Lucide icons for clothing types

7. **CountdownTimer.vue**
   - Real-time countdown: Days | Hours | Minutes | Seconds
   - Updates every second
   - Elegant number display

8. **EventFAQ.vue**
   - Accordion/collapsible sections
   - Question + expandable answer
   - Lucide chevron icons

9. **MonetaryGifts.vue**
   - QR code display (centered)
   - Payment instructions
   - Optional account details text

10. **GiftGuide.vue**
    - List or card layout
    - Gift suggestions with optional links
    - Store names, product names

### Middle Section Component

#### **InviteRSVP.vue**

- Section header: "Your Invitation"
- Display all guest names on this invite
- Per-guest controls:
  - Guest name (read-only)
  - Yes/No button toggle
  - Optional message textarea (shows when Yes/No selected)
  - Submit button per guest or batch submit
- Success/error feedback
- Disable after RSVP submitted

### Bottom Section Component

#### **ConfirmedGuests.vue**

- Section header: "Who's Coming"
- List of guest display names who RSVP'd "yes"
- Count: "X guests attending"
- Alphabetical sort
- Simple list or card layout
- Show message snippets if hosts allow

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
- Invite codes: 6-character alphanumeric (random all caps), unique per event. Format: `ABC123`, `XYZ789`, etc.
- Private event components are modular and can be toggled/reordered per event via `events.enabled_components` JSONB.
- Content for components uses hybrid approach:
  - Structured data in separate tables (schedule, FAQs, gallery, gifts)
  - Flexible/unique content in `events.custom_content` JSONB
  - See [modular-content-architecture.md](./modular-content-architecture.md) for details
- `CustomSection.vue` component handles all unique, host-specific sections without creating new components per event.
- First focus: wedding events with full component suite.
