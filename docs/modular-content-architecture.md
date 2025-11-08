# Modular Content Architecture

## Problem Statement

Each event is unique with different content needs:

- Some events need all sections (schedule, FAQs, gallery, dress code, etc.)
- Some only need a few sections
- Hosts want custom sections that don't fit predefined templates
- Content varies wildly (text length, media, styling preferences)
- We need flexibility WITHOUT creating custom components per event

## Solution: 3-Tier Hybrid Approach

We use a combination of fixed columns, structured tables, and flexible JSONB to achieve maximum modularity without complexity.

---

## Tier 1: Core Event Fields (Fixed Columns)

**Location:** `events` table columns

**Purpose:** Essential data that EVERY event needs, optimized for queries and indexing.

**Fields:**

```sql
events:
  id (uuid pk)
  type (text)                    -- wedding|birthday|party
  slug (text)                    -- URL-friendly identifier
  title (text)                   -- Event name
  description (text nullable)    -- Main event description
  is_public (bool)               -- Public vs private event
  cover_image_url (text nullable)
  starts_at (timestamptz)        -- Event start date/time
  location (text)                -- Venue name/address
  created_at (timestamptz)
```

**Why fixed columns?**

- Fast queries: `SELECT * FROM events WHERE type = 'wedding' AND starts_at > NOW()`
- Index optimization
- Strong typing and validation
- Required for all events

---

## Tier 2: Structured Component Data (Separate Tables)

**Purpose:** Repeatable, queryable content with consistent structure.

### Tables

#### event_schedule

```sql
id (uuid pk)
event_id (uuid fk)
time (timestamptz)           -- When this activity happens
title (text)                 -- Activity name
description (text nullable)  -- Details
order_index (int)            -- Display order
```

**Example data:**

```json
[
  { "time": "2024-06-15T14:00:00Z", "title": "Ceremony", "description": "Garden venue", "order_index": 1 },
  { "time": "2024-06-15T15:30:00Z", "title": "Cocktail Hour", "description": "Patio", "order_index": 2 },
  { "time": "2024-06-15T17:00:00Z", "title": "Reception", "description": "Ballroom", "order_index": 3 }
]
```

#### event_faqs

```sql
id (uuid pk)
event_id (uuid fk)
question (text)
answer (text)
order_index (int)
```

**Example data:**

```json
[
  { "question": "What's the dress code?", "answer": "Black tie optional", "order_index": 1 },
  { "question": "Is parking available?", "answer": "Free valet parking provided", "order_index": 2 }
]
```

#### event_gallery

```sql
id (uuid pk)
event_id (uuid fk)
media_type (text)            -- photo|video
media_url (text)             -- Supabase Storage URL
caption (text nullable)
order_index (int)
```

#### event_gifts

```sql
id (uuid pk)
event_id (uuid fk)
gift_type (text)             -- physical|monetary
title (text)                 -- Gift name or payment method
description (text nullable)  -- Details
link (text nullable)         -- Registry link or payment link
order_index (int)
```

**Why separate tables?**

- Easy CRUD operations per item
- Queryable: "Show all events with >5 FAQs"
- Ordered lists with `order_index`
- Relationship integrity via foreign keys
- Can be empty if event doesn't need that component

---

## Tier 3: Flexible Custom Content (JSONB)

**Location:** `events` table JSONB columns

**Purpose:** Truly unique, one-off content that doesn't fit structured tables.

### Schema Addition

```sql
events:
  ...existing columns...
  custom_content (jsonb nullable)      -- Flexible event-specific data
  enabled_components (jsonb nullable)  -- Component visibility & order
```

### custom_content Examples

#### Example 1: Dress Code Details

```json
{
  "dressCode": {
    "title": "Garden Formal",
    "description": "Formal attire with garden-appropriate footwear",
    "notes": "Ceremony is outdoors on grass. Heel protectors recommended.",
    "examples": ["Long dresses", "Suits", "No stilettos"],
    "iconColor": "#fca311"
  }
}
```

#### Example 2: Location Details

```json
{
  "locationDetails": {
    "parkingInfo": "Free parking in Lot A. Overflow parking in Lot B (5 min walk).",
    "accessibilityNotes": "Wheelchair accessible entrance on west side. Elevator available.",
    "mapEmbedUrl": "https://www.google.com/maps/embed?pb=..."
  }
}
```

#### Example 3: Monetary Gifts

