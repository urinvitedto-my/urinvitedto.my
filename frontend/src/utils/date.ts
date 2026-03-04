/** Shared date formatting and conversion utilities. */

/** Returns true if the given string produces a valid Date. */
function isValidDate(dateStr: string): boolean {
  return !Number.isNaN(new Date(dateStr).getTime())
}

/** Formats an ISO date string for display (e.g. "Mar 4, 2026"). */
export function formatDate(dateStr?: string | null, includeTime = false): string {
  if (!dateStr || !isValidDate(dateStr)) return '-'
  const opts: Intl.DateTimeFormatOptions = {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  }
  if (includeTime) {
    opts.hour = 'numeric'
    opts.minute = '2-digit'
  }
  return new Date(dateStr).toLocaleDateString('en-US', opts)
}

/** Formats time portion of an ISO string (e.g. "Mar 4, 10:30 AM"). */
export function formatTime(isoStr: string): string {
  if (!isValidDate(isoStr)) return '-'
  return new Date(isoStr).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}

/** Full date format with weekday (e.g. "Tuesday, March 4, 2026, 10:30 AM"). */
export function formatDateFull(dateStr?: string): string {
  if (!dateStr || !isValidDate(dateStr)) return ''
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}

/** Time-only format (e.g. "10:30 AM"). */
export function formatTimeOnly(dateStr: string): string {
  if (!isValidDate(dateStr)) return '-'
  return new Date(dateStr).toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
  })
}

/** Converts a datetime-local input value to ISO 8601 string. */
export function toISO(value: string): string {
  return new Date(value).toISOString()
}

/** Returns undefined if the value is empty, otherwise converts to ISO. */
export function toISOOrUndefined(value: string): string | undefined {
  if (!value) return undefined
  return toISO(value)
}

/** Converts an ISO string to datetime-local input format (YYYY-MM-DDTHH:MM). */
export function toDatetimeLocal(isoStr?: string): string {
  if (!isoStr) return ''
  const d = new Date(isoStr)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
