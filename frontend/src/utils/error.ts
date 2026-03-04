/** Extracts a user-friendly error message from an unknown catch value. */
export function errorMsg(e: unknown, fallback: string): string {
  return e instanceof Error ? e.message : fallback
}
