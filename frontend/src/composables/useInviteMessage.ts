import { computed, type Ref } from 'vue'
import type { AdminInvite } from '@/types'

interface InviteMessageEvent {
  type: string
  slug: string
  isPublic: boolean
}

interface InviteMessageHost {
  displayName: string
}

const MULTI_HOST_TYPES = ['wedding'] as const

/**
 * Composable for building copy-pasteable invite messages.
 * Works in both host and admin contexts — just pass in the event + hosts refs.
 */
export function useInviteMessage(
  event: Ref<InviteMessageEvent | null>,
  hosts: Ref<InviteMessageHost[]>,
) {
  const eventUrl = computed(() => {
    if (!event.value) return ''
    return `${window.location.origin}/${event.value.type}/${event.value.slug}`
  })

  const visibleHosts = computed(() =>
    hosts.value.filter((h) => h.displayName !== 'Admin'),
  )

  /** Formats host names with possessive (e.g. "Jester & Sione's"). */
  function formatHostNames(): string {
    const names = visibleHosts.value.map((h) => h.displayName)
    if (names.length === 0) return ''
    if (names.length === 1) return `${names[0]}'s`
    const last = names.pop()
    return `${names.join(', ')} & ${last}'s`
  }

  /**
   * Combines host names and event type into a display label.
   * "our wedding" for multi-host types, "{host}'s party" otherwise.
   */
  function formatEventWithHostNames(): string {
    const eventType = event.value?.type ?? 'event'
    if (MULTI_HOST_TYPES.includes(eventType as (typeof MULTI_HOST_TYPES)[number])) {
      return `our ${eventType}`
    }
    return `${formatHostNames()} ${eventType}`
  }

  /** Returns the opening lines tailored to the event type. */
  function buildMessageBody(eventType: string, label: string, eventLabel: string): string[] {
    if (eventType === 'wedding') {
      return [
        `Hello ${label}! ✨`,
        '',
        `We're so happy to share that you're invited to celebrate ${eventLabel} with us!`,
        'Your presence would truly mean a lot as we mark this special moment together.',
      ]
    }
    if (eventType === 'birthday') {
      return [
        `Hey ${label}! 🎂`,
        '',
        `You're invited to celebrate ${eventLabel}!`,
        "Come join the fun — it's going to be a day to remember!",
      ]
    }
    return [
      `Hey ${label}! 🎉`,
      '',
      `You're invited to ${eventLabel} — and it wouldn't be the same without you!`,
      "Get ready for a great time. We can't wait to see you there!",
    ]
  }

  /** Returns the closing line(s) with an emoji matching the event tone. */
  function getClosingLines(eventType: string): string[] {
    if (eventType === 'wedding') {
      const [a, b] = visibleHosts.value.map((h) => h.displayName)
      return ['💙', '', 'Love,', `${a} and ${b}`]
    }
    if (eventType === 'birthday') return ['🎈']
    return ['🥳']
  }

  /** Generates the full invite message text for a given invite. */
  function buildInviteMessage(invite: AdminInvite): string {
    const eventType = event.value?.type ?? 'event'
    const isPrivate = event.value && !event.value.isPublic
    const label = invite.label || 'there'
    const eventLabel = formatEventWithHostNames()

    const lines = buildMessageBody(eventType, label, eventLabel)

    lines.push('', 'View your invitation here:', eventUrl.value)

    if (isPrivate) {
      const guestNames = invite.guests.map((g) => `- ${g.displayName}`)
      lines.push('', `Invite code: ${invite.inviteCode}`, '', 'Invitation for:', ...guestNames)
    }

    lines.push(
      '',
      "At the end of the invitation, you'll find a quick RSVP. Just let us know if you can make it.",
      ...getClosingLines(eventType),
    )

    return lines.join('\n')
  }

  return { buildInviteMessage }
}
