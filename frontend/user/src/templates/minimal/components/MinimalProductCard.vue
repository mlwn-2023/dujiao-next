<template>
  <RouterLink
    :to="`/products/${product.slug}`"
    class="minimal-product-card group relative flex min-h-[164px] flex-col overflow-hidden rounded-2xl border border-border bg-card p-4 transition hover:-translate-y-0.5 hover:border-foreground/20 hover:shadow-lg"
    :class="{ 'opacity-60': soldOut }"
  >
    <div class="flex items-start justify-between gap-2">
      <span class="line-clamp-1 text-xs font-medium text-muted-foreground">{{ categoryName || t('products.categoryLabel') }}</span>
      <img v-if="categoryIcon" :src="categoryIcon" :alt="categoryName" class="h-7 w-7 shrink-0 rounded-lg object-contain opacity-80" />
    </div>
    <h3 class="mt-2 line-clamp-2 text-[15px] font-bold leading-snug tracking-tight sm:text-base">{{ title }}</h3>
    <p v-if="description" class="mt-1.5 line-clamp-2 text-xs leading-relaxed text-muted-foreground">{{ description }}</p>
    <div class="mt-auto flex items-end justify-between gap-2 pt-4">
      <div class="min-w-0">
        <span class="block truncate text-base font-bold text-primary sm:text-lg">{{ displayPrice }}</span>
        <span class="mt-0.5 block text-[11px] text-muted-foreground">{{ stockLabel }}</span>
      </div>
      <button
        type="button"
        class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-foreground text-background transition group-hover:scale-105 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="soldOut"
        :aria-label="t('products.quickBuyAria')"
        @click.prevent.stop="$emit('quickBuy', product)"
      >
        <ArrowUpRight class="h-4 w-4" />
      </button>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowUpRight } from 'lucide-vue-next'
import { useLocalized, useProductLabels } from '../../../composables/useProduct'
import { getImageUrl } from '../../../utils/image'

const props = defineProps<{ product: any }>()
defineEmits<{ quickBuy: [product: any] }>()

const { t } = useI18n()
const { getLocalizedText, siteCurrency, formatPrice } = useLocalized()
const { getStockStatusLabel, isSoldOut, hasPromotionPrice, getPromotionPriceAmount } = useProductLabels()

const title = computed(() => getLocalizedText(props.product?.title))
const description = computed(() => getLocalizedText(props.product?.description))
const categoryName = computed(() => getLocalizedText(props.product?.category?.name))
const categoryIcon = computed(() => {
  const raw = String(props.product?.category?.icon || '').trim()
  return raw ? getImageUrl(raw) : ''
})
const soldOut = computed(() => isSoldOut(props.product))
const stockLabel = computed(() => getStockStatusLabel(props.product))
const displayPrice = computed(() => formatPrice(
  hasPromotionPrice(props.product) ? getPromotionPriceAmount(props.product) : props.product?.price_amount,
  siteCurrency.value,
))
</script>
