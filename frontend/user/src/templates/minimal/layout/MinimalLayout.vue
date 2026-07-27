<template>
  <div class="minimal-shell flex min-h-screen flex-col bg-background text-foreground">
    <header class="sticky top-0 z-50 border-b border-border/80 bg-background/95 backdrop-blur">
      <div class="mx-auto flex h-16 w-full max-w-6xl items-center gap-3 px-4 sm:px-6">
        <RouterLink to="/" class="flex min-w-0 items-center gap-2.5">
          <img v-if="brandIcon" :src="brandIcon" :alt="brandName" class="h-9 w-9 rounded-xl object-cover" />
          <span v-else class="grid h-9 w-9 place-items-center rounded-xl bg-foreground text-sm font-bold text-background">{{ brandInitial }}</span>
          <span class="max-w-[42vw] truncate text-base font-bold tracking-tight sm:max-w-xs sm:text-lg">{{ brandName }}</span>
        </RouterLink>

        <nav class="ml-auto flex items-center gap-1.5">
          <RouterLink to="/products" class="grid h-10 w-10 place-items-center rounded-full text-muted-foreground transition hover:bg-secondary hover:text-foreground" :aria-label="t('nav.products')">
            <Search class="h-4.5 w-4.5" />
          </RouterLink>
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
    </header>

    <main class="flex-1">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ClipboardList, Search, ShoppingBag, UserRound } from 'lucide-vue-next'
import { useAppStore } from '../../../stores/app'
import { useCartStore } from '../../../stores/cart'
import { useUserAuthStore } from '../../../stores/userAuth'
import { getImageUrl } from '../../../utils/image'
import '../styles/minimal.css'

const { t } = useI18n()
const appStore = useAppStore()
const cartStore = useCartStore()
const userAuthStore = useUserAuthStore()

const brandName = computed(() => String(appStore.config?.brand?.site_name || '').trim() || 'Dujiao-Next')
const brandInitial = computed(() => brandName.value.charAt(0).toUpperCase())
const brandIcon = computed(() => {
  const raw = String(appStore.config?.brand?.site_icon || '').trim()
  return raw ? getImageUrl(raw) : ''
})
const cartCount = computed(() => cartStore.totalItems)
</script>
