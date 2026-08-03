type AuthInvalidatedListener = () => void

const authInvalidatedListeners = new Set<AuthInvalidatedListener>()

export const subscribeAuthInvalidated = (listener: AuthInvalidatedListener) => {
  authInvalidatedListeners.add(listener)
  return () => authInvalidatedListeners.delete(listener)
}

export const invalidateUserAuth = () => {
  localStorage.removeItem('user_token')
  localStorage.removeItem('user_profile')
  authInvalidatedListeners.forEach((listener) => listener())
}

export const canFallbackToGuestCheckout = (path: string) => String(path || '').trim() === '/checkout'