```json
{
  "monetaryGifts": {
    "enabled": true,
    "qrCodeUrl": "https://supabase.co/storage/event-media/qr-payment.png",
    "instructions": "Scan QR code or send to:",
    "accounts": [
      { "method": "GCash", "number": "09171234567", "name": "John Doe" },
      { "method": "PayMaya", "number": "09181234567", "name": "Jane Doe" }
    ]
  }
}
```

#### Example 4: Countdown Timer Config

```json
{
  "countdownTimer": {
    "enabled": true,
    "targetDate": "2024-06-15T14:00:00Z",
    "labels": {
      "days": "Days",
      "hours": "Hours", 
      "minutes": "Minutes",
      "seconds": "Seconds"
    },
    "customMessage": "Until we say 'I do'",
    "accentColor": "#fca311"
  }
}
```

#### Example 5: Custom Sections (The Magic!)

```json
{
  "customSections": [
    {
      "id": "accommodations",
      "title": "Accommodations",
      "content": "<p>We have reserved a block of rooms at the <strong>Marriott Hotel</strong>.</p><p>Use code: <code>WEDDING2024</code> for 20% off.</p>",
      "image": "https://supabase.co/storage/event-media/hotel.jpg",
      "order": 11,
      "bgColor": "#ececec"
    },
    {
      "id": "transportation",
      "title": "Transportation",
      "content": "<p>Shuttle buses will run every 30 minutes from 5:00 PM to 11:00 PM.</p><ul><li>Pick-up: Hotel main lobby</li><li>Drop-off: Venue entrance</li></ul>",
      "order": 12,
      "bgColor": "#ffffff"
    }
  ]
}
```

### enabled_components Schema

Controls which components show and in what order:

```json
{
  "components": [
    { "name": "EventDetails", "enabled": true, "order": 1 },
    { "name": "LocationPhoto", "enabled": true, "order": 2 },
    { "name": "CountdownTimer", "enabled": true, "order": 3 },
    { "name": "EventMap", "enabled": true, "order": 4 },
    { "name": "EventSchedule", "enabled": true, "order": 5 },
    { "name": "EventGallery", "enabled": true, "order": 6 },
    { "name": "DressCode", "enabled": true, "order": 7 },
    { "name": "EventFAQ", "enabled": true, "order": 8 },
    { "name": "MonetaryGifts", "enabled": true, "order": 9 },
    { "name": "GiftGuide", "enabled": false, "order": 10 },
    { "name": "CustomSections", "enabled": true, "order": 11 }
  ]
}
```

**Purpose:**

- Toggle components on/off per event
- Reorder sections dynamically
- Admin can drag-and-drop to reorder
- Frontend loops through enabled components in order

---

## Frontend Component Strategy

### Predefined Components (Reusable)

**ONE component handles ALL events of that type.**

#### EventSchedule.vue

```vue
<template>
  <section v-if="scheduleItems?.length" class="event-schedule">
    <h2>Schedule</h2>
    <div v-for="item in scheduleItems" :key="item.id" class="schedule-item">
      <time>{{ formatTime(item.time) }}</time>
      <h3>{{ item.title }}</h3>
      <p v-if="item.description">{{ item.description }}</p>
    </div>
  </section>
</template>

<script setup lang="ts">
// Reads from event_schedule table
const props = defineProps<{ scheduleItems: ScheduleItem[] }>()
</script>
```

#### DressCode.vue

```vue
<template>
  <section v-if="dressCode" class="dress-code">
    <h2>{{ dressCode.title }}</h2>
    <p>{{ dressCode.description }}</p>
    <p v-if="dressCode.notes" class="text-sm">{{ dressCode.notes }}</p>
    <ul v-if="dressCode.examples">
      <li v-for="ex in dressCode.examples" :key="ex">{{ ex }}</li>
    </ul>
  </section>
</template>

<script setup lang="ts">
// Reads from events.custom_content.dressCode JSONB
const props = defineProps<{ dressCode: DressCodeConfig }>()
</script>
```

### Generic Custom Section Component

**ONE component for ALL custom sections.**

#### CustomSection.vue

```vue
<template>
  <section 
    v-if="section"
    class="custom-section p-6 rounded-lg mb-6"
    :style="{ backgroundColor: section.bgColor || '#ececec' }"
  >
    <h2 class="text-2xl font-bold mb-4">{{ section.title }}</h2>
    
    <!-- Render HTML content -->
    <div v-html="section.content" class="prose max-w-none mb-4"></div>
    
    <!-- Optional image -->
    <img 
      v-if="section.image" 
      :src="section.image" 
      :alt="section.title"
      class="w-full rounded-lg mt-4"
    />
  </section>
</template>

<script setup lang="ts">
interface CustomSectionData {
  id: string
  title: string
  content: string // HTML string
  image?: string
  bgColor?: string
  order: number
}

const props = defineProps<{ section: CustomSectionData }>()
</script>
```

