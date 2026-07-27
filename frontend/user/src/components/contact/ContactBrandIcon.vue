<template>
  <svg v-if="kind === 'telegram'" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
    <path d="M21.6 3.1 18.4 19c-.2 1.1-.9 1.4-1.8.9l-4.9-3.6-2.4 2.3c-.3.3-.5.5-1 .5l.4-5 9.1-8.2c.4-.4-.1-.6-.6-.2L6 12.8 1.2 11.3c-1-.3-1.1-1 .2-1.5L20.3 2.5c.9-.3 1.7.2 1.3.6Z" />
  </svg>
  <svg v-else-if="kind === 'whatsapp'" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
    <path d="M12 2a9.8 9.8 0 0 0-8.4 14.9L2.2 22l5.2-1.4A10 10 0 1 0 12 2Zm0 17.9a8 8 0 0 1-4.1-1.1l-.3-.2-3 .8.8-2.9-.2-.3A8 8 0 1 1 12 20Zm4.4-6c-.2-.1-1.4-.7-1.7-.8-.2-.1-.4-.1-.5.1l-.7.9c-.1.2-.3.2-.5.1a6.4 6.4 0 0 1-1.9-1.2 7 7 0 0 1-1.3-1.7c-.1-.2 0-.4.1-.5l.4-.4.2-.4c.1-.2 0-.3 0-.4l-.7-1.8c-.2-.5-.4-.4-.6-.4h-.5c-.2 0-.4.1-.7.3-.2.3-.8.8-.8 2 0 1.1.8 2.3 1 2.4.1.2 1.6 2.6 4.1 3.5 2 .8 2.4.6 2.9.6.4-.1 1.4-.6 1.6-1.1.2-.6.2-1 .1-1.1l-.5-.3Z" />
  </svg>
  <svg v-else-if="kind === 'wechat'" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
    <path d="M9.4 3.2C4.8 3.2 1 6.3 1 10.1c0 2.2 1.3 4.2 3.3 5.4l-.8 2.4 2.9-1.4c.9.3 1.9.5 3 .5h.5a6.2 6.2 0 0 1-.2-1.5c0-3.6 3.3-6.5 7.5-6.5h.8c-.7-3.3-4.2-5.8-8.6-5.8ZM6.5 7.1a1.1 1.1 0 1 1 0 2.2 1.1 1.1 0 0 1 0-2.2Zm5.6 0a1.1 1.1 0 1 1 0 2.2 1.1 1.1 0 0 1 0-2.2Z" />
    <path d="M23 15.5c0-3-2.9-5.5-6.4-5.5s-6.4 2.5-6.4 5.5 2.9 5.5 6.4 5.5c.8 0 1.6-.1 2.3-.4l2.2 1.1-.6-1.9c1.5-1 2.5-2.5 2.5-4.3Zm-8.6-.7a.9.9 0 1 1 0-1.8.9.9 0 0 1 0 1.8Zm4.4 0a.9.9 0 1 1 0-1.8.9.9 0 0 1 0 1.8Z" />
  </svg>
  <svg v-else viewBox="0 0 24 24" fill="none" aria-hidden="true">
    <circle cx="12" cy="12" r="10" fill="currentColor" />
    <text x="12" y="15" fill="white" font-family="Arial, sans-serif" font-size="8" font-weight="800" text-anchor="middle">QQ</text>
  </svg>
</template>

<script setup lang="ts">
defineProps<{ kind: 'telegram' | 'whatsapp' | 'wechat' | 'qq' }>()
</script>
