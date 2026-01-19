<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="bg-[#14213d] text-white py-8 px-4">
    <div class="max-w-6xl mx-auto">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <div class="text-center md:text-left">
          <div class="font-bold text-lg mb-1">urinvitedto.my</div>
          <div class="text-[#e5e5e5] text-sm">Beautiful digital invitations</div>
        </div>

        <div class="text-[#e5e5e5] text-sm">
          &copy; {{ currentYear }} urinvitedto.my. All rights reserved.
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped></style>
