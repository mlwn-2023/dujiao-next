<template>
  <section
    v-if="showBanner && showHeroSection"
    data-testid="minimal-banner"
    class="mx-auto w-full max-w-6xl px-4 pt-4 sm:px-6 sm:pt-6"
  >
    <div
      class="relative overflow-hidden rounded-2xl border border-border bg-card"
      @touchstart="onBannerTouchStart"
      @touchend="onBannerTouchEnd"
    >
      <Transition name="minimal-banner-fade" mode="out-in">
        <img
          v-if="!bannerLoading && heroImage"
          :key="heroImage"
          :src="heroImage"
          :alt="heroTitle"
          class="absolute inset-0 h-full w-full object-cover"
        />
      </Transition>
      <div class="absolute inset-0 bg-gradient-to-r from-black/65 via-black/35 to-black/15" />

      <div
        v-if="bannerLoading"
        class="relative flex min-h-32 flex-col justify-end gap-2.5 p-4 sm:min-h-40 sm:p-5 md:min-h-[200px] md:p-6"
      >
        <div class="h-5 w-1/2 max-w-80 animate-pulse rounded bg-white/30" />
        <div class="h-3.5 w-2/3 max-w-96 animate-pulse rounded bg-white/25" />
      </div>

      <div
        v-else
        class="relative flex min-h-32 flex-col justify-between gap-3 p-4 sm:min-h-40 sm:p-5 md:min-h-[200px] md:p-6"
      >
        <div v-if="bannerCount > 1" class="flex items-center justify-end gap-1.5">
          <button
            type="button"
            class="grid h-8 w-8 place-items-center rounded-full border border-white/30 bg-black/25 text-white transition hover:bg-black/45"
            :aria-label="t('common.previousBanner')"
            @click="handlePrevHeroBanner"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>
          <button
            type="button"
            class="grid h-8 w-8 place-items-center rounded-full border border-white/30 bg-black/25 text-white transition hover:bg-black/45"
            :aria-label="t('common.nextBanner')"
            @click="handleNextHeroBanner"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>

        <div class="mt-auto min-w-0 max-w-2xl">
          <h2 class="truncate text-lg font-bold tracking-tight text-white sm:text-xl md:text-2xl">{{ heroTitle }}</h2>
          <p class="mt-1 line-clamp-1 text-xs text-white/80 sm:text-sm">{{ heroSubtitle }}</p>
        </div>

        <div v-if="bannerCount > 1 || hasHeroLink" class="flex items-center justify-between gap-3">
          <div v-if="bannerCount > 1" class="flex items-center gap-1.5">
            <button
              v-for="(_, index) in banners"
              :key="`minimal-banner-dot-${index}`"
              type="button"
              class="h-1.5 rounded-full transition-all"
              :class="index === currentBannerIndex ? 'w-5 bg-white' : 'w-1.5 bg-white/45 hover:bg-white/70'"
              :aria-label="t('common.switchBanner', { n: index + 1 })"
              @click="selectHeroBanner(index)"
            />
          </div>
          <button
            v-if="hasHeroLink"
            type="button"
            class="ml-auto inline-flex h-8 items-center gap-1.5 rounded-full bg-white px-3 text-xs font-semibold text-gray-900 transition hover:scale-[1.03]"
            @click="goToHeroLink"
          >
            {{ heroPrimaryButtonText }}
            <ArrowRight class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowRight, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { useBannerCarousel } from '../../../composables/useBannerCarousel'
import { useStorefrontControls } from '../../../composables/useStorefrontControls'

const { t } = useI18n()
const { showBanner } = useStorefrontControls()
const {
  banners,
  bannerLoading,
  currentBannerIndex,
  bannerCount,
  showHeroSection,
  heroImage,
  heroTitle,
  heroSubtitle,
  hasHeroLink,
  heroPrimaryButtonText,
  loadBanners,
  handleNextHeroBanner,
  handlePrevHeroBanner,
  selectHeroBanner,
  goToHeroLink,
  onBannerTouchStart,
  onBannerTouchEnd,
  stopHeroAutoPlay,
} = useBannerCarousel()

onMounted(() => { void loadBanners() })
onUnmounted(() => stopHeroAutoPlay())
</script>

<style scoped>
.minimal-banner-fade-enter-active,
.minimal-banner-fade-leave-active {
  transition: opacity 300ms ease;
}

.minimal-banner-fade-enter-from,
.minimal-banner-fade-leave-to {
  opacity: 0;
}
</style>
