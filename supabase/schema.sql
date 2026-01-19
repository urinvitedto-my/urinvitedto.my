-- urinvitedto.my Supabase Schema
-- Complete database schema with all tables, indexes, and RLS policies
-- ============================================================================
-- CORE TABLES
-- ============================================================================
-- Admins allowlist (email-based)
create table if not exists public.admins (email text primary key);
-- Events (main table with fixed columns for core data)
create table if not exists public.events (
  id uuid primary key default gen_random_uuid(),
  type text not null check (type in ('wedding', 'birthday', 'party')),
  slug text not null,
  title text not null,
  description text,
  is_public boolean not null default false,
  cover_image_url text,
  location_photo_url text,
  starts_at timestamptz,
  location text,
  custom_content jsonb,
  enabled_components jsonb,
  created_at timestamptz not null default now(),
  unique(type, slug)
);
-- Hosts (celebrants) - tied to Supabase user via auth_user_id
create table if not exists public.hosts (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  display_name text not null,
  contact_email text,
  auth_user_id uuid,
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
  rsvp_status text not null default 'pending' check (rsvp_status in ('pending', 'yes', 'no')),
  rsvp_message text,
  rsvp_at timestamptz,
  created_at timestamptz not null default now()
);
-- ============================================================================
-- EXTENDED TABLES (Modular Components)
-- ============================================================================
-- Event Schedule (timeline of activities)
create table if not exists public.event_schedule (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  time timestamptz not null,
  title text not null,
  description text,
  order_index int not null default 0,
  created_at timestamptz not null default now()
);
-- Event FAQs (frequently asked questions)
create table if not exists public.event_faqs (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  question text not null,
  answer text not null,
  order_index int not null default 0,
  created_at timestamptz not null default now()
);
-- Event Gallery (photos and videos)
create table if not exists public.event_gallery (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  media_type text not null check (media_type in ('photo', 'video')),
  media_url text not null,
  caption text,
  order_index int not null default 0,
  created_at timestamptz not null default now()
);
-- Event Gifts (physical and monetary gift suggestions)
create table if not exists public.event_gifts (
  id uuid primary key default gen_random_uuid(),
  event_id uuid not null references public.events(id) on delete cascade,
  gift_type text not null check (gift_type in ('physical', 'monetary')),
  title text not null,
  description text,
  link text,
  order_index int not null default 0,
  created_at timestamptz not null default now()
);
-- ============================================================================
-- INDEXES
-- ============================================================================
-- Core table indexes
create index if not exists idx_hosts_event on public.hosts(event_id);
create index if not exists idx_hosts_auth_user on public.hosts(auth_user_id);
create index if not exists idx_invites_event on public.invites(event_id);
create index if not exists idx_guests_event on public.guests(event_id);
create index if not exists idx_guests_invite on public.guests(invite_id);
create index if not exists idx_guests_event_invite on public.guests(event_id, invite_id);
-- Extended table indexes
create index if not exists idx_event_schedule_event_order on public.event_schedule(event_id, order_index);
create index if not exists idx_event_faqs_event_order on public.event_faqs(event_id, order_index);
create index if not exists idx_event_gallery_event_order on public.event_gallery(event_id, order_index);
create index if not exists idx_event_gifts_event_order on public.event_gifts(event_id, order_index);
-- ============================================================================
-- ROW LEVEL SECURITY
-- ============================================================================
-- Enable RLS on all tables
alter table public.admins enable row level security;
alter table public.events enable row level security;
alter table public.hosts enable row level security;
alter table public.invites enable row level security;
alter table public.guests enable row level security;
alter table public.event_schedule enable row level security;
alter table public.event_faqs enable row level security;
alter table public.event_gallery enable row level security;
alter table public.event_gifts enable row level security;
-- ============================================================================
-- RLS POLICIES: ADMINS (full access)
-- ============================================================================
-- Admins can see admins table
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'admins'
    and policyname = 'admins_select_admins'
) then create policy "admins_select_admins" on public.admins for
select using (
    exists (
      select 1
      from public.admins a
      where a.email = auth.email()
    )
  );
