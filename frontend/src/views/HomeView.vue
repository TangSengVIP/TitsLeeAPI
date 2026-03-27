<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
    ref="pageRef"
    class="relative min-h-screen overflow-hidden bg-[#0a0a0f] text-white"
    :style="{ '--scroll-y': scrollY + 'px' }"
  >
    <!-- ── Animated Background Mesh ──────────────────────────── -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Layer 1: slow drift -->
      <div
        class="blob blob-1"
        :style="{ transform: `translate(${scrollY * 0.04}px, ${-scrollY * 0.03}px)` }"
      ></div>
      <!-- Layer 2: faster drift -->
      <div
        class="blob blob-2"
        :style="{ transform: `translate(${-scrollY * 0.06}px, ${scrollY * 0.04}px)` }"
      ></div>
      <!-- Layer 3: reverse drift -->
      <div
        class="blob blob-3"
        :style="{ transform: `translate(${scrollY * 0.03}px, ${scrollY * 0.05}px)` }"
      ></div>
      <!-- Grid pattern with parallax -->
      <div
        class="grid-overlay"
        :style="{ transform: `translateY(${scrollY * 0.02}px)` }"
      ></div>
    </div>

    <!-- ── Header ────────────────────────────────────────────── -->
    <header class="header-animate relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="logo-wrap flex h-10 max-h-10 w-auto max-w-[200px] items-center overflow-hidden rounded-lg shadow-lg shadow-[#14b8a6]/20 ring-1 ring-[#14b8a6]/20">
            <img
              :src="siteLogo || '/logo.png'"
              alt="Logo"
              class="max-h-10 w-auto max-w-full object-contain object-left"
            />
          </div>
          <span
            class="inline-block bg-gradient-to-r from-white via-[#e0fffa] to-[#5eead4] bg-clip-text text-transparent text-lg font-bold tracking-tight drop-shadow-[0_0_14px_rgba(94,234,212,0.45)] transition-all duration-300 hover:scale-[1.02] hover:drop-shadow-[0_0_22px_rgba(94,234,212,0.6)]"
          >{{ siteName }}</span>
        </div>

        <div class="flex items-center gap-2">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-lg p-2 text-gray-400 transition-colors hover:bg-white/5 hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="ml-2 inline-flex items-center gap-2 rounded-full bg-[#14b8a6] px-4 py-1.5 text-sm font-medium text-[#0a0a0f] transition-all hover:bg-[#2dd4bf] hover:shadow-lg hover:shadow-[#14b8a6]/30"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-black/20 text-[10px] font-bold"
            >
              {{ userInitial }}
            </span>
            {{ t('home.dashboard') }}
            <svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="ml-2 inline-flex items-center rounded-full bg-[#14b8a6] px-4 py-1.5 text-sm font-medium text-[#0a0a0f] transition-all hover:bg-[#2dd4bf] hover:shadow-lg hover:shadow-[#14b8a6]/30"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- ── Main Content ─────────────────────────────────────── -->
    <main class="relative z-10 px-6 pb-24 pt-8">

      <!-- ─── Hero ─────────────────────────────────────────── -->
      <div class="hero-wrap mx-auto max-w-6xl">
        <!-- Tags -->
        <div class="mb-6 flex items-center justify-center gap-2">
          <span class="hero-tag hero-tag-1 inline-flex items-center gap-1.5 rounded-full border border-[#14b8a6]/30 bg-[#14b8a6]/10 px-3 py-1 text-xs font-medium text-[#14b8a6]">
            <span class="tag-dot"></span>
            {{ t('home.tags.subscriptionToApi') }}
          </span>
          <span class="hero-tag hero-tag-2 inline-flex items-center gap-1.5 rounded-full border border-[#0ea5e9]/30 bg-[#0ea5e9]/10 px-3 py-1 text-xs font-medium text-[#0ea5e9]">
            <span class="tag-dot"></span>
            {{ t('home.tags.stickySession') }}
          </span>
          <span class="hero-tag hero-tag-3 inline-flex items-center gap-1.5 rounded-full border border-[#8b5cf6]/30 bg-[#8b5cf6]/10 px-3 py-1 text-xs font-medium text-[#8b5cf6]">
            <span class="tag-dot"></span>
            {{ t('home.tags.realtimeBilling') }}
          </span>
        </div>

        <!-- Title -->
        <!-- 径向发光：标题背后柔和青绿光晕，整体 Hero 不显沉闷 -->
        <div class="relative mx-auto mb-4 max-w-4xl" style="padding-top: 2.5rem; padding-bottom: 2.5rem;">
          <div class="pointer-events-none absolute inset-0 -z-10 rounded-3xl bg-gradient-radial from-[#5eead4]/12 via-[#14b8a6]/6 to-transparent blur-2xl"></div>
          <h1
            ref="heroTitleRef"
            class="scroll-card hero-title group cursor-default text-center text-5xl font-bold tracking-tight md:text-6xl lg:text-7xl"
            :class="{ 'is-visible': visibleSections.heroTitle }"
          >
            <span
              class="inline-block bg-gradient-to-r from-white via-[#e0fffa] to-[#5eead4] bg-clip-text text-transparent drop-shadow-[0_0_22px_rgba(94,234,212,0.55),0_0_56px_rgba(20,184,166,0.28)] transition-[filter,transform] duration-[700ms] ease-out group-hover:scale-[1.02] group-hover:-translate-y-0.5 group-hover:drop-shadow-[0_0_32px_rgba(94,234,212,0.7),0_0_72px_rgba(20,184,166,0.38)]"
            >
              {{ siteName }}
            </span>
          </h1>
        </div>

        <!-- Subtitle -->
        <p class="hero-subtitle-1 mx-auto mb-4 max-w-2xl text-center text-lg text-gray-300 md:text-xl">
          {{ siteSubtitle }}
        </p>

        <p class="hero-subtitle-2 mx-auto mb-10 max-w-xl text-center text-base text-gray-400">
          {{ t('home.heroDescription') }}
        </p>

        <!-- CTA Buttons -->
        <div class="hero-cta mb-16 flex items-center justify-center gap-4">
          <router-link
            :to="isAuthenticated ? dashboardPath : '/register'"
            class="group inline-flex items-center gap-2 rounded-full bg-[#14b8a6] px-8 py-3 text-base font-semibold text-[#0a0a0f] shadow-lg shadow-[#14b8a6]/30 transition-all hover:bg-[#2dd4bf] hover:shadow-[#14b8a6]/50 hover:-translate-y-0.5"
          >
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
            <svg class="h-4 w-4 transition-transform group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" />
            </svg>
          </router-link>
        </div>
      </div>

      <!-- ─── Pain Points ─────────────────────────────────── -->
      <div ref="painPointsRef" class="scroll-section mx-auto mb-24 max-w-6xl">
        <div class="mb-12 text-center">
          <h2 class="mb-3 text-3xl font-bold tracking-tight text-white md:text-4xl text-3d-hover">
            {{ t('home.painPoints.title') }}
          </h2>
        </div>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div
            v-for="(item, idx) in painPointItems"
            :key="idx"
            class="scroll-card group rounded-2xl border border-white/5 bg-white/[0.03] p-6 backdrop-blur-sm transition-all duration-300 hover:border-[#ef4444]/20 hover:bg-white/[0.06]"
            :class="{ 'is-visible': visibleSections.painPoints }"
            :style="{ transitionDelay: idx ? `${idx * 80}ms` : '0ms' }"
          >
            <div
              class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl border border-[#ef4444]/20 bg-[#ef4444]/10 transition-transform group-hover:scale-110"
            >
              <component :is="item.icon" class="h-5 w-5 text-[#ef4444]" />
            </div>
            <h3 class="mb-2 text-base font-semibold text-white text-3d-hover">{{ item.title }}</h3>
            <p class="text-sm leading-relaxed text-gray-500">{{ item.desc }}</p>
          </div>
        </div>
      </div>

      <!-- ─── Solutions Steps ──────────────────────────────── -->
      <div ref="solutionsRef" class="scroll-section mx-auto mb-24 max-w-5xl">
        <div class="mb-12 text-center">
          <h2 class="mb-3 text-3xl font-bold tracking-tight text-white md:text-4xl text-3d-hover">
            {{ t('home.solutions.title') }}
          </h2>
          <p class="text-gray-400">{{ t('home.solutions.subtitle') }}</p>
        </div>
        <div class="grid gap-8 sm:grid-cols-3">
          <div
            v-for="(step, idx) in solutionSteps"
            :key="idx"
            class="scroll-card relative flex flex-col items-center text-center"
            :class="{ 'is-visible': visibleSections.solutions }"
            :style="{ transitionDelay: `${idx * 100}ms` }"
          >
            <div class="step-connector" v-if="idx < 2">
              <svg class="absolute top-7 -right-4 h-8 w-full text-[#14b8a6]/20" fill="none" viewBox="0 0 50 20">
                <path d="M0 10h45M45 5l5 5-5 5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div class="mb-5 flex h-14 w-14 items-center justify-center rounded-2xl border border-[#14b8a6]/30 bg-[#14b8a6]/10 font-mono text-xl font-bold text-[#14b8a6] shadow-lg shadow-[#14b8a6]/10">
              {{ idx + 1 }}
            </div>
            <h3 class="mb-2 text-lg font-semibold text-white text-3d-hover">{{ step.title }}</h3>
            <p class="text-sm leading-relaxed text-gray-500">{{ step.desc }}</p>
          </div>
        </div>
      </div>

      <!-- ─── Features Grid ────────────────────────────────── -->
      <div ref="featuresRef" class="scroll-section mx-auto mb-24 max-w-6xl">
        <div class="mb-12 grid gap-6 md:grid-cols-3">
          <div
            v-for="(feature, idx) in featureItems"
            :key="feature.title"
            class="scroll-card group rounded-2xl border border-white/5 bg-white/[0.03] p-7 backdrop-blur-sm transition-all duration-300 hover:border-[#14b8a6]/20 hover:bg-white/[0.06] hover:shadow-xl hover:shadow-[#14b8a6]/5"
            :class="{ 'is-visible': visibleSections.features }"
            :style="{ transitionDelay: `${idx * 100}ms` }"
          >
            <div
              class="mb-5 flex h-12 w-12 items-center justify-center rounded-xl transition-transform group-hover:scale-110"
              :class="feature.iconBg"
            >
              <component :is="feature.icon" class="h-6 w-6 text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-white text-3d-hover">{{ feature.title }}</h3>
            <p class="text-sm leading-relaxed text-gray-500">{{ feature.desc }}</p>
          </div>
        </div>
      </div>

      <!-- ─── Comparison Table ─────────────────────────────── -->
      <div ref="comparisonRef" class="scroll-section mx-auto mb-24 max-w-4xl">
        <div class="mb-12 text-center">
          <h2 class="mb-3 text-3xl font-bold tracking-tight text-white md:text-4xl text-3d-hover">
            {{ t('home.comparison.title') }}
          </h2>
        </div>
        <div
          :class="['scroll-card', 'overflow-hidden', 'rounded-2xl', 'border', 'border-white/5', 'bg-white/[0.02]', 'backdrop-blur-sm', { 'is-visible': visibleSections.comparison }]"
        >
          <table class="w-full">
            <thead>
              <tr class="border-b border-white/5">
                <th class="px-6 py-4 text-left text-sm font-medium text-gray-400">{{ t('home.comparison.headers.feature') }}</th>
                <th class="px-6 py-4 text-center text-sm font-medium text-gray-500">{{ t('home.comparison.headers.official') }}</th>
                <th class="px-6 py-4 text-center text-sm font-medium text-[#14b8a6]">{{ t('home.comparison.headers.us') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="row in comparisonRows"
                :key="row.feature"
                class="border-b border-white/5 transition-colors hover:bg-white/[0.02]"
              >
                <td class="px-6 py-4 text-sm text-gray-300">{{ row.feature }}</td>
                <td class="px-6 py-4 text-center text-sm text-gray-600">
                  <span class="inline-flex items-center gap-1">
                    <svg class="h-4 w-4 text-[#ef4444]/60" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    {{ row.official }}
                  </span>
                </td>
                <td class="px-6 py-4 text-center text-sm text-[#2dd4bf]">
                  <span class="inline-flex items-center gap-1">
                    <svg class="h-4 w-4 text-[#14b8a6]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    {{ row.us }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ─── Supported Providers ─────────────────────────── -->
      <div ref="providersRef" class="scroll-section mx-auto mb-24 max-w-4xl">
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold tracking-tight text-white md:text-3xl text-3d-hover">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-500">{{ t('home.providers.description') }}</p>
        </div>
        <div class="flex flex-wrap items-center justify-center gap-4">
          <div
            v-for="(provider, idx) in providerItems"
            :key="provider.name"
            class="scroll-card flex items-center gap-2.5 rounded-xl border border-white/10 bg-white/[0.04] px-5 py-3 backdrop-blur-sm transition-all duration-200 hover:border-white/20 hover:bg-white/[0.07]"
            :class="{ 'is-visible': visibleSections.providers }"
            :style="{ transitionDelay: `${idx * 80}ms` }"
          >
            <div class="flex h-8 w-8 items-center justify-center rounded-lg text-xs font-bold text-white" :class="provider.bg">
              {{ provider.initial }}
            </div>
            <span class="text-sm font-medium text-gray-300">{{ provider.name }}</span>
            <span
              v-if="provider.supported"
              class="rounded bg-[#14b8a6]/15 px-1.5 py-0.5 text-[10px] font-medium text-[#14b8a6]"
            >
              {{ t('home.providers.supported') }}
            </span>
            <span
              v-else
              class="rounded bg-white/5 px-1.5 py-0.5 text-[10px] font-medium text-gray-500"
            >
              {{ t('home.providers.soon') }}
            </span>
          </div>
        </div>
      </div>

      <!-- ─── CTA ──────────────────────────────────────────── -->
      <div ref="ctaRef" class="scroll-section mx-auto mb-16 max-w-3xl">
        <div
          class="scroll-card relative overflow-hidden rounded-3xl border border-[#14b8a6]/20 bg-gradient-to-br from-[#14b8a6]/10 via-[#0ea5e9]/5 to-[#8b5cf6]/5 p-12 text-center backdrop-blur-sm"
        >
          <div class="pointer-events-none absolute inset-0 overflow-hidden">
            <div class="cta-blob cta-blob-1"></div>
            <div class="cta-blob cta-blob-2"></div>
          </div>
          <div class="relative z-10">
            <h2 class="mb-3 text-3xl font-bold tracking-tight text-white md:text-4xl text-3d-hover">
              {{ t('home.cta.title') }}
            </h2>
            <p class="mx-auto mb-8 max-w-lg text-base text-gray-400">
              {{ t('home.cta.description') }}
            </p>
            <router-link
              v-if="!isAuthenticated"
              to="/register"
              class="group inline-flex items-center gap-2 rounded-full bg-[#14b8a6] px-8 py-3.5 text-base font-semibold text-[#0a0a0f] shadow-lg shadow-[#14b8a6]/30 transition-all hover:bg-[#2dd4bf] hover:shadow-[#14b8a6]/50 hover:-translate-y-0.5"
            >
              {{ t('home.cta.button') }}
              <svg class="h-4 w-4 transition-transform group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" />
              </svg>
            </router-link>
            <router-link
              v-else
              :to="dashboardPath"
              class="group inline-flex items-center gap-2 rounded-full bg-[#14b8a6] px-8 py-3.5 text-base font-semibold text-[#0a0a0f] shadow-lg shadow-[#14b8a6]/30 transition-all hover:bg-[#2dd4bf] hover:shadow-[#14b8a6]/50 hover:-translate-y-0.5"
            >
              {{ t('home.goToDashboard') }}
              <svg class="h-4 w-4 transition-transform group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" />
              </svg>
            </router-link>
          </div>
        </div>
      </div>

    </main>

    <!-- ── Footer ─────────────────────────────────────────── -->
    <footer class="relative z-10 border-t border-white/5 px-6 py-8">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left">
        <p class="text-sm text-gray-600">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-600 transition-colors hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, h, type Ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})
const currentYear = computed(() => new Date().getFullYear())

// ── Scroll tracking ─────────────────────────────────────────────
const pageRef = ref<HTMLElement>()
const scrollY = ref(0)
const heroTitleRef = ref<HTMLElement>()
const painPointsRef = ref<HTMLElement>()
const solutionsRef = ref<HTMLElement>()
const featuresRef = ref<HTMLElement>()
const comparisonRef = ref<HTMLElement>()
const providersRef = ref<HTMLElement>()
const ctaRef = ref<HTMLElement>()

function onScroll() {
  scrollY.value = window.scrollY
}

// ── Intersection Observer ────────────────────────────────────────
const visibleSections = reactive({
  heroTitle: false,
  painPoints: false,
  solutions: false,
  features: false,
  comparison: false,
  providers: false,
  cta: false,
})

let observer: IntersectionObserver | null = null

function setupObserver() {
  const options = {
    root: null,
    rootMargin: '0px 0px -80px 0px',
    threshold: 0.1,
  }

  const sections: { ref: Ref<HTMLElement | undefined>; key: keyof typeof visibleSections }[] = [
    { ref: heroTitleRef, key: 'heroTitle' },
    { ref: painPointsRef, key: 'painPoints' },
    { ref: solutionsRef, key: 'solutions' },
    { ref: featuresRef, key: 'features' },
    { ref: comparisonRef, key: 'comparison' },
    { ref: providersRef, key: 'providers' },
    { ref: ctaRef, key: 'cta' },
  ]

  observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const key = entry.target.getAttribute('data-section') as keyof typeof visibleSections
        if (key) visibleSections[key] = true
      }
    })
  }, options)

  sections.forEach(({ ref: sectionRef, key }) => {
    if (sectionRef.value) {
      sectionRef.value.setAttribute('data-section', key)
      observer!.observe(sectionRef.value)
    }
  })
}

