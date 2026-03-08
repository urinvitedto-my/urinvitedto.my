<script setup lang="ts">
/**
 * Landing page - homepage for urinvitedto.my
 */
</script>

<template>
  <div class="landing-view">
    <!-- Hero -->
    <section class="min-h-screen flex items-center px-6 md:px-16 pt-16 relative overflow-hidden">
      <div class="hero-logo-bg" aria-hidden="true">
        <img src="/footer_logo.png" alt="" class="w-full h-full object-contain scale-125" />
      </div>
      <div class="max-w-6xl w-full relative z-10">
        <p class="text-base uppercase tracking-[0.3em] text-accent font-semibold mb-6">
          Digital Invitations
        </p>
        <h1
          class="text-7xl md:text-8xl lg:text-9xl font-bold text-primary leading-[0.95] mb-8"
        >
          Invite<br />beautifully
        </h1>
        <p class="text-xl md:text-2xl text-gray-400 max-w-md mb-10">
          Weddings. Birthdays. Parties.
        </p>
        <RouterLink
          to="/host/login"
          class="inline-block bg-accent text-white font-semibold px-10 py-4 rounded-full hover:bg-accent-dark transition-colors text-xl"
        >
          Get Started
        </RouterLink>
      </div>
    </section>

    <!-- How It Works -->
    <section class="py-28 px-6 bg-white">
      <div class="max-w-5xl mx-auto">
        <p
          class="text-base uppercase tracking-[0.3em] text-accent font-semibold text-center mb-4"
        >
          How It Works
        </p>
        <h2 class="text-4xl md:text-5xl font-bold text-primary text-center mb-20">
          Three simple steps
        </h2>

        <div class="grid md:grid-cols-3 gap-16 md:gap-8">
          <div class="text-center">
            <div
              class="w-14 h-14 rounded-full bg-primary text-white flex items-center justify-center mx-auto mb-5 text-xl font-bold"
            >
              1
            </div>
            <h3 class="text-2xl font-semibold text-primary mb-2">Reach Out</h3>
            <p class="text-lg text-gray-500 leading-relaxed">
              Email us at
              <a
                href="mailto:urinvitedto.my@gmail.com"
                class="text-accent hover:text-accent-dark underline transition-colors"
                >urinvitedto.my@gmail.com</a
              >
              with your event details.
            </p>
          </div>

          <div class="text-center">
            <div
              class="w-14 h-14 rounded-full bg-primary text-white flex items-center justify-center mx-auto mb-5 text-xl font-bold"
            >
              2
            </div>
            <h3 class="text-2xl font-semibold text-primary mb-2">We Set It Up</h3>
            <p class="text-lg text-gray-500 leading-relaxed">
              We set up your event page and send you your invite link.
            </p>
          </div>

          <div class="text-center">
            <div
              class="w-14 h-14 rounded-full bg-primary text-white flex items-center justify-center mx-auto mb-5 text-xl font-bold"
            >
              3
            </div>
            <h3 class="text-2xl font-semibold text-primary mb-2">Share & Track</h3>
            <p class="text-lg text-gray-500 leading-relaxed">
              Send invitations and track RSVPs from your dashboard.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Event Types -->
    <section class="py-28 px-6">
      <div class="max-w-5xl mx-auto">
        <p
          class="text-base uppercase tracking-[0.3em] text-accent font-semibold text-center mb-4"
        >
          For Every Occasion
        </p>
        <h2 class="text-4xl md:text-5xl font-bold text-primary text-center mb-20">
          One platform, any celebration
        </h2>

        <div class="grid md:grid-cols-3 gap-6">
          <div
            class="group bg-white rounded-2xl p-8 shadow-sm hover:shadow-md transition-shadow"
          >
            <div
              class="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center mb-5"
            >
              <svg
                class="w-6 h-6 text-accent"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
                />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold text-primary mb-2">Weddings</h3>
            <p class="text-lg text-gray-500 leading-relaxed">
              Elegant invitations for your big day.
            </p>
          </div>

          <div
            class="group bg-white rounded-2xl p-8 shadow-sm hover:shadow-md transition-shadow"
          >
            <div
              class="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center mb-5"
            >
              <svg
                class="w-6 h-6 text-accent"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M12 8.25v-1.5m0 1.5c-1.355 0-2.697.056-4.024.166C6.845 8.51 6 9.473 6 10.608v2.513m6-4.871c1.355 0 2.697.056 4.024.166C17.155 8.51 18 9.473 18 10.608v2.513M15 8.25v-1.5m-6 1.5v-1.5m12 9.75-1.5.75a3.354 3.354 0 0 1-3 0 3.354 3.354 0 0 0-3 0 3.354 3.354 0 0 1-3 0 3.354 3.354 0 0 0-3 0 3.354 3.354 0 0 1-3 0L3 16.5m15-3.379a48.474 48.474 0 0 0-6-.371c-2.032 0-4.034.126-6 .371m12 0c.39.049.777.102 1.163.16 1.07.16 1.837 1.094 1.837 2.175v5.169c0 .621-.504 1.125-1.125 1.125H4.125A1.125 1.125 0 0 1 3 20.625v-5.17c0-1.08.768-2.014 1.837-2.174A47.78 47.78 0 0 1 6 13.12M12.265 3.11a.375.375 0 1 1-.53 0L12 2.845l.265.265Z"
                />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold text-primary mb-2">Birthdays</h3>
            <p class="text-lg text-gray-500 leading-relaxed">
              Fun, festive invites for any age.
            </p>
          </div>

          <div
            class="group bg-white rounded-2xl p-8 shadow-sm hover:shadow-md transition-shadow"
          >
            <div
              class="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center mb-5"
            >
              <svg
                class="w-6 h-6 text-accent"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 0 0-2.455 2.456ZM16.894 20.567 16.5 21.75l-.394-1.183a2.25 2.25 0 0 0-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 0 0 1.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 0 0 1.423 1.423l1.183.394-1.183.394a2.25 2.25 0 0 0-1.423 1.423Z"
                />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold text-primary mb-2">Parties</h3>
            <p class="text-lg text-gray-500 leading-relaxed">
              Set the vibe before the party starts.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="py-28 px-6 bg-primary">
      <div class="max-w-2xl mx-auto text-center">
        <h2 class="text-4xl md:text-5xl font-bold text-white mb-6">
          Ready to send your first invite?
        </h2>
        <p class="text-xl text-white/60 mb-6">
          Get in touch and we'll set everything up for you.
        </p>
        <a
          href="mailto:urinvitedto.my@gmail.com"
          class="inline-block bg-accent text-white font-semibold px-10 py-4 rounded-full hover:bg-accent-dark transition-colors text-xl"
        >
          Email Us
        </a>
        <p class="mt-6 text-xl text-white/40">urinvitedto.my@gmail.com</p>
      </div>
    </section>
  </div>
</template>

<style scoped>
.landing-view {
  font-family: var(--font-gelasio);
}

.hero-logo-bg {
  position: absolute;
  right: -5%;
  top: 50%;
  transform: translateY(-50%);
  width: 55%;
  height: 80%;
  opacity: 0.4;
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 75%);
  -webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 75%);
  pointer-events: none;
}
</style>
