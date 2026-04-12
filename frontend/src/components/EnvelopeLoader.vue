<script setup lang="ts">
/**
 * Full-screen loading overlay with an animated envelope SVG.
 * Covers the entire viewport (including navbar/footer) during page data fetches.
 */
</script>

<template>
  <div
    class="fixed inset-0 z-60 bg-white flex items-center justify-center"
    role="status"
    aria-label="Loading"
  >
    <svg
      viewBox="0 0 200 170"
      class="w-28 md:w-36"
      xmlns="http://www.w3.org/2000/svg"
    >
      <defs>
        <!-- Prevents card from poking below the envelope -->
        <clipPath id="card-bounds">
          <rect x="20" y="0" width="160" height="160" />
        </clipPath>
      </defs>

      <!-- Soft shadow -->
      <ellipse cx="100" cy="163" rx="65" ry="4" fill="rgba(20,33,61,0.10)" />

      <!-- Envelope body (blue) -->
      <rect x="20" y="55" width="160" height="105" rx="8" fill="#14213d" />

      <!-- Envelope flap (renders BEFORE card so card appears in front when flap is open) -->
      <path class="flap" d="M21 56.5 L100 10 L179 56.5 Z" fill="#14213d" stroke="#000000" stroke-width="1.5" />

      <!-- Invitation card (peeking at rest; V-fold covers most of it) -->
      <g clip-path="url(#card-bounds)">
        <g class="letter">
          <rect x="32" y="125" width="136" height="110" rx="5" fill="#fca311" stroke="#c47e0a" stroke-width="1.5" />
          <rect x="48" y="143" width="76" height="2.5" rx="1.25" fill="rgba(255,255,255,0.45)" />
          <rect x="48" y="152" width="58" height="2.5" rx="1.25" fill="rgba(255,255,255,0.45)" />
          <rect x="48" y="161" width="66" height="2.5" rx="1.25" fill="rgba(255,255,255,0.45)" />
        </g>
      </g>

      <!-- Envelope front V-fold (covers card below fold line, gives card a V-shaped bottom) -->
      <path
        d="M21.5 56 L100 120 L178.5 56 V152 Q178.5 160 172 160 L28 160 Q21.5 160 21.5 152 Z"
        fill="#1a2a4d"
        stroke="#000000"
        stroke-width="1.5"
        stroke-linejoin="round"
      />
    </svg>
  </div>
</template>

<style scoped>
.letter {
  transform-box: view-box;
  animation: letter-slide 1.8s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

.flap {
  transform-box: view-box;
  transform-origin: 100px 55px;
  animation: flap-fold 1.8s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes flap-fold {
  /* sealed (flap hidden/flat) */
  0%,
  8% {
    transform: scaleY(0);
  }
  /* opens (flap lifts up) */
  25%,
  50% {
    transform: scaleY(1);
  }
  /* seals as card slides back in */
  72%,
  100% {
    transform: scaleY(0);
  }
}

@keyframes letter-slide {
  /* inside (fully hidden behind V-fold) */
  0%,
  8% {
    transform: translateY(0);
  }
  /* slides up with flap opening */
  35%,
  48% {
    transform: translateY(-100px);
  }
  /* slides back in as flap closes */
  72%,
  100% {
    transform: translateY(0);
  }
}
</style>

<!-- Unscoped: transition classes for parent <Transition name="loader-fade"> wrappers -->
<style>
.loader-fade-leave-active {
  transition: opacity 0.7s ease;
}
.loader-fade-leave-to {
  opacity: 0;
}
</style>
