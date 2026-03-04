<script setup lang="ts">
defineProps<{
  size?: 'sm' | 'md' | 'lg'
}>()

const sizeClasses: Record<string, string> = {
  sm: 'h-6 w-6 border-2',
  md: 'h-8 w-8 border-4',
  lg: 'h-12 w-12 border-4',
}
</script>

<template>
  <div
    role="status"
    aria-label="Loading"
    :class="['animate-spin rounded-full border-accent border-t-transparent', sizeClasses[size ?? 'lg']]"
  ></div>
</template>
