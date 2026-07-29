<template>
  <RouterLink
    :to="`/products/${product.slug}`"
    :data-layout-mode="mode"
    class="minimal-product-card group relative overflow-hidden rounded-2xl border border-border bg-card transition hover:-translate-y-0.5 hover:border-foreground/20 hover:shadow-lg"
    :class="[
      mode === 'list' ? 'flex flex-row items-center' : 'flex min-h-[260px] flex-col',
      { 'opacity-60': soldOut },
    ]"
  >
    <div
      data-testid="minimal-product-cover"
      class="shrink-0 overflow-hidden bg-muted"
      :class="mode === 'list' ? 'w-16 self-stretch sm:w-[72px]' : 'aspect-[4/3] w-full'"
    >
      <img
        v-if="coverImage && !imageErrored"
        :src="coverImage"
        :alt="title"
        loading="lazy"
        class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
        @error="imageErrored = true"
      />
      <div v-else class="grid h-full w-full place-items-center text-muted-foreground">
        <ImageIcon :class="mode === 'list' ? 'h-5 w-5' : 'h-7 w-7'" :stroke-width="1.5" />
      </div>
    </div>
    <div v-if="mode === 'list'" class="flex min-w-0 flex-1 items-center gap-2.5 p-2.5 sm:gap-3 sm:p-3">
      <div class="min-w-0 flex-1">
        <h3 class="truncate text-sm font-semibold leading-tight tracking-tight sm:text-base">{{ title }}</h3>
        <div class="mt-1 flex min-w-0 items-center gap-2">
          <span class="truncate text-sm font-bold text-primary sm:text-base">{{ displayPrice }}</span>
          <span class="shrink-0 text-[11px] text-muted-foreground">{{ stockLabel }}</span>
        </div>
      </div>
      <button
        type="button"
        class="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-[var(--minimal-button-bg)] text-[var(--minimal-button-text)] transition group-hover:scale-105 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="soldOut"
        :aria-label="t('products.quickBuyAria')"
        @click.prevent.stop="$emit('quickBuy', product)"
      >
        <ArrowUpRight class="h-3.5 w-3.5" />
      </button>
    </div>
    <div v-else class="flex min-w-0 flex-1 flex-col p-3 sm:p-4">
      <span class="line-clamp-1 text-xs font-medium text-muted-foreground">{{ categoryName || t('products.categoryLabel') }}</span>
      <h3 class="mt-1.5 line-clamp-2 text-[15px] font-bold leading-snug tracking-tight sm:text-base">{{ title }}</h3>
      <p v-if="description" class="mt-1.5 line-clamp-2 text-xs leading-relaxed text-muted-foreground">{{ description }}</p>
      <div class="mt-auto flex items-end justify-between gap-2 pt-3">
        <div class="min-w-0">
          <span class="block truncate text-base font-bold text-primary sm:text-lg">{{ displayPrice }}</span>
          <span class="mt-0.5 block text-[11px] text-muted-foreground">{{ stockLabel }}</span>
        </div>
        <button
          type="button"
          class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-[var(--minimal-button-bg)] text-[var(--minimal-button-text)] transition group-hover:scale-105 disabled:cursor-not-allowed disabled:opacity-40"
          :disabled="soldOut"
          :aria-label="t('products.quickBuyAria')"
          @click.prevent.stop="$emit('quickBuy', product)"
        >
          <ArrowUpRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowUpRight, Image as ImageIcon } from 'lucide-vue-next'
import { useLocalized, useProductLabels } from '../../../composables/useProduct'
import { getFirstImageUrl } from '../../../utils/image'

const props = withDefaults(defineProps<{ product: any; mode?: 'card' | 'list' }>(), {
  mode: 'card',
})
defineEmits<{ quickBuy: [product: any] }>()

const { t } = useI18n()
const { getLocalizedText, siteCurrency, formatPrice } = useLocalized()
const { getStockStatusLabel, isSoldOut, hasPromotionPrice, getPromotionPriceAmount } = useProductLabels()

const title = computed(() => getLocalizedText(props.product?.title))
const description = computed(() => getLocalizedText(props.product?.description))
const categoryName = computed(() => getLocalizedText(props.product?.category?.name))
const coverImage = computed(() => getFirstImageUrl(props.product?.images))
const imageErrored = ref(false)
const soldOut = computed(() => isSoldOut(props.product))
const stockLabel = computed(() => getStockStatusLabel(props.product))
const displayPrice = computed(() => formatPrice(
  hasPromotionPrice(props.product) ? getPromotionPriceAmount(props.product) : props.product?.price_amount,
  siteCurrency.value,
))

watch(coverImage, () => {
  imageErrored.value = false
})
</script>
