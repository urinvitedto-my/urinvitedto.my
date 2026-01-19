<script setup lang="ts">
// 404 page
</script>

<template>
  <div class="not-found-view min-h-screen flex items-center justify-center py-12 px-4">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-[#14213d] mb-4">404</h1>
      <p class="text-xl text-gray-600 mb-8">Page not found</p>
      <RouterLink
        to="/"
        class="inline-block bg-[#fca311] text-black font-semibold px-6 py-3 rounded-lg hover:bg-[#e5930f] transition-colors"
      >
        Go Home
      </RouterLink>
    </div>
  </div>
</template>

<style scoped></style>
