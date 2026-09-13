import { ref } from 'vue'
import { api } from '@/shared/api'
import type { ParseProgress } from '@/shared/api/parse'
import { useApiKey, usePoll } from '@/shared/lib'

const progress = ref<ParseProgress | null>(null)
const rebuilding = ref(false)

const fetchProgress = async () => {
  try {
    progress.value = await api.parse.getProgress()
  } catch (e) {
    console.error(e)
  }
}

export const useRebuildDb = () => {
  const { apiKey } = useApiKey()
  const poll = usePoll(fetchProgress, 500)
  const start = async () => {
    rebuilding.value = true
    poll.start()

    try {
      progress.value = await api.parse.rebuild(apiKey.value)
    } catch (e) {
      console.error(e)
    } finally {
      rebuilding.value = false
      poll.stop()
    }
  }

  return {
    progress,
    rebuilding,
    start,
  }
}
