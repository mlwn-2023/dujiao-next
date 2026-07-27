<template>
  <div class="minimal-shell flex min-h-screen flex-col bg-background text-foreground">
    <header class="sticky top-0 z-50 border-b border-border/80 bg-background/95 backdrop-blur">
      <div class="mx-auto flex h-16 w-full max-w-6xl items-center gap-3 px-4 sm:px-6">
        <RouterLink to="/" class="flex min-w-0 items-center gap-2.5">
          <img v-if="brandIcon" :src="brandIcon" :alt="brandName" class="h-9 w-9 rounded-xl object-cover" />
          <span v-else class="grid h-9 w-9 place-items-center rounded-xl bg-foreground text-sm font-bold text-background">{{ brandInitial }}</span>
          <span class="hidden max-w-xs truncate text-base font-bold tracking-tight sm:inline sm:text-lg">{{ brandName }}</span>
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
          <RouterLink to="/guest/orders" class="hidden items-center gap-1.5 rounded-full border border-border px-3 py-2 text-sm font-medium transition hover:bg-secondary sm:inline-flex">
            <ClipboardList class="h-4 w-4" />
            {{ t('navbar.guestOrders') }}
          </RouterLink>
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
          <button type="submit" class="rounded-full bg-foreground px-4 py-2 text-sm font-semibold text-background">{{ t('products.searchLabel') }}</button>
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
import { ClipboardList, Languages, Moon, Search, ShoppingBag, Sun, UserRound, X } from 'lucide-vue-next'
import { useAppStore } from '../../../stores/app'
import { useCartStore } from '../../../stores/cart'
import { useUserAuthStore } from '../../../stores/userAuth'
import { useStorefrontControls } from '../../../composables/useStorefrontControls'
import { getImageUrl } from '../../../utils/image'
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
const brandInitial = computed(() => brandName.value.charAt(0).toUpperCase())
const brandIcon = computed(() => {
  const raw = String(appStore.config?.brand?.site_icon || '').trim()
  return raw ? getImageUrl(raw) : ''
})
const cartCount = computed(() => cartStore.totalItems)
</script>
