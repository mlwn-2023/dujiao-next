<template>
  <div class="minimal-home">
    <section v-if="announcement" class="mx-auto w-full max-w-6xl px-4 pt-4 sm:px-6 sm:pt-6">
      <div class="rounded-2xl border border-border bg-[var(--minimal-announcement-bg)] px-4 py-4 text-[var(--minimal-announcement-text)] sm:px-6">
        <div class="flex items-start gap-3">
          <Megaphone class="mt-0.5 h-5 w-5 shrink-0 opacity-80" />
          <div class="min-w-0">
            <h2 v-if="announcementTitle" class="text-sm font-bold">{{ announcementTitle }}</h2>
            <div class="minimal-announcement prose prose-sm mt-1 max-w-none text-current opacity-90" v-html="announcementContent" />
          </div>
        </div>
      </div>
    </section>

    <section class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 sm:py-8">
      <div class="mb-5 flex flex-wrap gap-2">
        <button
          type="button"
          class="shrink-0 rounded-full border px-4 py-2 text-sm font-medium transition"
          :class="selectedCategory === null ? 'border-[var(--minimal-button-bg)] bg-[var(--minimal-button-bg)] text-[var(--minimal-button-text)]' : 'border-border bg-card text-muted-foreground hover:text-foreground'"
          @click="selectCategory(null)"
        >
          {{ t('products.allCategories') }}
        </button>
        <button
          v-for="category in categories"
          :key="category.id"
          type="button"
          class="shrink-0 rounded-full border px-4 py-2 text-sm font-medium transition"
          :class="selectedCategory === category.id ? 'border-[var(--minimal-button-bg)] bg-[var(--minimal-button-bg)] text-[var(--minimal-button-text)]' : 'border-border bg-card text-muted-foreground hover:text-foreground'"
          @click="selectCategory(category.id)"
        >
          {{ getLocalizedText(category.name) }}
        </button>
      </div>

      <p data-testid="minimal-product-count" class="mb-4 text-xs text-muted-foreground">
        {{ t('minimal.productCount', { count: products.length }) }}
      </p>

      <div
        v-if="loading"
        :class="templateMode === 'list' ? 'space-y-3' : 'grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4'"
      >
        <div
          v-for="index in 8"
          :key="index"
          class="animate-pulse rounded-2xl border border-border bg-muted/60"
          :class="templateMode === 'list' ? 'h-16 sm:h-[72px]' : 'h-[260px]'"
        />
      </div>
      <div
        v-else-if="products.length"
        :class="templateMode === 'list' ? 'space-y-3' : 'grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4'"
      >
        <MinimalProductCard
          v-for="product in products"
          :key="product.id"
          :product="product"
          :mode="templateMode"
          @quick-buy="openQuickBuy"
        />
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
import { computed, onMounted, ref, watch } from 'vue'
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

const announcement = computed(() => appStore.config?.announcement || null)
const announcementTitle = computed(() => getLocalizedText(announcement.value?.title))
const announcementContent = computed(() => DOMPurify.sanitize(
  processHtmlForDisplay(getLocalizedText(announcement.value?.content)),
  {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 's', 'code', 'a', 'ul', 'ol', 'li', 'span'],
    ALLOWED_ATTR: ['href', 'target', 'rel'],
  },
))
const searchQuery = computed(() => String(route.query.search || '').trim())
const templateMode = computed<'card' | 'list'>(() => appStore.config?.template_mode === 'list' ? 'list' : 'card')

const loadProducts = async () => {
  loading.value = true
  try {
    const params: Record<string, unknown> = { page: 1, page_size: 40 }
    if (selectedCategory.value) params.category_id = selectedCategory.value
    if (searchQuery.value) params.search = searchQuery.value
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

watch(() => route.query.search, () => {
  void loadProducts()
})
</script>