// ── Pain Points ────────────────────────────────────────────────
const painPointItems = computed(() => [
  {
    title: t('home.painPoints.items.expensive.title'),
    desc: t('home.painPoints.items.expensive.desc'),
    icon: DollarIcon,
  },
  {
    title: t('home.painPoints.items.complex.title'),
    desc: t('home.painPoints.items.complex.desc'),
    icon: UsersIcon,
  },
  {
    title: t('home.painPoints.items.unstable.title'),
    desc: t('home.painPoints.items.unstable.desc'),
    icon: AlertIcon,
  },
  {
    title: t('home.painPoints.items.noControl.title'),
    desc: t('home.painPoints.items.noControl.desc'),
    icon: ChartIcon,
  },
])

// ── Solution Steps ────────────────────────────────────────────
const solutionSteps = computed(() => [
  {
    title: t('home.solutions.steps.s1.title'),
    desc: t('home.solutions.steps.s1.desc'),
  },
  {
    title: t('home.solutions.steps.s2.title'),
    desc: t('home.solutions.steps.s2.desc'),
  },
  {
    title: t('home.solutions.steps.s3.title'),
    desc: t('home.solutions.steps.s3.desc'),
  },
])

// ── Features ─────────────────────────────────────────────────
const featureItems = [
  {
    title: t('home.features.unifiedGateway'),
    desc: t('home.features.unifiedGatewayDesc'),
    icon: ServerIcon,
    iconBg: 'bg-gradient-to-br from-blue-500 to-blue-600 shadow-lg shadow-blue-500/30',
  },
  {
    title: t('home.features.multiAccount'),
    desc: t('home.features.multiAccountDesc'),
    icon: UsersIcon,
    iconBg: 'bg-gradient-to-br from-[#14b8a6] to-[#0d9488] shadow-lg shadow-[#14b8a6]/30',
  },
  {
    title: t('home.features.balanceQuota'),
    desc: t('home.features.balanceQuotaDesc'),
    icon: DollarIcon,
    iconBg: 'bg-gradient-to-br from-purple-500 to-purple-600 shadow-lg shadow-purple-500/30',
  },
]

