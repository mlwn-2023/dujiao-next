<template>
  <div class="minimal-shell flex min-h-screen flex-col bg-background text-foreground" :style="minimalThemeStyle">
    <header class="sticky top-0 z-50 border-b border-border/80 bg-background/95 backdrop-blur">
      <div class="mx-auto flex h-16 w-full max-w-6xl items-center gap-3 px-4 sm:px-6">
        <RouterLink to="/" class="min-w-0">
          <span class="flex h-9 max-w-[8.5rem] items-center rounded-full bg-[var(--minimal-site-icon-bg)] px-3 text-[var(--minimal-site-icon-text)] sm:max-w-xs sm:px-4">
            <span class="truncate text-sm font-bold tracking-tight sm:text-base">{{ brandName }}</span>
          </span>
        </RouterLink>

        <nav class="ml-auto flex items-center gap-1.5">
          <button
            type="button"
            data-testid="minimal-search-button"
            :data-search-mode="userAuthStore.isAuthenticated ? 'products' : 'guest-orders'"
            class="grid h-10 w-10 place-items-center rounded-full text-muted-foreground transition hover:bg-secondary hover:text-foreground"
            :aria-label="userAuthStore.isAuthenticated ? t('products.searchLabel') : t('navbar.guestOrders')"
            @click="handleSearchAction"
          >
            <Search class="h-4.5 w-4.5" />
          </button>
          <button v-if="showThemeSwitcher" type="button" data-testid="minimal-theme-button" class="grid h-10 w-10 place-items-center rounded-full text-muted-foreground transition hover:bg-secondary hover:text-foreground" :aria-label="t('resellerConsole.common.toggleTheme')" @click="toggleTheme">
            <Sun v-if="theme === 'dark'" class="h-4.5 w-4.5" />
            <Moon v-else class="h-4.5 w-4.5" />
          </button>
          <div v-if="showLanguageSwitcher" class="relative">
            <button type="button" data-testid="minimal-language-button" class="grid h-10 w-10 place-items-center rounded-full text-muted-foreground transition hover:bg-secondary hover:text-foreground" :aria-label="t('navbar.selectLanguage')" @click="languageOpen = !languageOpen">
              <Languages class="h-4.5 w-4.5" />
            </button>
            <div v-if="languageOpen" class="absolute right-0 top-[calc(100%+8px)] z-[60] w-40 rounded-xl border border-border bg-card p-2 shadow-xl">
              <button v-for="language in languages" :key="language.code" type="button" class="flex w-full items-center justify-between rounded-lg px-3 py-2 text-left text-sm transition hover:bg-secondary" :class="appStore.locale === language.code ? 'font-semibold text-primary' : 'text-muted-foreground'" @click="changeLanguage(language.code)">
                {{ language.name }}
                <span v-if="appStore.locale === language.code" class="h-1.5 w-1.5 rounded-full bg-primary" />
              </button>
            </div>
          </div>
          <RouterLink to="/cart" class="relative grid h-10 w-10 place-items-center rounded-full text-muted-foreground transition hover:bg-secondary hover:text-foreground" :aria-label="t('navbar.cart')">
            <ShoppingBag class="h-4.5 w-4.5" />
            <span v-if="cartCount" class="absolute -right-0.5 -top-0.5 grid h-5 min-w-5 place-items-center rounded-full bg-primary px-1 text-[10px] font-bold text-primary-foreground">{{ cartCount }}</span>
          </RouterLink>
          <RouterLink :to="userAuthStore.isAuthenticated ? '/me' : '/auth/login'" class="grid h-10 w-10 place-items-center rounded-full text-muted-foreground transition hover:bg-secondary hover:text-foreground" :aria-label="userAuthStore.isAuthenticated ? t('navbar.personalCenter') : t('navbar.login')">
            <UserRound class="h-4.5 w-4.5" />
          </RouterLink>
        </nav>
      </div>
      <form v-if="searchOpen" data-testid="minimal-search-form" class="border-t border-border/70 bg-background" @submit.prevent="submitSearch">
        <div class="mx-auto flex w-full max-w-6xl items-center gap-2 px-4 py-3 sm:px-6">
          <Search class="h-4 w-4 shrink-0 text-muted-foreground" />
          <input ref="searchInputRef" v-model="searchText" type="search" class="h-10 min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground" :placeholder="t('products.searchBoxPlaceholder')" />
          <button v-if="searchText" type="button" class="grid h-8 w-8 place-items-center rounded-full text-muted-foreground hover:bg-secondary" :aria-label="t('blog.searchClear')" @click="clearSearch">
            <X class="h-4 w-4" />
          </button>
          <button type="submit" class="rounded-full bg-[var(--minimal-button-bg)] px-4 py-2 text-sm font-semibold text-[var(--minimal-button-text)]">{{ t('products.searchLabel') }}</button>
        </div>
      </form>
    </header>

    <main class="flex-1">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { Languages, Moon, Search, ShoppingBag, Sun, UserRound, X } from 'lucide-vue-next'
