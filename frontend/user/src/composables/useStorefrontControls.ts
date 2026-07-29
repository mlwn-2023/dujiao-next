import { computed } from 'vue'
import { useAppStore } from '../stores/app'

export function useStorefrontControls() {
  const appStore = useAppStore()
  const controls = computed(() => appStore.config?.ui_controls || {})

  return {
    showBanner: computed(() => controls.value.show_banner !== false),
    showLanguageSwitcher: computed(() => controls.value.show_language_switcher !== false),
    showThemeSwitcher: computed(() => controls.value.show_theme_switcher !== false),
    showContactFloatingButton: computed(() => controls.value.show_contact_floating_button !== false),
  }
}