// ── Comparison ────────────────────────────────────────────────
const comparisonRows = computed(() => [
  {
    feature: t('home.comparison.items.pricing.feature'),
    official: t('home.comparison.items.pricing.official'),
    us: t('home.comparison.items.pricing.us'),
  },
  {
    feature: t('home.comparison.items.models.feature'),
    official: t('home.comparison.items.models.official'),
    us: t('home.comparison.items.models.us'),
  },
  {
    feature: t('home.comparison.items.management.feature'),
    official: t('home.comparison.items.management.official'),
    us: t('home.comparison.items.management.us'),
  },
  {
    feature: t('home.comparison.items.stability.feature'),
    official: t('home.comparison.items.stability.official'),
    us: t('home.comparison.items.stability.us'),
  },
  {
    feature: t('home.comparison.items.control.feature'),
    official: t('home.comparison.items.control.official'),
    us: t('home.comparison.items.control.us'),
  },
])

// ── Providers ─────────────────────────────────────────────────
const providerItems = computed(() => [
  { name: 'Claude', initial: 'C', bg: 'bg-gradient-to-br from-orange-400 to-orange-500', supported: true },
  { name: 'GPT', initial: 'G', bg: 'bg-gradient-to-br from-green-500 to-green-600', supported: true },
  { name: 'Gemini', initial: 'G', bg: 'bg-gradient-to-br from-blue-500 to-blue-600', supported: true },
  { name: 'Antigravity', initial: 'A', bg: 'bg-gradient-to-br from-rose-500 to-pink-600', supported: true },
  { name: t('home.providers.more'), initial: '+', bg: 'bg-gradient-to-br from-gray-500 to-gray-600', supported: false },
])

