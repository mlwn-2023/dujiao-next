<template>
  <div class="editorial-home min-h-screen bg-[#fffaf7] text-slate-900">
    <section class="mx-auto grid max-w-6xl gap-4 px-4 pb-10 pt-6 md:grid-cols-[1.35fr_.65fr] md:gap-6 md:pt-16">
      <div class="relative min-h-[300px] overflow-hidden rounded-2xl bg-rose-900 p-6 text-white sm:min-h-[360px] md:rounded-[2rem] md:p-12">
        <img v-if="heroImage" :src="heroImage" class="absolute inset-0 h-full w-full object-cover opacity-45" alt="" />
        <div class="relative flex h-full flex-col justify-end">
          <p class="mb-3 text-sm font-semibold uppercase tracking-[.22em] text-rose-200">{{ heroTitle || t('home.featured.title') }}</p>
          <h1 class="max-w-xl text-3xl font-black leading-tight sm:text-4xl md:text-6xl">{{ heroSubtitle || t('home.featured.description') }}</h1>
          <RouterLink to="/products" class="mt-7 inline-flex h-11 w-fit items-center rounded-full bg-white px-6 text-sm font-bold text-rose-900">{{ t('home.featured.viewAll') }}</RouterLink>
        </div>
      </div>
      <aside class="flex flex-col justify-between rounded-2xl border border-rose-200 bg-white p-6 md:rounded-[2rem] md:p-9">
        <div><span class="text-xs font-bold uppercase tracking-[.2em] text-rose-700">{{ t('nav.notice') }}</span><h2 class="mt-4 text-2xl font-black">{{ announcementTitle || t('home.latest.title') }}</h2><p class="mt-3 text-sm leading-7 text-slate-600">{{ announcementText || t('home.latest.description') }}</p></div>
        <RouterLink to="/notice" class="mt-8 text-sm font-bold text-rose-700">{{ t('blog.readMore') }} -></RouterLink>
      </aside>
    </section>
    <section class="mx-auto max-w-6xl px-4 pb-14"><div class="mb-6 flex items-end justify-between gap-4"><div><p class="text-xs font-bold uppercase tracking-[.2em] text-rose-700">01 / {{ t('home.featured.title') }}</p><h2 class="mt-2 text-2xl font-black sm:text-3xl">{{ t('home.featured.title') }}</h2></div><RouterLink to="/products" class="shrink-0 text-sm font-bold text-rose-700">{{ t('home.featured.viewAll') }}</RouterLink></div><TemplateProductCollection :products="products" :mode="templateMode" grid-class="grid grid-cols-1 gap-4 sm:grid-cols-2" @click="go" @quick-buy="quickBuy" /></section>
    <ProductQuickBuy v-if="selected" :product="selected" :visible="quickVisible" @update:visible="quickVisible = $event" />
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { productAPI } from '../../api'
import { useBannerCarousel } from '../../composables/useBannerCarousel'
import { useAppStore } from '../../stores/app'
import { useLocalized } from '../../composables/useProduct'
import TemplateProductCollection from '../components/TemplateProductCollection.vue'
import ProductQuickBuy from '../../components/ProductQuickBuy.vue'
const { t } = useI18n(); const router = useRouter(); const app = useAppStore(); const { getLocalizedText } = useLocalized(); const products = ref<any[]>([]); const selected = ref<any>(null); const quickVisible = ref(false)
const { heroImage, heroTitle, heroSubtitle, loadBanners } = useBannerCarousel(); const announcementTitle = computed(() => getLocalizedText(app.config?.announcement?.title)); const announcementText = computed(() => getLocalizedText(app.config?.announcement?.content)?.replace(/<[^>]+>/g, '').slice(0, 180))
const templateMode = computed<'card' | 'list'>(() => app.config?.template_mode === 'list' ? 'list' : 'card')
const go = (slug: string) => router.push(`/products/${slug}`); const quickBuy = (p: any) => { selected.value = p; quickVisible.value = true }
onMounted(async () => { await app.loadConfig(); await Promise.all([loadBanners(), productAPI.list({ page: 1, page_size: 6 }).then(r => { products.value = r.data.data || [] })]) })
</script>
