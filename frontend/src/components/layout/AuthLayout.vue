<template>
  <!-- 与 HomeView 完全一致：底色 #0a0a0f + 三层光球 + 网格 -->
  <div
    class="dark relative flex min-h-screen items-center justify-center overflow-hidden bg-[#0a0a0f] p-4 text-white"
    :style="{ '--scroll-y': scrollY + 'px' }"
  >
    <!-- ── Animated Background Mesh（与 HomeView 同参） ──────────── -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Layer 1: slow drift -->
      <div
        class="auth-blob auth-blob-1"
        :style="{ transform: `translate(${scrollY * 0.04}px, ${-scrollY * 0.03}px)` }"
      ></div>
      <!-- Layer 2: faster drift -->
      <div
        class="auth-blob auth-blob-2"
        :style="{ transform: `translate(${-scrollY * 0.06}px, ${scrollY * 0.04}px)` }"
      ></div>
      <!-- Layer 3: reverse drift -->
      <div
        class="auth-blob auth-blob-3"
        :style="{ transform: `translate(${scrollY * 0.03}px, ${scrollY * 0.05}px)` }"
      ></div>
      <!-- Grid pattern with parallax -->
      <div
        class="auth-grid-overlay"
        :style="{ transform: `translateY(${scrollY * 0.02}px)` }"
      ></div>
    </div>

    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div class="mb-8 text-center">
        <template v-if="settingsLoaded">
          <!-- Logo -->
          <div
            class="auth-logo-wrap mb-5 inline-flex max-h-28 w-auto max-w-[240px] items-center justify-center overflow-hidden rounded-2xl shadow-lg shadow-[#14b8a6]/20 ring-1 ring-[#14b8a6]/20"
          >
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="max-h-28 w-auto max-w-full object-contain" />
          </div>

          <!-- Site name — gradient text + 3D hover, entrance animation on load -->
          <h1
            class="text-gradient auth-site-name mb-2 text-3xl font-bold tracking-tight"
            :class="{ 'is-visible': titleVisible }"
          >
            {{ siteName }}
          </h1>

          <!-- Subtitle -->
          <p class="auth-subtitle text-sm text-gray-400">
            {{ siteSubtitle }}
          </p>
        </template>
      </div>

      <!-- Card Container — 与首页卡片区同风格 -->
      <div
        class="auth-card rounded-2xl border border-white/5 bg-white/[0.03] p-8 shadow-xl shadow-black/20 backdrop-blur-sm"
      >
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-8 text-center text-xs text-gray-500">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

const appStore = useAppStore()

const siteName = appStore.siteName || 'Sub2API'
const siteLogo = sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true })
const siteSubtitle = appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform'
const settingsLoaded = appStore.publicSettingsLoaded
const currentYear = new Date().getFullYear()

// ── Scroll tracking ─────────────────────────────────────────────
const scrollY = ref(0)
const titleVisible = ref(false)

function onScroll() {
  scrollY.value = window.scrollY
}

onMounted(() => {
  appStore.fetchPublicSettings()
  window.addEventListener('scroll', onScroll, { passive: true })
  // Trigger title entrance animation after a brief delay
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      titleVisible.value = true
    })
  })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
/* ── Background blobs：与 HomeView `.blob` / `.blob-1`… 完全一致 ── */
.auth-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
}
.auth-blob-1 {
  width: 500px;
  height: 500px;
  top: -128px;
  right: -128px;
  background: radial-gradient(circle, rgba(20, 184, 166, 0.18) 0%, transparent 70%);
  animation: authBlobFloat1 14s ease-in-out infinite alternate;
}
.auth-blob-2 {
  width: 600px;
  height: 600px;
  bottom: -160px;
  left: -160px;
  background: radial-gradient(circle, rgba(14, 165, 233, 0.14) 0%, transparent 70%);
  animation: authBlobFloat2 18s ease-in-out infinite alternate;
}
.auth-blob-3 {
  width: 400px;
  height: 400px;
  top: 33%;
  left: 50%;
  transform: translateX(-50%);
  background: radial-gradient(circle, rgba(139, 92, 246, 0.12) 0%, transparent 70%);
  animation: authBlobFloat3 22s ease-in-out infinite alternate;
}