// ── Inline SVG Icon Components ────────────────────────────────
function ServerIcon({ class: cls }: { class?: string }) {
  return h('svg', { class: cls, fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01' })
  ])
}
function UsersIcon({ class: cls }: { class?: string }) {
  return h('svg', { class: cls, fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z' })
  ])
}
function DollarIcon({ class: cls }: { class?: string }) {
  return h('svg', { class: cls, fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
  ])
}
function ChartIcon({ class: cls }: { class?: string }) {
  return h('svg', { class: cls, fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z' })
  ])
}
function AlertIcon({ class: cls }: { class?: string }) {
  return h('svg', { class: cls, fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z' })
  ])
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  setupObserver()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  observer?.disconnect()
})
</script>

<style scoped>
/* ════════════════════════════════════════════════
   HERO PAGE-LOAD ANIMATIONS
   ════════════════════════════════════════════════ */

/* Header: slide down from top */
.header-animate {
  animation: headerSlide 0.6s cubic-bezier(0.22, 1, 0.36, 1) both;
}
@keyframes headerSlide {
  from { opacity: 0; transform: translateY(-16px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* Hero tags: staggered fade + scale in */
.hero-tag {
  animation: tagPop 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}
.hero-tag-1 { animation-delay: 0.3s; }
.hero-tag-2 { animation-delay: 0.45s; }
.hero-tag-3 { animation-delay: 0.6s; }
@keyframes tagPop {
  from { opacity: 0; transform: scale(0.8) translateY(8px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

/* Tag pulsing dot */
.tag-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  animation: dotPulse 2s ease-in-out infinite;
}
@keyframes dotPulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50%       { opacity: 0.5; transform: scale(0.7); }
}

/* Title: 与下方 Pain Points 卡片使用完全相同的淡入上浮动画 */
.hero-title {
  animation: none; /* 由 .scroll-card / .is-visible 接管 */
}

/* Subtitle lines: staggered fade + up */
.hero-subtitle-1 {
  animation: fadeUp 0.6s ease-out 0.5s both;
}
.hero-subtitle-2 {
  animation: fadeUp 0.6s ease-out 0.65s both;
}
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(16px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* CTA buttons: fade + slide up */
.hero-cta {
  animation: fadeUp 0.6s ease-out 0.8s both;
}

/* ════════════════════════════════════════════════
   BACKGROUND BLOB ANIMATIONS
   ════════════════════════════════════════════════ */
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
  animation: blobFloat 12s ease-in-out infinite alternate;
}
.blob-1 {
  width: 500px;
  height: 500px;
  top: -128px;
  right: -128px;
  background: radial-gradient(circle, rgba(20,184,166,0.18) 0%, transparent 70%);
  animation-duration: 14s;
  animation-name: blobFloat1;
}
.blob-2 {
  width: 600px;
  height: 600px;
  bottom: -160px;
  left: -160px;
  background: radial-gradient(circle, rgba(14,165,233,0.14) 0%, transparent 70%);
  animation-duration: 18s;
  animation-name: blobFloat2;
}
.blob-3 {
  width: 400px;
  height: 400px;
  top: 33%;
  left: 50%;
  transform: translateX(-50%);
  background: radial-gradient(circle, rgba(139,92,246,0.12) 0%, transparent 70%);
  animation-duration: 22s;
  animation-name: blobFloat3;
}

@keyframes blobFloat1 {
  0%   { transform: translate(0, 0) scale(1); }
  33%  { transform: translate(-30px, 20px) scale(1.05); }
  66%  { transform: translate(20px, -15px) scale(0.97); }
  100% { transform: translate(-10px, 30px) scale(1.02); }
}
@keyframes blobFloat2 {
  0%   { transform: translate(0, 0) scale(1); }
  40%  { transform: translate(25px, -20px) scale(1.06); }
  80%  { transform: translate(-15px, 15px) scale(0.95); }
  100% { transform: translate(10px, -5px) scale(1.03); }
}
@keyframes blobFloat3 {
  0%   { transform: translateX(-50%) scale(1); }
  50%  { transform: translateX(-45%) scale(1.08); }
  100% { transform: translateX(-55%) scale(0.94); }
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(20,184,166,0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(20,184,166,0.04) 1px, transparent 1px);
  background-size: 80px 80px;
}

/* ════════════════════════════════════════════════
   SCROLL-TRIGGERED SECTION ANIMATIONS
   ════════════════════════════════════════════════ */
.scroll-card {
  opacity: 0;
  transform: translateY(32px);
  transition: opacity 0.6s cubic-bezier(0.22, 1, 0.36, 1),
              transform 0.6s cubic-bezier(0.22, 1, 0.36, 1);
}
.scroll-card.is-visible {
  opacity: 1;
  transform: translateY(0);
}

/* Comparison table: rows always visible (container animation handles entrance) */

/* CTA blob pulse */
.cta-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}
.cta-blob-1 {
  width: 240px;
  height: 240px;
  top: -80px;
  left: -80px;
  background: radial-gradient(circle, rgba(20,184,166,0.25) 0%, transparent 70%);
  animation: ctaPulse1 4s ease-in-out infinite;
}
.cta-blob-2 {
  width: 240px;
  height: 240px;
  bottom: -80px;
  right: -80px;
  background: radial-gradient(circle, rgba(14,165,233,0.25) 0%, transparent 70%);
  animation: ctaPulse2 5s ease-in-out infinite 1s;
}
@keyframes ctaPulse1 {
  0%, 100% { transform: scale(1); opacity: 0.6; }
  50%       { transform: scale(1.15); opacity: 1; }
}
@keyframes ctaPulse2 {
  0%, 100% { transform: scale(1.1); opacity: 0.5; }
  50%       { transform: scale(0.9); opacity: 0.8; }
}
</style>
