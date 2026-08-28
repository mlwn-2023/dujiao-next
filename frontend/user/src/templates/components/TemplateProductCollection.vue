<template>
  <div :class="mode === 'list' ? 'space-y-3' : gridClass">
    <component
      :is="mode === 'list' ? ProductListItem : ProductCard"
      v-for="(product, index) in products"
      :key="product.id"
      :product="product"
      :index="index"
      :animation-step="50"
      @click="emit('click', product.slug)"
      @quick-buy="emit('quick-buy', product)"
    />
  </div>
</template>

<script setup lang="ts">
import ProductCard from '../../components/ProductCard.vue'
import ProductListItem from '../../components/ProductListItem.vue'

withDefaults(defineProps<{
  products: any[]
  mode: 'card' | 'list'
  gridClass?: string
}>(), {
  gridClass: 'grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3',
})

const emit = defineEmits<{
  click: [slug: string]
  'quick-buy': [product: any]
}>()
</script>
