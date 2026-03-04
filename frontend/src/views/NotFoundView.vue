<script setup lang="ts">
// 404 page
</script>

<template>
  <div class="not-found-view min-h-screen flex items-center justify-center pt-24 pb-12 px-4">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-primary mb-4">404</h1>
      <p class="text-xl text-gray-600 mb-8">Page not found</p>
      <RouterLink
        to="/"
        class="inline-block bg-accent text-black font-semibold px-6 py-3 rounded-lg hover:bg-accent-dark transition-colors"
      >
        Go Home
      </RouterLink>
    </div>
  </div>
</template>
