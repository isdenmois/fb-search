import { onMounted, onUnmounted } from 'vue'

export const usePoll = (fn: () => void, timeout = 500) => {
  let interval: ReturnType<typeof setInterval> | null = null

  onMounted(fn)
  onUnmounted(() => {
    if (interval) {
      clearInterval(interval)
      interval = null
    }
  })

  return {
    start() {
      interval = setInterval(fn, timeout)
    },
    stop() {
      if (interval) {
        clearInterval(interval)
        interval = null
      }
    },
  }
}
