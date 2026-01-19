# urinvitedto.my

Minimalist e-invitation platform.

## Stack

- **Backend:** Go (chi), JSON API at `/api/v1`
- **Frontend:** Vue 3 + Vite + TypeScript
- **DB/Auth/Storage:** Supabase (Postgres + Auth + Storage)

> Full architecture details: [docs/architecture.md](./docs/architecture.md)

## Design System

**Style:** Modern, minimalist, elegant, mobile-first

**Tools:**

- Icons: [Lucide](https://lucide.dev/icons/)
- Styling: Tailwind CSS

**Color Palette:**

```css
--black: #000000;
--oxford-blue: #14213d;
--orange-web: #fca311;
--platinum: #e5e5e5;
--antiflash-white: #ececec;
```

## Setup

### 1. Supabase

Run `supabase/schema.sql` in Supabase SQL editor, then:

```sql
insert into public.admins (email) values ('you@example.com');
```

### 2. Backend

Create `backend/.env`:

```bash
PORT=8080
ENV=dev
DATABASE_URL=postgresql://postgres:[password]@[host]:5432/postgres
SUPABASE_URL=https://YOUR_PROJECT.supabase.co
SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE
CORS_ORIGIN=http://localhost:5173
```

Run:

```bash
cd backend
go mod tidy
go run .
```

### 3. Frontend

Create `frontend/.env`:

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_SUPABASE_URL=https://YOUR_PROJECT.supabase.co
VITE_SUPABASE_ANON_KEY=YOUR_ANON_KEY
```

Run:

```bash
cd frontend
npm i
npm run dev
```

## Routes

| Path                             | Description                                    |
| -------------------------------- | ---------------------------------------------- |
| `/`                              | Landing page                                   |
| `/:type/:slug`                   | Event landing (public details or invite entry) |
| `/:type/:slug/guest?invite=CODE` | Private event page                             |
| `/host/login`                    | Host authentication                            |
| `/host/dashboard`                | Host's read-only dashboard                     |
| `/admin`                         | Admin CRUD interface                           |

## API

| Method | Endpoint                                      | Description                                     |
| ------ | --------------------------------------------- | ----------------------------------------------- |
| GET    | `/api/v1/events/:type/:slug/summary`          | Basic event info                                |
| GET    | `/api/v1/events/:type/:slug/details`          | Full event details (`?invite=CODE` for private) |
| GET    | `/api/v1/events/:type/:slug/confirmed-guests` | List of confirmed attendees                     |
| POST   | `/api/v1/events/:type/:slug/rsvp`             | Submit RSVP                                     |

## Private Event Components

After entering 6-char alphanumeric invite code (all caps), guests see modular sections:

**Configurable sections (order adjustable):**

- Event details (title, description, hosts)
- Location photo
- Countdown timer
- Map (embedded/interactive)
- Event schedule/timeline
- Photo/video gallery
- Dress code
- FAQ
- Monetary gifts (QR code)
- Physical gift guide
- Custom sections (host-defined)

**Fixed sections (always at bottom):**

- Invite details + RSVP form
- Confirmed guests list

## Conventions

- API JSON: camelCase (`isPublic`, `startsAt`, `displayName`)
- DB columns: snake_case (`is_public`, `starts_at`, `display_name`)
- Invite codes: 6-char alphanumeric, all caps (e.g., `ABC123`)
