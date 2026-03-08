<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="bg-surface text-primary py-8 px-4">
    <div class="max-w-full mx-16">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <div class="h-12 w-12 overflow-hidden rounded mb-1">
          <img
            src="/footer_logo.png"
            alt="urinvitedto.my"
            class="h-full w-full object-contain scale-140"
          />
        </div>

        <div class="text-primary text-sm text-center md:text-right">
          <p>&copy; {{ currentYear }} urinvitedto.my</p>
          <p>All rights reserved</p>
        </div>
      </div>
    </div>
  </footer>
</template>