### Main Event Page Component

#### GuestPage.vue (Simplified)

```vue
<template>
  <div class="event-page">
    <!-- Loop through enabled components in order -->
    <component
      v-for="comp in orderedComponents"
      :key="comp.name"
      :is="componentMap[comp.name]"
      v-bind="getComponentProps(comp.name)"
    />
    
    <!-- Always show RSVP and confirmed guests at bottom -->
    <InviteRSVP :guests="inviteGuests" />
    <ConfirmedGuests :count="confirmedCount" />
  </div>
</template>

<script setup lang="ts">
const componentMap = {
  EventDetails: EventDetails,
  LocationPhoto: LocationPhoto,
  EventMap: EventMap,
  EventSchedule: EventSchedule,
  EventGallery: EventGallery,
  DressCode: DressCode,
  CountdownTimer: CountdownTimer,
  EventFAQ: EventFAQ,
  MonetaryGifts: MonetaryGifts,
  GiftGuide: GiftGuide,
  CustomSections: CustomSection
}

// Get enabled components from event.enabled_components JSONB
const orderedComponents = computed(() => 
  event.value.enabledComponents
    .filter(c => c.enabled)
    .sort((a, b) => a.order - b.order)
)

// Pass appropriate data to each component
const getComponentProps = (componentName: string) => {
  switch(componentName) {
    case 'EventSchedule': return { scheduleItems: event.value.schedule }
    case 'EventFAQ': return { faqs: event.value.faqs }
    case 'DressCode': return { dressCode: event.value.customContent?.dressCode }
    case 'CustomSections': return { sections: event.value.customContent?.customSections }
    // ... etc
  }
}
</script>
```

---

## Admin Dashboard Flow

### Creating an Event

1. **Fill Core Fields** (fixed columns)
   - Event type, title, slug
   - Date, location, description
   - Cover image upload

2. **Enable/Order Components** (enabled_components JSONB)
   - Toggle switches for each component
   - Drag-and-drop to reorder
   - Preview order

3. **Fill Component Content**

   **For structured components:**
   - Add Schedule: Form to add items → saves to `event_schedule` table
   - Add FAQs: Q&A pairs → saves to `event_faqs` table  
   - Upload Gallery: File uploads → saves to `event_gallery` table
   **For JSONB components:**
   - Dress Code: Form fields → saves to `custom_content.dressCode`
   - Monetary Gifts: Upload QR, add accounts → saves to `custom_content.monetaryGifts`

4. **Add Custom Sections** (custom_content.customSections JSONB)
   - Click "Add Custom Section"
   - Title input
   - Rich text editor (WYSIWYG) for content
   - Optional image upload
   - Background color picker
   - Save → adds to customSections array

### Admin UI Mock

```plaintext
┌─────────────────────────────────────────┐
│  Create Event                           │
├─────────────────────────────────────────┤
│  Title: [ John & Jane's Wedding      ] │
│  Type:  [wedding ▼]                     │
│  Date:  [2024-06-15 2:00 PM]            │
│                                         │
│  ✓ Enable/Order Components              │
├─────────────────────────────────────────┤
│  [✓] Event Details        [↕] Order: 1  │
│  [✓] Location Photo       [↕] Order: 2  │
│  [✓] Countdown Timer      [↕] Order: 3  │
│  [✓] Schedule             [↕] Order: 4  │
│  [✓] Gallery              [↕] Order: 5  │
│  [✓] Dress Code           [↕] Order: 6  │
│  [✓] FAQ                  [↕] Order: 7  │
│  [✓] Monetary Gifts       [↕] Order: 8  │
│  [ ] Gift Guide           [↕] Order: 9  │
├─────────────────────────────────────────┤
│  📅 Add Schedule Items                  │
│  [+ Add Activity]                       │
│                                         │
│  ❓ Add FAQs                            │
│  [+ Add Question]                       │
│                                         │
│  🎨 Custom Sections                     │
│  [+ Add Custom Section]                 │
│                                         │
│  [Save Event]                           │
└─────────────────────────────────────────┘
```

---

## Benefits of This Approach

### ✅ Maximum Flexibility

