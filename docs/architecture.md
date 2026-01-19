# urinvitedto.my - Complete Architecture

> Single source of truth for the entire project. Combines planning, data modeling, and implementation details.

---

## Table of Contents

1. [Overview](#overview)
2. [Terminology](#terminology)
3. [Tech Stack](#tech-stack)
4. [Design System](#design-system)
5. [Data Model](#data-model)
6. [Modular Content Architecture](#modular-content-architecture)
7. [API Specification](#api-specification)
8. [Frontend Architecture](#frontend-architecture)
9. [Access Control & RLS](#access-control--rls)
10. [Configuration](#configuration)
11. [Conventions](#conventions)

---

## Overview

**urinvitedto.my** is a minimalist e-invitation platform. Hosts create events, generate invite codes for guests, and guests RSVP via a beautiful, mobile-first experience.

### Scope Decisions

- **Public events:** Single page with event details; no RSVP, no guest lists.
- **Private events:** Access via 6-character alphanumeric invite code (all caps); shows full event details and per-guest RSVP.
- **Hosts:** Read-only dashboard to view approved guests and messages. Cannot modify data.
- **Admin:** Full CRUD access to all data.

---

## Terminology

| Term   | Definition                                             |
| ------ | ------------------------------------------------------ |
| Host   | Celebrant(s) of the event (e.g., bride & groom)        |
| Guest  | Individual invitee (person on an invite)               |
| Invite | A party/group invitation containing one or more guests |
| Admin  | Platform administrator (you) with full access          |

---

## Tech Stack

| Layer    | Technology                                      |
| -------- | ----------------------------------------------- |
| Frontend | Vue 3 + Vite + TypeScript (SPA)                 |
| Backend  | Go (chi router), JSON API under `/api/v1`       |
| Database | Supabase (Postgres)                             |
| Auth     | Supabase Auth (email/password)                  |
| Storage  | Supabase Storage (`event-media` bucket, public) |

---

## Design System

### Style

Modern, minimalist, elegant.

### Mobile-First Responsive

All views must be mobile-first. Majority of users access via mobile. Design for:

- Mobile: 320px+
- Tablet: 768px+
- Desktop: 1024px+

### Tools

- **Icons:** [Lucide](https://lucide.dev/icons/)
- **Styling:** Tailwind CSS

### Color Palette

```css
--black: #000000; /* primary text, headers */
--oxford-blue: #14213d; /* secondary dark, accents */
--orange-web: #fca311; /* primary CTA, highlights */
--platinum: #e5e5e5; /* borders, dividers */
--antiflash-white: #ececec; /* backgrounds, cards */
```

---

## Data Model

### Entity Relationship

```plaintext
admins (email allowlist)

events
  ├── hosts (1:many - celebrants)
  ├── invites (1:many - invite codes)
  │     └── guests (1:many - individuals on each invite)
  ├── event_schedule (1:many - timeline items)
  ├── event_faqs (1:many - Q&A items)
  ├── event_gallery (1:many - photos/videos)
  └── event_gifts (1:many - gift suggestions)
```

### Core Tables

#### admins

Admin allowlist (email-based).

| Column | Type | Notes |
| ------ | ---- | ----- |
| email  | text | PK    |

#### events

Main event table with fixed columns for core data.

| Column             | Type        | Notes                                  |
| ------------------ | ----------- | -------------------------------------- |
| id                 | uuid        | PK, auto-generated                     |
| type               | text        | `wedding` \| `birthday` \| `party`     |
| slug               | text        | URL-friendly identifier                |
| title              | text        | Event name                             |
| description        | text        | Nullable, main event description       |
| is_public          | boolean     | Default false                          |
| cover_image_url    | text        | Nullable, hero image for event landing |
| location_photo_url | text        | Nullable, venue/location image         |
| starts_at          | timestamptz | Event start date/time                  |
| location           | text        | Venue name/address                     |
| custom_content     | jsonb       | Nullable, flexible event-specific data |
| enabled_components | jsonb       | Nullable, component visibility & order |
| created_at         | timestamptz | Auto-generated                         |

**Indexes:**

- `unique(type, slug)`

#### hosts

Celebrants tied to events. Multiple hosts per event allowed (e.g., bride & groom).

| Column        | Type        | Notes                    |
| ------------- | ----------- | ------------------------ |
| id            | uuid        | PK                       |
| event_id      | uuid        | FK → events.id           |
| display_name  | text        | Name shown on event page |
| contact_email | text        | Nullable                 |
| auth_user_id  | uuid        | Maps to auth.users.id    |
| created_at    | timestamptz |                          |

#### invites

One invite code per party/group per event.

| Column      | Type        | Notes                                 |
| ----------- | ----------- | ------------------------------------- |
| id          | uuid        | PK                                    |
| event_id    | uuid        | FK → events.id                        |
| invite_code | text        | 6-char alphanumeric, all caps, random |
| label       | text        | Nullable, e.g., "Smith Family"        |
| created_at  | timestamptz |                                       |

**Indexes:**

- `unique(event_id, invite_code)`

#### guests

Individual invitees. RSVP stored per guest.

| Column       | Type        | Notes                      |
| ------------ | ----------- | -------------------------- |
| id           | uuid        | PK                         |
| event_id     | uuid        | FK → events.id             |
| invite_id    | uuid        | FK → invites.id            |
| display_name | text        | Guest's name               |
| rsvp_status  | text        | `pending` \| `yes` \| `no` |
| rsvp_message | text        | Nullable, optional message |
| rsvp_at      | timestamptz | When RSVP was submitted    |
| created_at   | timestamptz |                            |

**Indexes:**

- `index(event_id, invite_id)`

### Extended Tables (Modular Components)

#### event_schedule

Timeline of event activities.

| Column      | Type        | Notes                |
| ----------- | ----------- | -------------------- |
| id          | uuid        | PK                   |
| event_id    | uuid        | FK → events.id       |
| time        | timestamptz | When activity occurs |
| title       | text        | Activity name        |
| description | text        | Nullable             |
| order_index | int         | Display order        |

**Indexes:**

- `index(event_id, order_index)`

#### event_faqs

Frequently asked questions.

| Column      | Type | Notes          |
| ----------- | ---- | -------------- |
| id          | uuid | PK             |
| event_id    | uuid | FK → events.id |
| question    | text |                |
| answer      | text |                |
| order_index | int  | Display order  |

**Indexes:**

- `index(event_id, order_index)`

#### event_gallery

Photos and videos for the event.

| Column      | Type | Notes                |
| ----------- | ---- | -------------------- |
| id          | uuid | PK                   |
| event_id    | uuid | FK → events.id       |
| media_type  | text | `photo` \| `video`   |
| media_url   | text | Supabase Storage URL |
| caption     | text | Nullable             |
| order_index | int  | Display order        |

**Indexes:**

- `index(event_id, order_index)`

#### event_gifts

Gift suggestions (physical or monetary).

| Column      | Type | Notes                           |
| ----------- | ---- | ------------------------------- |
| id          | uuid | PK                              |
| event_id    | uuid | FK → events.id                  |
| gift_type   | text | `physical` \| `monetary`        |
| title       | text | Gift name or payment method     |
| description | text | Nullable                        |
| link        | text | Nullable, registry/payment link |
| order_index | int  | Display order                   |

**Indexes:**

- `index(event_id, order_index)`

---

## Modular Content Architecture

### The 3-Tier Hybrid Approach

We use a combination of fixed columns, structured tables, and flexible JSONB:

| Tier | Location                  | Purpose                          | Example                     |
| ---- | ------------------------- | -------------------------------- | --------------------------- |
| 1    | Fixed columns in `events` | Core data, fast queries, indexed | title, starts_at, location  |
| 2    | Separate tables           | Repeatable, queryable content    | schedule, FAQs, gallery     |
| 3    | JSONB columns             | Flexible, unique content         | dress code, custom sections |

### custom_content JSONB Schema

Stores flexible, event-specific configuration.

```typescript
interface CustomContent {
  // dress code configuration
  dressCode?: {
    title: string // e.g., "Garden Formal"
    description: string
    notes?: string
    examples?: string[] // e.g., ["Long dresses", "Suits"]
  }

  // location details beyond the address
  locationDetails?: {
    parkingInfo?: string
    accessibilityNotes?: string
    mapEmbedUrl?: string // Google Maps embed URL
  }

  // monetary gift configuration
  monetaryGifts?: {
    enabled: boolean
    qrCodeUrl?: string // QR code image URL
    instructions?: string
    accounts?: Array<{
      method: string // e.g., "GCash", "PayMaya"
      number: string
      name: string
    }>
  }

  // countdown timer customization
  countdownTimer?: {
    enabled: boolean
    customMessage?: string // e.g., "Until we say 'I do'"
  }

  // custom sections (the power feature)
  customSections?: Array<{
    id: string
    title: string
    content: string // HTML content
    image?: string
    bgColor?: string
    order: number
  }>
}
```

### enabled_components JSONB Schema

Controls which components show and in what order.

```typescript
interface EnabledComponents {
  components: Array<{
    name: string // Component name (e.g., "EventDetails")
    enabled: boolean
    order: number
  }>
}
```

**Default component order:**

| Order | Component      | Data Source                       |
| ----- | -------------- | --------------------------------- |
| 1     | EventDetails   | events table                      |
| 2     | LocationPhoto  | events.location_photo_url         |
| 3     | CountdownTimer | events.starts_at + custom_content |
| 4     | EventMap       | custom_content.locationDetails    |
| 5     | EventSchedule  | event_schedule table              |
| 6     | EventGallery   | event_gallery table               |
| 7     | DressCode      | custom_content.dressCode          |
| 8     | EventFAQ       | event_faqs table                  |
| 9     | MonetaryGifts  | custom_content.monetaryGifts      |
| 10    | GiftGuide      | event_gifts table                 |
| 11    | CustomSections | custom_content.customSections     |

**Fixed sections (always at bottom):**

- InviteRSVP - guest names + RSVP form
- ConfirmedGuests - list of confirmed attendees

---

## API Specification

Base URL: `/api/v1`

### Public Endpoints

#### GET `/events/:type/:slug/summary`

Returns basic event info to determine public/private.

**Response:**

```json
{
  "id": "uuid",
  "type": "wedding",
  "slug": "john-jane-2024",
  "title": "John & Jane's Wedding",
  "isPublic": false,
  "coverImageUrl": "https://...",
  "startsAt": "2024-06-15T14:00:00Z",
  "location": "The Garden Venue"
}
```

#### GET `/events/:type/:slug/details`

Returns full event details.

**Query params:**

- `invite` (required for private events): 6-char invite code

**Response (private event with valid invite):**

```json
{
  "event": {
    "id": "uuid",
    "type": "wedding",
    "slug": "john-jane-2024",
    "title": "John & Jane's Wedding",
    "description": "Join us for our special day...",
    "isPublic": false,
    "coverImageUrl": "https://...",
    "locationPhotoUrl": "https://...",
    "startsAt": "2024-06-15T14:00:00Z",
    "location": "The Garden Venue",
    "customContent": { ... },
    "enabledComponents": { ... }
  },
  "hosts": [
    { "id": "uuid", "displayName": "John Doe" },
    { "id": "uuid", "displayName": "Jane Doe" }
  ],
  "schedule": [
    { "id": "uuid", "time": "2024-06-15T14:00:00Z", "title": "Ceremony", "description": "..." }
  ],
  "faqs": [
    { "id": "uuid", "question": "What's the dress code?", "answer": "..." }
  ],
  "gallery": [
    { "id": "uuid", "mediaType": "photo", "mediaUrl": "https://...", "caption": "..." }
  ],
  "gifts": [
    { "id": "uuid", "giftType": "physical", "title": "...", "link": "..." }
  ],
  "invite": {
    "id": "uuid",
    "label": "Smith Family",
    "guests": [
      { "id": "uuid", "displayName": "John Smith", "rsvpStatus": "pending" },
      { "id": "uuid", "displayName": "Jane Smith", "rsvpStatus": "yes", "rsvpMessage": "..." }
    ]
  },
  "confirmedGuestsCount": 42
}
```

#### GET `/events/:type/:slug/confirmed-guests`

Returns list of confirmed guest names.

**Response:**

```json
{
  "guests": [{ "displayName": "John Smith" }, { "displayName": "Jane Smith" }],
  "count": 42
}
```

#### POST `/events/:type/:slug/rsvp`

Submit RSVP for a guest.

**Request body:**

```json
{
  "inviteCode": "ABC123",
  "guestId": "uuid",
  "status": "yes",
  "message": "Can't wait!"
}
```

**Response:**

```json
{
  "id": "uuid",
  "displayName": "John Smith",
  "rsvpStatus": "yes",
  "rsvpMessage": "Can't wait!",
  "rsvpAt": "2024-01-15T10:30:00Z"
}
```

---

## Frontend Architecture

### Directory Structure

```plaintext
frontend/src/
├── views/
│   ├── LandingView.vue         # / (homepage)
│   ├── EventLandingView.vue    # /:type/:slug (public details or invite entry)
│   ├── GuestView.vue           # /:type/:slug/guest (private event page)
│   ├── HostLoginView.vue       # /host/login
│   ├── HostDashboardView.vue   # /host/dashboard
│   └── AdminView.vue           # /admin
├── components/
│   ├── Navbar.vue
│   ├── Footer.vue
│   └── event/
│       ├── EventDetails.vue
│       ├── LocationPhoto.vue
│       ├── EventMap.vue
│       ├── EventSchedule.vue
│       ├── EventGallery.vue
│       ├── DressCode.vue
│       ├── CountdownTimer.vue
│       ├── EventFAQ.vue
│       ├── MonetaryGifts.vue
│       ├── GiftGuide.vue
│       ├── CustomSection.vue
│       ├── InviteRSVP.vue
│       └── ConfirmedGuests.vue
├── services/
│   ├── api.ts                  # Backend API calls
│   └── supabase.ts             # Supabase client
├── types/
│   └── index.ts                # TypeScript interfaces
├── router/
│   └── index.ts
├── stores/                     # Pinia stores (if needed)
├── assets/
│   └── style.css
├── App.vue
└── main.ts
```

### Routes

| Path                                     | View              | Description                                  |
| ---------------------------------------- | ----------------- | -------------------------------------------- |
| `/`                                      | LandingView       | Homepage with hero, features, CTA            |
| `/:type(wedding\|birthday\|party)/:slug` | EventLandingView  | Public event page or invite code entry       |
| `/:type/:slug/guest`                     | GuestView         | Private event page (requires `?invite=CODE`) |
| `/host/login`                            | HostLoginView     | Host authentication                          |
| `/host/dashboard`                        | HostDashboardView | Host's read-only dashboard                   |
| `/admin`                                 | AdminView         | Admin CRUD interface                         |

### Component Details

#### Top Section Components (modular, order configurable)

| Component      | Data Source                    | Description                                 |
| -------------- | ------------------------------ | ------------------------------------------- |
| EventDetails   | events table                   | Title, type badge, description, hosts, date |
| LocationPhoto  | events.location_photo_url      | Hero-style venue image with overlay         |
| EventMap       | custom_content.locationDetails | Embedded map + address + directions CTA     |
| EventSchedule  | event_schedule table           | Vertical timeline of activities             |
| EventGallery   | event_gallery table            | Masonry/grid with lightbox                  |
| DressCode      | custom_content.dressCode       | Attire guidelines with icons                |
| CountdownTimer | events.starts_at               | Real-time countdown                         |
| EventFAQ       | event_faqs table               | Accordion-style Q&A                         |
| MonetaryGifts  | custom_content.monetaryGifts   | QR code + payment info                      |
| GiftGuide      | event_gifts table              | Gift suggestions with links                 |
| CustomSection  | custom_content.customSections  | Host-defined HTML sections                  |

#### Fixed Section Components (always at bottom)

| Component       | Description                                             |
| --------------- | ------------------------------------------------------- |
| InviteRSVP      | Guest names on this invite + per-guest Yes/No + message |
| ConfirmedGuests | List of all confirmed attendees across all invites      |

---

## Access Control & RLS

### Policy Summary

| Table          | Admin     | Host     | Public  |
| -------------- | --------- | -------- | ------- |
| events         | Full CRUD | Read own | Via API |
| hosts          | Full CRUD | Read own | Via API |
| invites        | Full CRUD | Read own | Via API |
| guests         | Full CRUD | Read own | Via API |
| event_schedule | Full CRUD | Read own | Via API |
| event_faqs     | Full CRUD | Read own | Via API |
| event_gallery  | Full CRUD | Read own | Via API |
| event_gifts    | Full CRUD | Read own | Via API |
| admins         | Read      | -        | -       |

### RLS Strategy

- **Admin:** Check `exists (select 1 from admins where email = auth.email())`
- **Host:** Check `exists (select 1 from hosts where event_id = table.event_id and auth_user_id = auth.uid())`
- **Public:** No direct access. Frontend calls backend API; backend uses service role key (bypasses RLS).

---

## Configuration

### Backend Environment Variables

File: `backend/.env`

```bash
PORT=8080
ENV=dev
DATABASE_URL=postgresql://postgres:[password]@[host]:5432/postgres
SUPABASE_URL=https://YOUR_PROJECT.supabase.co
SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE
CORS_ORIGIN=http://localhost:5173
```

### Frontend Environment Variables

File: `frontend/.env`

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_SUPABASE_URL=https://YOUR_PROJECT.supabase.co
VITE_SUPABASE_ANON_KEY=YOUR_ANON_KEY
```

---

## Conventions

### JSON Casing

- **API responses:** camelCase (e.g., `isPublic`, `coverImageUrl`, `startsAt`, `displayName`, `rsvpStatus`)
- **Database columns:** snake_case (e.g., `is_public`, `cover_image_url`, `starts_at`)

### Invite Code Format

- 6 characters
- Alphanumeric (A-Z, 0-9)
- All uppercase
- Random, unique per event
- Examples: `ABC123`, `XYZ789`, `JW2024`

### Event Types

Currently supported:

- `wedding`
- `birthday`
- `party`

Extensible via `:type` route parameter.

---

## Implementation Priority

### Phase 1: Core Infrastructure

1. Database schema (all tables, indexes, RLS)
2. Backend API endpoints
3. Frontend routing and basic views

### Phase 2: Event Display

1. Event landing page (public/private detection)
2. Invite code entry
3. Guest page with all modular components

### Phase 3: RSVP Flow

1. InviteRSVP component
2. RSVP API endpoint
3. ConfirmedGuests component

### Phase 4: Host Dashboard

1. Host authentication
2. Read-only dashboard (guests, RSVPs, messages)

### Phase 5: Admin Dashboard

1. Admin authentication
2. Full CRUD for events, hosts, invites, guests
3. Component enable/disable/reorder UI

### Phase 6: Polish

1. Custom sections with rich text editor
2. Drag-and-drop reordering
3. Image uploads via Supabase Storage