end if;
end $$;
-- Admins full access on events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'events'
    and policyname = 'admins_all_events'
) then create policy "admins_all_events" on public.events for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on hosts
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'hosts'
    and policyname = 'admins_all_hosts'
) then create policy "admins_all_hosts" on public.hosts for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on invites
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'invites'
    and policyname = 'admins_all_invites'
) then create policy "admins_all_invites" on public.invites for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on guests
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'guests'
    and policyname = 'admins_all_guests'
) then create policy "admins_all_guests" on public.guests for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on event_schedule
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_schedule'
    and policyname = 'admins_all_event_schedule'
) then create policy "admins_all_event_schedule" on public.event_schedule for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on event_faqs
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_faqs'
    and policyname = 'admins_all_event_faqs'
) then create policy "admins_all_event_faqs" on public.event_faqs for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on event_gallery
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_gallery'
    and policyname = 'admins_all_event_gallery'
) then create policy "admins_all_event_gallery" on public.event_gallery for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- Admins full access on event_gifts
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_gifts'
    and policyname = 'admins_all_event_gifts'
) then create policy "admins_all_event_gifts" on public.event_gifts for all using (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
) with check (
  exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
-- ============================================================================
-- RLS POLICIES: HOSTS (read-only on their own events)
-- ============================================================================
-- Hosts can select their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'events'
    and policyname = 'hosts_select_own_events'
) then create policy "hosts_select_own_events" on public.events for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = events.id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- Hosts can select invites for their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'invites'
    and policyname = 'hosts_select_own_invites'
) then create policy "hosts_select_own_invites" on public.invites for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = invites.event_id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- Hosts can select guests for their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'guests'
    and policyname = 'hosts_select_own_guests'
) then create policy "hosts_select_own_guests" on public.guests for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = guests.event_id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- Hosts can select schedule for their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_schedule'
    and policyname = 'hosts_select_own_schedule'
) then create policy "hosts_select_own_schedule" on public.event_schedule for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = event_schedule.event_id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- Hosts can select FAQs for their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_faqs'
    and policyname = 'hosts_select_own_faqs'
) then create policy "hosts_select_own_faqs" on public.event_faqs for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = event_faqs.event_id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- Hosts can select gallery for their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_gallery'
    and policyname = 'hosts_select_own_gallery'
) then create policy "hosts_select_own_gallery" on public.event_gallery for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = event_gallery.event_id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- Hosts can select gifts for their own events
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'public'
    and tablename = 'event_gifts'
    and policyname = 'hosts_select_own_gifts'
) then create policy "hosts_select_own_gifts" on public.event_gifts for
select using (
    exists (
      select 1
      from public.hosts h
      where h.event_id = event_gifts.event_id
        and h.auth_user_id = auth.uid()
    )
  );
end if;
end $$;
-- ============================================================================
-- STORAGE
-- ============================================================================
-- Storage bucket for event media (public read)
insert into storage.buckets (id, name, public)
values ('event-media', 'event-media', true) on conflict (id) do nothing;
-- Storage policy: anyone can read from event-media bucket
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'storage'
    and tablename = 'objects'
    and policyname = 'public_read_event_media'
) then create policy "public_read_event_media" on storage.objects for
select using (bucket_id = 'event-media');
end if;
end $$;
-- Storage policy: admins can upload to event-media bucket
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'storage'
    and tablename = 'objects'
    and policyname = 'admins_upload_event_media'
) then create policy "admins_upload_event_media" on storage.objects for
insert with check (
    bucket_id = 'event-media'
    and exists (
      select 1
      from public.admins a
      where a.email = auth.email()
    )
  );
end if;
end $$;
-- Storage policy: admins can update event-media bucket
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'storage'
    and tablename = 'objects'
    and policyname = 'admins_update_event_media'
) then create policy "admins_update_event_media" on storage.objects for
update using (
    bucket_id = 'event-media'
    and exists (
      select 1
      from public.admins a
      where a.email = auth.email()
    )
  );
end if;
end $$;
-- Storage policy: admins can delete from event-media bucket
do $$ begin if not exists (
  select 1
  from pg_policies
  where schemaname = 'storage'
    and tablename = 'objects'
    and policyname = 'admins_delete_event_media'
) then create policy "admins_delete_event_media" on storage.objects for delete using (
  bucket_id = 'event-media'
  and exists (
    select 1
    from public.admins a
    where a.email = auth.email()
  )
);
end if;
end $$;
