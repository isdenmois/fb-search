import { ref } from 'vue'

const STORAGE_KEY = 'admin-api-key'

const readStoredKey = () => {
  try {
    return sessionStorage.getItem(STORAGE_KEY) ?? ''
  } catch {
    return ''
  }
}

// Shared singleton: every consumer sees and mutates the same key.
const apiKey = ref(readStoredKey())

const setApiKey = (key: string) => {
  apiKey.value = key

  try {
    if (key) {
      sessionStorage.setItem(STORAGE_KEY, key)
    } else {
      sessionStorage.removeItem(STORAGE_KEY)
    }
  } catch {
    // sessionStorage unavailable — keep the in-memory value only
  }
}

export const useApiKey = () => ({
  apiKey,
  setApiKey,
})
