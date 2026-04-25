/** How to order the guest list. */
export type GuestSortMode =
  | "name-asc"
  | "name-desc"
  | "responded-newest"
  | "responded-oldest"

export interface GuestSortableRow {
  displayName: string
  rsvpAt?: string | null
}

/** Compares two guests for stable list ordering. */
export function compareGuests(
  a: GuestSortableRow,
  b: GuestSortableRow,
  mode: GuestSortMode,
): number {
  switch (mode) {
    case "name-asc":
      return a.displayName.localeCompare(b.displayName)
    case "name-desc":
      return b.displayName.localeCompare(a.displayName)
    case "responded-newest":
      return compareByRespondedAt(a, b, "newest")
    case "responded-oldest":
      return compareByRespondedAt(a, b, "oldest")
  }
}

/** Orders by response time; ties break on name. */
function compareByRespondedAt(
  a: GuestSortableRow,
  b: GuestSortableRow,
  order: "newest" | "oldest",
): number {
  const na = rsvpSortKey(a.rsvpAt, order)
  const nb = rsvpSortKey(b.rsvpAt, order)
  if (na !== nb) {
    return order === "newest" ? nb - na : na - nb
  }
  return a.displayName.localeCompare(b.displayName)
}

/**
 * Numeric key for ordering by `rsvpAt`. Missing or invalid dates sort after real
 * timestamps for both directions.
 */
function rsvpSortKey(
  rsvpAt: string | null | undefined,
  order: "newest" | "oldest",
): number {
  if (!rsvpAt) {
    return order === "newest" ? Number.NEGATIVE_INFINITY : Number.POSITIVE_INFINITY
  }
  const t = new Date(rsvpAt).getTime()
  if (Number.isNaN(t)) {
    return order === "newest" ? Number.NEGATIVE_INFINITY : Number.POSITIVE_INFINITY
  }
  return t
}
