<template>
  <div v-if="hasContact" class="fixed bottom-20 right-4 z-[70] lg:bottom-7 lg:right-7">
    <button
      type="button"
      data-testid="contact-fab"
      class="group flex h-14 w-14 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-[0_12px_32px_rgba(0,0,0,0.24)] transition hover:-translate-y-1 hover:shadow-[0_16px_36px_rgba(0,0,0,0.28)] focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-primary/25"
      :style="minimal ? contactButtonStyle : undefined"
      :aria-label="t('contact.openButton')"
      @click="open = true"
    >
      <MessagesSquare class="h-6 w-6" />
      <span class="absolute right-0 top-0 h-3.5 w-3.5 rounded-full border-2 border-background bg-emerald-500" />
    </button>
  </div>

  <Teleport to="body">
    <Transition name="contact-fade">
      <div
        v-if="open"
        class="fixed inset-0 z-[100] flex items-end justify-center bg-black/50 p-0 backdrop-blur-sm sm:items-center sm:p-4"
        role="presentation"
        @click.self="open = false"
      >
        <section
          ref="panelRef"
          tabindex="-1"
          data-testid="contact-dialog"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="titleId"
          class="max-h-[88vh] w-full overflow-y-auto rounded-t-3xl border border-border bg-card p-5 text-card-foreground shadow-2xl sm:max-w-md sm:rounded-3xl sm:p-6"
        >
          <header class="mb-5 flex items-start justify-between gap-4">
            <div>
              <h2 :id="titleId" class="text-xl font-bold tracking-tight">{{ t('contact.title') }}</h2>
              <p class="mt-1 text-sm text-muted-foreground">{{ t('contact.subtitle') }}</p>
            </div>
            <button type="button" class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-secondary text-muted-foreground transition hover:text-foreground" :aria-label="t('announcement.close')" @click="open = false">
              <X class="h-4 w-4" />
            </button>
          </header>

          <div class="space-y-2.5">
            <component
              :is="item.href ? 'a' : 'button'"
              v-for="item in contactItems"
              :key="item.kind"
              :data-contact-kind="item.kind"
              :href="item.href || undefined"
              :target="item.href ? '_blank' : undefined"
              :rel="item.href ? 'noopener noreferrer' : undefined"
              type="button"
              class="flex w-full items-center gap-3 rounded-2xl border border-border bg-background p-3.5 text-left transition hover:border-primary/30 hover:bg-secondary/50"
              @click="item.href ? undefined : copyValue(item)"
            >
              <span class="grid h-11 w-11 shrink-0 place-items-center rounded-xl text-white" :class="item.iconClass">
                <span v-if="item.kind === 'qq'" class="text-[11px] font-black tracking-[-0.08em]">QQ</span>
                <ContactBrandIcon v-else :kind="item.kind" class="h-6 w-6" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="block text-sm font-semibold">{{ item.label }}</span>
                <span class="block truncate text-sm text-muted-foreground">{{ item.value }}</span>
              </span>
              <span class="shrink-0 text-xs font-medium text-primary">
                {{ copiedKind === item.kind ? t('contact.copied') : item.href ? t('contact.open') : t('contact.copy') }}
              </span>
            </component>
          </div>

          <div v-if="qrCode" class="mt-5 border-t border-border pt-5 text-center">
            <p class="mb-3 text-sm font-semibold">{{ t('contact.qrCode') }}</p>
            <div class="mx-auto w-fit rounded-2xl border border-border bg-white p-3 shadow-sm">
              <img :src="getImageUrl(qrCode)" :alt="t('contact.qrCode')" class="h-44 w-44 object-contain" />
            </div>
          </div>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { MessagesSquare, X } from 'lucide-vue-next'
import { useAppStore } from '../../stores/app'
import { getImageUrl } from '../../utils/image'
import { useTheme } from '../../utils/theme'
import ContactBrandIcon from './ContactBrandIcon.vue'

type ContactKind = 'telegram' | 'whatsapp' | 'wechat' | 'qq'
type ContactItem = {
  kind: ContactKind
  label: string
  value: string
  href: string
  iconClass: string
}

const props = withDefaults(defineProps<{ minimal?: boolean }>(), {
  minimal: false,
})
const { t } = useI18n()
const appStore = useAppStore()
const { theme } = useTheme()
const open = ref(false)
const copiedKind = ref<ContactKind | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const titleId = `contact-dialog-${Math.random().toString(36).slice(2)}`
let copiedTimer = 0

const contact = computed(() => appStore.config?.contact || {})
const defaultButtonColors = {
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
    return channel <= 0.03928 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4
  }
  const luminance = 0.2126 * toLinearChannel(1)
    + 0.7152 * toLinearChannel(3)
    + 0.0722 * toLinearChannel(5)
  return luminance > 0.42 ? '#0B1120' : '#FFFFFF'
}
const contactButtonStyle = computed(() => {
  if (!props.minimal) return undefined
  const configured = appStore.config?.minimal_theme_colors || {}
  const dark = theme.value === 'dark'
  const backgroundColor = normalizeHexColor(
    dark ? configured.button_dark : configured.button_light,
    dark ? defaultButtonColors.button_dark : defaultButtonColors.button_light,
  )
  return {
    backgroundColor,
    color: contrastTextColor(backgroundColor),
  }
})
const qrCode = computed(() => String(contact.value?.qr_code || '').trim())
const isLink = (value: string) => /^(https?:\/\/|mailto:|tel:)/i.test(value)

const contactItems = computed<ContactItem[]>(() => {
  const definitions: Array<[ContactKind, string, string]> = [
    ['telegram', t('contact.telegram'), 'bg-[#229ED9]'],
    ['whatsapp', t('contact.whatsapp'), 'bg-[#25D366]'],
    ['wechat', t('contact.wechat'), 'bg-[#07C160]'],
    ['qq', t('contact.qq'), 'bg-[#12B7F5]'],
  ]
  return definitions.flatMap(([kind, label, iconClass]) => {
    const value = String(contact.value?.[kind] || '').trim()
    return value ? [{ kind, label, value, href: isLink(value) ? value : '', iconClass }] : []
  })
})

const hasContact = computed(() => contactItems.value.length > 0 || Boolean(qrCode.value))

const copyValue = async (item: ContactItem) => {
  try {
    await navigator.clipboard.writeText(item.value)
    copiedKind.value = item.kind
    window.clearTimeout(copiedTimer)
    copiedTimer = window.setTimeout(() => { copiedKind.value = null }, 1800)
  } catch {
    // Clipboard permissions may be denied; keep the value visible for manual copy.
  }
}

const onKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') open.value = false
}

watch(open, async (visible) => {
  document.body.style.overflow = visible ? 'hidden' : ''
  if (visible) {
    window.addEventListener('keydown', onKeydown)
    await nextTick()
    panelRef.value?.focus()
  } else {
    window.removeEventListener('keydown', onKeydown)
  }
})

onBeforeUnmount(() => {
  document.body.style.overflow = ''
  window.removeEventListener('keydown', onKeydown)
  window.clearTimeout(copiedTimer)
})
</script>

<style scoped>
.contact-fade-enter-active,
.contact-fade-leave-active {
  transition: opacity 180ms ease;
}
.contact-fade-enter-from,
.contact-fade-leave-to {
  opacity: 0;
}
</style>
