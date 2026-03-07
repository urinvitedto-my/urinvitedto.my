// API response types matching backend models

export type EventType = 'wedding' | 'birthday' | 'party'

export interface EventSummary {
  id: string
  type: EventType
  slug: string
  title: string
  isPublic: boolean
  coverImageUrl?: string
  musicUrl?: string
  startsAt?: string
  location?: string
}

export interface Event {
  id: string
  type: EventType
  slug: string
  title: string
  description?: string
  isPublic: boolean
  coverImageUrl?: string
  locationPhotoUrl?: string
  musicUrl?: string
  startsAt?: string
  location?: string
  customContent?: CustomContent
  enabledComponents?: EnabledComponents
  createdAt: string
}

export interface Host {
  id: string
  displayName: string
}

export interface Guest {
  id: string
  displayName: string
  rsvpStatus: 'pending' | 'yes' | 'no'
  rsvpMessage?: string
  rsvpAt?: string
}

export interface Invite {
  id: string
  label?: string
  guests: Guest[]
}

export interface ScheduleItem {
  id: string
  time: string
  title: string
  description?: string
  orderIndex: number
}

export interface FAQ {
  id: string
  question: string
  answer: string
  orderIndex: number
}

export interface GalleryItem {
  id: string
  mediaType: 'photo' | 'video'
  mediaUrl: string
  caption?: string
  orderIndex: number
}

export interface Gift {
  id: string
  giftType: 'physical' | 'monetary'
  title: string
  description?: string
  link?: string
  orderIndex: number
}

export interface EventDetailsResponse {
  event: Event
  hosts: Host[]
  schedule: ScheduleItem[]
  faqs: FAQ[]
  gallery: GalleryItem[]
  gifts: Gift[]
  invite?: Invite
  confirmedGuestsCount: number
}

export interface ConfirmedGuestsResponse {
  guests: { displayName: string }[]
  count: number
}

export interface RSVPRequest {
  inviteCode: string
  guestId: string
  status: 'yes' | 'no'
  message?: string
}

export interface RSVPResponse {
  id: string
  displayName: string
  rsvpStatus: 'pending' | 'yes' | 'no'
  rsvpMessage?: string
  rsvpAt?: string
}

// JSONB types for custom_content and enabled_components

export interface DressCode {
  title: string
  description: string
  notes?: string
  examples?: string[]
}

export interface LocationDetails {
  parkingInfo?: string
  accessibilityNotes?: string
  mapEmbedUrl?: string
}

export interface MonetaryGiftsConfig {
  enabled: boolean
  accounts?: {
    method: string
    number?: string
    name: string
    qrCodeUrl?: string
  }[]
}

export interface CountdownTimerConfig {
  enabled: boolean
  customMessage?: string
}

export interface CustomSection {
  id: string
  title: string
  content: string
  image?: string
  bgColor?: string
  order: number
}

export interface CustomContent {
  dressCode?: DressCode
  locationDetails?: LocationDetails
  monetaryGifts?: MonetaryGiftsConfig
  countdownTimer?: CountdownTimerConfig
  customSections?: CustomSection[]
}

export interface ComponentConfig {
  name: string
  enabled: boolean
  order: number
}

export interface EnabledComponents {
  components: ComponentConfig[]
}

// --- Host Types ---

export interface HostEvent {
  id: string
  type: EventType
  slug: string
  title: string
  isPublic: boolean
  startsAt?: string
  location?: string
  createdAt: string
}

export interface HostGuest {
  id: string
  displayName: string
  rsvpStatus: 'pending' | 'yes' | 'no'
  rsvpMessage: string | null
  rsvpAt: string | null
}

// --- Admin Types ---

export interface AdminHost {
  id: string
  displayName: string
  contactEmail?: string
  authUserId?: string
}

export interface AdminGuest {
  id: string
  displayName: string
  rsvpStatus: 'pending' | 'yes' | 'no'
  rsvpMessage?: string
  rsvpAt?: string
  createdAt: string
}

export interface AdminInvite {
  id: string
  inviteCode: string
  label?: string
  createdAt: string
  guests: AdminGuest[]
}

export interface AdminScheduleItem {
  id: string
  time: string
  title: string
  description?: string
  orderIndex: number
  createdAt: string
}

export interface AdminFAQ {
  id: string
  question: string
  answer: string
  orderIndex: number
  createdAt: string
}

export interface AdminGift {
  id: string
  giftType: 'physical' | 'monetary'
  title: string
  description?: string
  link?: string
  orderIndex: number
  createdAt: string
}

export interface AdminGalleryItem {
  id: string
  mediaType: 'photo' | 'video'
  mediaUrl: string
  caption?: string
  orderIndex: number
  createdAt: string
}

export interface AdminEvent {
  id: string
  type: EventType
  slug: string
  title: string
  description?: string
  isPublic: boolean
  coverImageUrl?: string
  locationPhotoUrl?: string
  musicUrl?: string
  startsAt?: string
  location?: string
  createdAt: string
  hosts: AdminHost[]
}
