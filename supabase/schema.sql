-- urinvitedto.my Supabase schema (MVP)
-- Hosts (celebrants), Invites (party access), Guests (individuals)

-- Admins allowlist (email-based)
create table if not exists public.admins (
  email text primary key
);

-- Events
create table if not exists public.events (
  id uuid primary key default gen_random_uuid(),
  type text not null check (type in ('wedding','birthday','party')),
  slug text not null,
  title text not null,
  is_public boolean not null default false,
  cover_image_url text,
  starts_at timestamptz,
  location text,
  created_at timestamptz not null default now(),
  unique(type, slug)
);

-- Hosts (celebrants) - tied to Supabase user via auth_user_id (no FK for simplicity)
create table if not exists public.hosts (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  display_name text not null,
  contact_email text,
  auth_user_id uuid, -- maps to auth.users.id
  created_at timestamptz not null default now()
);

-- Invites (party access) - one code per party per event
create table if not exists public.invites (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  invite_code text not null,
  label text,
  created_at timestamptz not null default now(),
  unique(event_id, invite_code)
);

-- Guests (individual invitees) - RSVP stored per guest
create table if not exists public.guests (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  invite_id uuid not null references public.invites(id) on delete cascade,
  display_name text not null,
  rsvp_status text not null default 'pending' check (rsvp_status in ('pending','yes','no')),
  rsvp_message text,
  rsvp_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists idx_invites_event on public.invites(event_id);
create index if not exists idx_guests_event on public.guests(event_id);
create index if not exists idx_guests_invite on public.guests(invite_id);

-- Enable RLS
alter table public.events enable row level security;
alter table public.hosts enable row level security;
alter table public.invites enable row level security;
alter table public.guests enable row level security;
alter table public.admins enable row level security;

-- RLS: admins have full access (email-based)
do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'events' and policyname = 'admins full access events'
  ) then
    create policy "admins full access events" on public.events
      for all
      using (exists (select 1 from public.admins a where a.email = auth.email()))
      with check (exists (select 1 from public.admins a where a.email = auth.email()));
  end if;
end $$;

do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'hosts' and policyname = 'admins full access hosts'
  ) then
    create policy "admins full access hosts" on public.hosts
      for all
      using (exists (select 1 from public.admins a where a.email = auth.email()))
      with check (exists (select 1 from public.admins a where a.email = auth.email()));
  end if;
end $$;

do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'invites' and policyname = 'admins full access invites'
  ) then
    create policy "admins full access invites" on public.invites
      for all
      using (exists (select 1 from public.admins a where a.email = auth.email()))
      with check (exists (select 1 from public.admins a where a.email = auth.email()));
  end if;
end $$;

do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'guests' and policyname = 'admins full access guests'
  ) then
    create policy "admins full access guests" on public.guests
      for all
      using (exists (select 1 from public.admins a where a.email = auth.email()))
      with check (exists (select 1 from public.admins a where a.email = auth.email()));
  end if;
end $$;

-- RLS: hosts read-only on their own events
do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'events' and policyname = 'hosts can select own events'
  ) then
    create policy "hosts can select own events" on public.events
      for select
      using (
        exists (
          select 1 from public.hosts h
          where h.event_id = events.id and h.auth_user_id = auth.uid()
        )
      );
  end if;
end $$;

do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'invites' and policyname = 'hosts can select invites for own events'
  ) then
    create policy "hosts can select invites for own events" on public.invites
      for select
      using (
        exists (
          select 1 from public.hosts h
          where h.event_id = invites.event_id and h.auth_user_id = auth.uid()
        )
      );
  end if;
end $$;

do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'guests' and policyname = 'hosts can select guests for own events'
  ) then
    create policy "hosts can select guests for own events" on public.guests
      for select
      using (
        exists (
          select 1 from public.hosts h
          where h.event_id = guests.event_id and h.auth_user_id = auth.uid()
        )
      );
  end if;
end $$;

-- RLS: admins table visible only to admins themselves (harmless; service role bypasses RLS)
do $$ begin
  if not exists (
    select 1 from pg_policies where schemaname = 'public' and tablename = 'admins' and policyname = 'admins can see admins'
  ) then
    create policy "admins can see admins" on public.admins for select
      using (exists (select 1 from public.admins a where a.email = auth.email()));
  end if;
end $$;

-- Storage bucket for event media (public read)
insert into storage.buckets (id, name, public)
values ('event-media', 'event-media', true)
on conflict (id) do nothing;


