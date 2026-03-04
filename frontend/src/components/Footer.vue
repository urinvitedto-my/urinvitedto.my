<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="bg-surface text-primary py-8 px-4">
    <div class="max-w-6xl mx-auto">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <div class="text-center md:text-left">
          <div class="h-10 w-40 overflow-hidden rounded mb-1">
            <img
              src="/nav_logo.png"
              alt="urinvitedto.my"
              class="h-full w-full object-cover object-center"
            />
          </div>
          <div class="text-primary text-sm pl-2">Digital invitations</div>
        </div>

        <div class="text-primary text-sm text-center md:text-right">
          <p>&copy; {{ currentYear }} urinvitedto.my</p>
          <p>All rights reserved</p>
        </div>
      </div>
    </div>
  </footer>
</template>