@keyframes authBlobFloat1 {
  0%   { transform: translate(0, 0) scale(1); }
  33%  { transform: translate(-30px, 20px) scale(1.05); }
  66%  { transform: translate(20px, -15px) scale(0.97); }
  100% { transform: translate(-10px, 30px) scale(1.02); }
}
@keyframes authBlobFloat2 {
  0%   { transform: translate(0, 0) scale(1); }
  40%  { transform: translate(25px, -20px) scale(1.06); }
  80%  { transform: translate(-15px, 15px) scale(0.95); }
  100% { transform: translate(10px, -5px) scale(1.03); }
}
@keyframes authBlobFloat3 {
  0%   { transform: translateX(-50%) scale(1); }
  50%  { transform: translateX(-45%) scale(1.08); }
  100% { transform: translateX(-55%) scale(0.94); }
}

/* Grid overlay — 与 HomeView `.grid-overlay` 一致 */
.auth-grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(20, 184, 166, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(20, 184, 166, 0.04) 1px, transparent 1px);
  background-size: 80px 80px;
}

/* ── Logo entrance ── */
.auth-logo-wrap {
  animation: authLogoReveal 0.7s cubic-bezier(0.22, 1, 0.36, 1) 0.1s both;
}
@keyframes authLogoReveal {
  from { opacity: 0; transform: scale(0.85) translateY(10px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

/* ── Site name: gradient text (.text-gradient) + subtle entrance animation ── */
.auth-site-name {
  opacity: 0;
  transform: translateY(32px);
  transition: opacity 0.6s cubic-bezier(0.22, 1, 0.36, 1),
              transform 0.6s cubic-bezier(0.22, 1, 0.36, 1),
              filter 0.85s ease;
}
.auth-site-name.is-visible {
  opacity: 1;
  transform: translateY(0);
}
.auth-site-name:hover {
  filter: brightness(1.18);
}

/* ── Subtitle entrance ── */
.auth-subtitle {
  opacity: 0;
  transform: translateY(16px);
  transition: opacity 0.6s ease 0.15s, transform 0.6s cubic-bezier(0.22, 1, 0.36, 1) 0.15s;
}
.auth-subtitle {
  animation: authSubtitleFade 0.6s ease 0.2s both;
}
@keyframes authSubtitleFade {
  from { opacity: 0; transform: translateY(12px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ── Card entrance ── */
.auth-card {
  opacity: 0;
  transform: translateY(24px);
  animation: authCardReveal 0.65s cubic-bezier(0.22, 1, 0.36, 1) 0.3s both;
}
@keyframes authCardReveal {
  from { opacity: 0; transform: translateY(24px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ── 3D text hover for headings inside the auth card ── */
/* Respect original light/dark colors; add text-shadow on top */
.auth-card :deep(h2) {
  text-shadow: none;
  transition: all 0.35s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  cursor: pointer;
}
/* Dark mode: add subtle teal 3D text-shadow */
.dark .auth-card :deep(h2),
.auth-card :deep(h2.text-white) {
  color: #fff;
  text-shadow:
    0 0 10px rgba(20, 184, 166, 0.15),
    0 0 20px rgba(20, 184, 166, 0.08),
    1px 1px 0 rgba(0, 0, 0, 0.5),
    2px 2px 0 rgba(0, 0, 0, 0.38),
    3px 3px 10px rgba(0, 0, 0, 0.3);
}
.dark .auth-card :deep(h2:hover) {
  text-shadow:
    0 0 16px rgba(94, 234, 212, 0.4),
    0 0 32px rgba(20, 184, 166, 0.25),
    0 0 48px rgba(20, 184, 166, 0.1),
    1px 1px 0 rgba(0, 0, 0, 0.5),
    2px 2px 0 rgba(0, 0, 0, 0.4),
    3px 3px 0 rgba(0, 0, 0, 0.32),
    4px 4px 0 rgba(0, 0, 0, 0.26),
    5px 8px 22px rgba(0, 0, 0, 0.4);
  color: #99f6e4;
}

/* Reduce motion */
@media (prefers-reduced-motion: reduce) {
  .auth-blob { animation: none; }
  .auth-logo-wrap { animation: none; opacity: 1; transform: none; }
  .auth-site-name { animation: none; opacity: 1; transform: none; }
  .auth-subtitle { animation: none; opacity: 1; transform: none; }
  .auth-card { animation: none; opacity: 1; transform: none; }
}
</style>
