<template>
  <!-- dark：在应用壳内启用 Tailwind dark:，避免未开系统深色时出现白侧栏 + 深主区 -->
  <!-- 与 HomeView 完全一致：#0a0a0f 底色 + 三层光球 + 网格 -->
  <div
    class="dark relative min-h-screen overflow-hidden bg-[#0a0a0f] text-gray-100"
    :style="{ '--scroll-y': scrollY + 'px' }"
  >
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div
        class="app-blob app-blob-1"
        :style="{ transform: `translate(${scrollY * 0.04}px, ${-scrollY * 0.03}px)` }"
      ></div>
      <div
        class="app-blob app-blob-2"
        :style="{ transform: `translate(${-scrollY * 0.06}px, ${scrollY * 0.04}px)` }"
      ></div>
      <div
        class="app-blob app-blob-3"
        :style="{ transform: `translate(${scrollY * 0.03}px, ${scrollY * 0.05}px)` }"
      ></div>
      <div
        class="app-grid-overlay"
        :style="{ transform: `translateY(${scrollY * 0.02}px)` }"
      ></div>
    </div>

    <div class="relative z-10 min-h-screen">
      <AppSidebar />

      <div
        class="relative min-h-screen transition-all duration-300"
        :class="[sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64']"
      >
        <AppHeader />

        <main class="p-4 md:p-6 lg:p-8">
          <slot />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.user?.role === 'admin')

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

const scrollY = ref(0)
function onScroll() {
  scrollY.value = window.scrollY
}

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
  window.addEventListener('scroll', onScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})

defineExpose({ replayTour })
</script>

<style scoped>
/* 与 HomeView `.blob` / `.blob-1`… / `.grid-overlay` 完全一致 */
.app-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
}
.app-blob-1 {
  width: 500px;
  height: 500px;
  top: -128px;
  right: -128px;
  background: radial-gradient(circle, rgba(20, 184, 166, 0.18) 0%, transparent 70%);
  animation: appBlobFloat1 14s ease-in-out infinite alternate;
}
.app-blob-2 {
  width: 600px;
  height: 600px;
  bottom: -160px;
  left: -160px;
  background: radial-gradient(circle, rgba(14, 165, 233, 0.14) 0%, transparent 70%);
  animation: appBlobFloat2 18s ease-in-out infinite alternate;
}
.app-blob-3 {
  width: 400px;
  height: 400px;
  top: 33%;
  left: 50%;
  transform: translateX(-50%);
  background: radial-gradient(circle, rgba(139, 92, 246, 0.12) 0%, transparent 70%);
  animation: appBlobFloat3 22s ease-in-out infinite alternate;
}

@keyframes appBlobFloat1 {
  0%   { transform: translate(0, 0) scale(1); }
  33%  { transform: translate(-30px, 20px) scale(1.05); }
  66%  { transform: translate(20px, -15px) scale(0.97); }
  100% { transform: translate(-10px, 30px) scale(1.02); }
}
@keyframes appBlobFloat2 {
  0%   { transform: translate(0, 0) scale(1); }
  40%  { transform: translate(25px, -20px) scale(1.06); }
  80%  { transform: translate(-15px, 15px) scale(0.95); }
  100% { transform: translate(10px, -5px) scale(1.03); }
}
@keyframes appBlobFloat3 {
  0%   { transform: translateX(-50%) scale(1); }
  50%  { transform: translateX(-45%) scale(1.08); }
  100% { transform: translateX(-55%) scale(0.94); }
}

.app-grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(20, 184, 166, 0.07) 1px, transparent 1px),
    linear-gradient(90deg, rgba(20, 184, 166, 0.07) 1px, transparent 1px);
  background-size: 80px 80px;
}

@media (prefers-reduced-motion: reduce) {
  .app-blob { animation: none; }
}
</style>