- Hosts get predefined components that "just work"
- Can add unlimited custom sections for unique needs
- No code changes needed for new content

### ✅ One Component, Many Events

- `EventSchedule.vue` works for ALL events
- `CustomSection.vue` handles ALL unique content
- No per-host component creation

### ✅ Performance

- Fixed columns optimized with indexes
- JSONB is fast in Postgres (GIN indexes)
- Only fetch what's enabled

### ✅ Queryability

- "Find all weddings with galleries" → query `event_gallery` table
- "Events with >10 schedule items" → JOIN and COUNT
- Can't do this with pure JSONB

### ✅ Validation

- Fixed columns: strong typing
- Structured tables: relationships enforced
- JSONB: flexible but validated on frontend

### ✅ Cost Efficiency

- Empty tables cost nothing
- JSONB null if not used
- No wasted storage

### ✅ Future-Proof

- New component type? Add new table or JSONB path
- No schema changes for content variations
- Easy to export/import events

---

## Data Flow Example

### Wedding Event with Custom "Hotel Block" Section

**Admin creates event:**

1. Fills core fields → `events` table
2. Adds 3 schedule items → `event_schedule` table (3 rows)
3. Adds 5 FAQs → `event_faqs` table (5 rows)
4. Uploads 10 photos → `event_gallery` table (10 rows)
5. Fills dress code form → `custom_content.dressCode` JSONB
6. Clicks "Add Custom Section":
   - Title: "Hotel Accommodations"
   - Rich text: "We have rooms at..."
   - Uploads hotel photo
   - Saves → `custom_content.customSections[0]` JSONB

**Database state:**

```sql
-- events table (1 row)
{
  id: "uuid-123",
  title: "John & Jane's Wedding",
  custom_content: {
    "dressCode": { ... },
    "customSections": [
      {
        "id": "hotel-block",
        "title": "Hotel Accommodations",
        "content": "<p>We have reserved...</p>",
        "image": "https://...",
        "bgColor": "#ececec"
      }
    ]
  },
  enabled_components: {
    "components": [
      { "name": "EventDetails", "enabled": true, "order": 1 },
      { "name": "EventSchedule", "enabled": true, "order": 4 },
      // ...
      { "name": "CustomSections", "enabled": true, "order": 10 }
    ]
  }
}

-- event_schedule table (3 rows)
-- event_faqs table (5 rows)
-- event_gallery table (10 rows)
```

**Frontend renders:**

- Loops through `enabled_components` in order
- Passes schedule data to `EventSchedule.vue`
- Passes FAQs to `EventFAQ.vue`
- Passes `customSections` array to `CustomSection.vue` (loops and renders each)

---

## Implementation Priority (MVP)

### Phase 1: Core + Structured Tables

1. `events` table with fixed columns
2. `event_schedule`, `event_faqs`, `event_gallery` tables
3. Basic Vue components for each
4. Admin forms to populate

### Phase 2: JSONB for Simple Config

1. Add `custom_content` JSONB column
2. Implement `DressCode.vue` reading from JSONB
3. Implement `MonetaryGifts.vue` reading from JSONB
4. Admin forms for these

### Phase 3: Component Toggles

1. Add `enabled_components` JSONB column
2. Admin UI with checkboxes to enable/disable
3. Frontend respects enabled flags

### Phase 4: Custom Sections (The Power Move)

1. Implement `CustomSection.vue` component
2. Admin form: "Add Custom Section" button
3. Rich text editor (TipTap or Quill)
4. Image upload integration
5. Save to `custom_content.customSections` array

### Phase 5: Reordering

1. Drag-and-drop in admin
2. Updates `order` fields
3. Frontend respects order

---

## Code Never Changes Per Host

**Wrong approach:**

```plaintext
// ❌ BAD: Creating components per host
WeddingHost1CustomSection.vue
WeddingHost2HotelInfo.vue
WeddingHost3Transportation.vue
```

**Right approach:**

```plaintext
// ✅ GOOD: One reusable component
CustomSection.vue (handles ALL custom content)
```

The **data** is unique per host, not the **code**.

---

## Summary

- **Fixed columns:** Core event data, fast queries
- **Structured tables:** Repeatable content (schedule, FAQs, gallery)
- **JSONB:** Flexible config and custom HTML sections
- **One component per type:** Works for all events
- **No code changes:** Just data entry via admin dashboard

This architecture gives you:

- Flexibility of a CMS
- Performance of structured data
- Simplicity of reusable components
- Zero code per host
