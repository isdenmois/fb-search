<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

import { api } from 'shared/api'
import type { ParseProgress } from 'shared/api/fb'

const progress = ref<ParseProgress | null>(null)
const rebuilding = ref(false)
let interval: ReturnType<typeof setInterval> | null = null

const fetchProgress = async () => {
  try {
    progress.value = await api.parse.getProgress()
  } catch (e) {
    console.error(e)
  }
}

const startRebuild = async () => {
  rebuilding.value = true
  interval = setInterval(fetchProgress, 500)

  try {
    progress.value = await api.parse.rebuild()
  } catch (e) {
    console.error(e)
  } finally {
    rebuilding.value = false
    if (interval) {
      clearInterval(interval)
      interval = null
    }
  }
}

onMounted(fetchProgress)
onUnmounted(() => {
  if (interval) {
    clearInterval(interval)
    interval = null
  }
})
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Admin</h1>

    <button
      class="px-4 py-2 bg-red-500 text-white rounded cursor-pointer disabled:opacity-50"
      :disabled="rebuilding"
      @click="startRebuild"
    >
      {{ rebuilding ? 'Rebuilding...' : 'Rebuild Database' }}
    </button>

    <div v-if="progress?.books" class="mt-4 text-white">
      <p>Files: {{ progress.files }}</p>
      <p>Books: {{ progress.books }}</p>
      <p>Time: {{ progress.time }} ms</p>
    </div>
  </div>
</template>