import { useAppStore } from '../../../stores/app'
import { useCartStore } from '../../../stores/cart'
import { useUserAuthStore } from '../../../stores/userAuth'
import { useStorefrontControls } from '../../../composables/useStorefrontControls'
import { useTheme } from '../../../utils/theme'
import '../styles/minimal.css'

const { t } = useI18n()
const appStore = useAppStore()
const cartStore = useCartStore()
const userAuthStore = useUserAuthStore()
const route = useRoute()
const router = useRouter()
const { theme, toggleTheme } = useTheme()
const { showLanguageSwitcher, showThemeSwitcher } = useStorefrontControls()
const searchOpen = ref(userAuthStore.isAuthenticated && Boolean(route.query.search))
const searchText = ref(String(route.query.search || ''))
const searchInputRef = ref<HTMLInputElement | null>(null)
const languageOpen = ref(false)

const languages = [
  { code: 'zh-CN', name: '简体中文' },
  { code: 'zh-TW', name: '繁體中文' },
  { code: 'en-US', name: 'English' },
]

const defaultMinimalColors = {
  announcement_light: '#2563EB',
  announcement_dark: '#1E3A8A',
  button_light: '#2563EB',
  button_dark: '#60A5FA',
}

const normalizeHexColor = (value: unknown, fallback: string) => {
  const color = String(value || '').trim().toUpperCase()
  return /^#[0-9A-F]{6}$/.test(color) ? color : fallback
}

const contrastTextColor = (hex: string) => {
  const toLinearChannel = (offset: number) => {
    const channel = Number.parseInt(hex.slice(offset, offset + 2), 16) / 255
    return channel <= 0.03928
      ? channel / 12.92
      : ((channel + 0.055) / 1.055) ** 2.4
  }
  const red = toLinearChannel(1)
  const green = toLinearChannel(3)
  const blue = toLinearChannel(5)
  const luminance = 0.2126 * red + 0.7152 * green + 0.0722 * blue
  return luminance > 0.42 ? '#0B1120' : '#FFFFFF'
}

const minimalThemeStyle = computed<Record<string, string>>(() => {
  const configured = appStore.config?.minimal_theme_colors || {}
  const dark = theme.value === 'dark'
  const announcement = normalizeHexColor(
    dark ? configured.announcement_dark : configured.announcement_light,
    dark ? defaultMinimalColors.announcement_dark : defaultMinimalColors.announcement_light,
  )
  const button = normalizeHexColor(
    dark ? configured.button_dark : configured.button_light,
    dark ? defaultMinimalColors.button_dark : defaultMinimalColors.button_light,
  )
  return {
    '--minimal-announcement-bg': announcement,
    '--minimal-announcement-text': contrastTextColor(announcement),
    '--minimal-button-bg': button,
    '--minimal-button-text': contrastTextColor(button),
    '--minimal-site-icon-bg': button,
    '--minimal-site-icon-text': contrastTextColor(button),
  }
})

const changeLanguage = (code: string) => {
  appStore.setLocale(code)
  languageOpen.value = false
}

const handleSearchAction = async () => {
  if (!userAuthStore.isAuthenticated) {
    searchOpen.value = false
    await router.push('/guest/orders')
    return
  }
  searchOpen.value = !searchOpen.value
}

const submitSearch = async () => {
  const search = searchText.value.trim()
  await router.push({ path: '/products', query: search ? { search } : {} })
  searchOpen.value = Boolean(search)
}

const clearSearch = async () => {
  searchText.value = ''
  await router.push({ path: '/products' })
  await nextTick()
  searchInputRef.value?.focus()
}

watch(searchOpen, async (open) => {
  if (open) {
    await nextTick()
    searchInputRef.value?.focus()
  }
})

watch(() => route.query.search, (value) => {
  searchText.value = String(value || '')
})

watch(() => userAuthStore.isAuthenticated, (authenticated) => {
  if (!authenticated) searchOpen.value = false
})

const brandName = computed(() => String(appStore.config?.brand?.site_name || '').trim() || 'Dujiao-Next')
const cartCount = computed(() => cartStore.totalItems)
</script>
