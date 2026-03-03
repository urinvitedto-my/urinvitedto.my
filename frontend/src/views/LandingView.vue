<script setup lang="ts">
// Landing page - homepage with hero section
</script>

<template>
  <div class="landing-view">
    <!-- Hero Section -->
    <section class="bg-[#14213d] text-white pt-34 pb-20 px-4 mt-16">
      <div class="max-w-4xl mx-auto text-center">
        <h1 class="text-4xl md:text-6xl font-bold mb-6">
          Digital Invitations
        </h1>
        <p class="text-xl md:text-2xl text-[#e5e5e5] mb-8">
          Create e-invitations for weddings, birthdays, and parties.
        </p>
        <a
          href="/host/login"
          class="inline-block bg-[#fca311] text-black font-semibold px-8 py-4 rounded-lg hover:bg-[#e5930f] transition-colors"
        >
          Get Started
        </a>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-24 px-4">
      <div class="max-w-6xl mx-auto">
        <h2 class="text-3xl font-bold text-center text-[#14213d] mb-12">
          Why Choose Us
        </h2>
        <div class="grid md:grid-cols-3 gap-8">
          <div class="bg-white p-6 rounded-lg shadow-sm text-center">
            <div
              class="w-16 h-16 bg-[#fca311] rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <svg
                class="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-[#14213d] mb-2">Easy to Create</h3>
            <p class="text-gray-600">Simple setup with beautiful templates</p>
          </div>
          <div class="bg-white p-6 rounded-lg shadow-sm text-center">
            <div
              class="w-16 h-16 bg-[#fca311] rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <svg
                class="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-[#14213d] mb-2">Mobile First</h3>
            <p class="text-gray-600">Perfect on any device, anywhere</p>
          </div>
          <div class="bg-white p-6 rounded-lg shadow-sm text-center">
            <div
              class="w-16 h-16 bg-[#fca311] rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <svg
                class="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-[#14213d] mb-2">Track RSVPs</h3>
            <p class="text-gray-600">Know who's coming in real-time</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped></style>
