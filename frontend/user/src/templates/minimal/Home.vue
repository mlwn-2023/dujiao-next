<template>
  <div class="minimal-home">
    <section v-if="announcement" class="border-b border-border bg-foreground text-background">
      <div class="mx-auto w-full max-w-6xl px-4 py-4 sm:px-6">
        <div class="flex items-start gap-3">
          <Megaphone class="mt-0.5 h-5 w-5 shrink-0 opacity-80" />
          <div class="min-w-0">
            <h2 v-if="announcementTitle" class="text-sm font-bold">{{ announcementTitle }}</h2>
            <div class="minimal-announcement prose prose-sm mt-1 max-w-none text-current opacity-90" v-html="announcementContent" />
          </div>
        </div>
      </div>
    </section>

    <section class="border-b border-border">
      <div class="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6 sm:py-12">
        <div class="max-w-2xl">
          <p class="text-xs font-semibold uppercase tracking-[0.18em] text-primary">{{ t('minimal.storeLabel') }}</p>
          <h1 class="mt-2 text-3xl font-bold tracking-[-0.04em] sm:text-5xl">{{ brandName }}</h1>
          <p v-if="brandDescription" class="mt-3 max-w-xl text-sm leading-relaxed text-muted-foreground sm:text-base">{{ brandDescription }}</p>
        </div>
      </div>
    </section>

    <section class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 sm:py-8">
      <div class="mb-5 flex gap-2 overflow-x-auto pb-1 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        <button
          type="button"
          class="shrink-0 rounded-full border px-4 py-2 text-sm font-medium transition"
          :class="selectedCategory === null ? 'border-foreground bg-foreground text-background' : 'border-border bg-card text-muted-foreground hover:text-foreground'"
          @click="selectCategory(null)"
        >
          {{ t('products.allCategories') }}
        </button>
        <button
          v-for="category in categories"
          :key="category.id"
          type="button"
          class="shrink-0 rounded-full border px-4 py-2 text-sm font-medium transition"
          :class="selectedCategory === category.id ? 'border-foreground bg-foreground text-background' : 'border-border bg-card text-muted-foreground hover:text-foreground'"
          @click="selectCategory(category.id)"
        >
          {{ getLocalizedText(category.name) }}
        </button>
      </div>

      <div class="mb-4 flex items-end justify-between gap-3">
        <div>
          <h2 class="text-xl font-bold tracking-tight sm:text-2xl">{{ selectedCategoryName || t('home.featured.title') }}</h2>
          <p class="mt-1 text-xs text-muted-foreground">{{ t('minimal.productCount', { count: products.length }) }}</p>
        </div>
      </div>

      <div v-if="loading" class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        <div v-for="index in 8" :key="index" class="h-[164px] animate-pulse rounded-2xl border border-border bg-muted/60" />
      </div>
      <div v-else-if="products.length" class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        <MinimalProductCard v-for="product in products" :key="product.id" :product="product" @quick-buy="openQuickBuy" />
      </div>
      <div v-else class="rounded-2xl border border-dashed border-border px-5 py-16 text-center text-sm text-muted-foreground">
        {{ t('home.featured.empty') }}
      </div>
    </section>

    <ProductQuickBuy
      v-if="quickBuyProduct"
      :product="quickBuyProduct"
      :visible="quickBuyVisible"
      @update:visible="quickBuyVisible = $event"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import DOMPurify from 'dompurify'
import { Megaphone } from 'lucide-vue-next'
import { categoryAPI, productAPI } from '../../api'
import { useLocalized } from '../../composables/useProduct'
import { usePageSeo } from '../../composables/usePageSeo'
import { useAppStore } from '../../stores/app'
import { processHtmlForDisplay } from '../../utils/content'
import { type PublicCategory } from '../../utils/category'
import ProductQuickBuy from '../../components/ProductQuickBuy.vue'
import MinimalProductCard from './components/MinimalProductCard.vue'

const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const { getLocalizedText } = useLocalized()

const products = ref<any[]>([])
const categories = ref<PublicCategory[]>([])
const loading = ref(true)
const selectedCategory = ref<number | null>(null)
const quickBuyProduct = ref<any>(null)
const quickBuyVisible = ref(false)

const brandName = computed(() => String(appStore.config?.brand?.site_name || '').trim() || 'Dujiao-Next')
const brandDescription = computed(() => getLocalizedText(appStore.config?.brand?.site_description))
const announcement = computed(() => appStore.config?.announcement || null)
const announcementTitle = computed(() => getLocalizedText(announcement.value?.title))
const announcementContent = computed(() => DOMPurify.sanitize(
  processHtmlForDisplay(getLocalizedText(announcement.value?.content)),
  {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 's', 'code', 'a', 'ul', 'ol', 'li', 'span'],
    ALLOWED_ATTR: ['href', 'target', 'rel'],
  },
))
const selectedCategoryName = computed(() => {
  const category = categories.value.find((item) => item.id === selectedCategory.value)
  return category ? getLocalizedText(category.name) : ''
})

const loadProducts = async () => {
  loading.value = true
  try {
    const params: Record<string, unknown> = { page: 1, page_size: 40 }
    if (selectedCategory.value) params.category_id = selectedCategory.value
    const response = await productAPI.list(params)
    products.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load minimal products:', error)
    products.value = []
  } finally {
    loading.value = false
  }
}

const selectCategory = async (categoryId: number | null) => {
  selectedCategory.value = categoryId
  await loadProducts()
}

const openQuickBuy = (product: any) => {
  quickBuyProduct.value = product
  quickBuyVisible.value = true
}

usePageSeo({ canonicalPath: () => route.path })

onMounted(async () => {
  await appStore.loadConfig()
  try {
    const response = await categoryAPI.list()
    categories.value = response.data.data || []
    const slug = String(route.params.slug || '')
    if (slug) {
      selectedCategory.value = categories.value.find((item) => item.slug === slug)?.id || null
    }
  } catch (error) {
    console.error('Failed to load minimal categories:', error)
  }
  await loadProducts()
})
</script>
