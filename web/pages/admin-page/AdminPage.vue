<script setup lang="ts">
import { useRebuildDb } from '@/features/rebuild-db'

const { start: startRebuild, progress, rebuilding } = useRebuildDb()
</script>

<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">Admin</h1>

    <button
      class="px-4 py-2 bg-red-500 text-white rounded cursor-pointer disabled:opacity-50"
      :disabled="rebuilding"
      @click="startRebuild"
    >
      {{ rebuilding ? "Rebuilding..." : "Rebuild Database" }}
    </button>

    <div v-if="progress?.books" class="mt-4 text-white">
      <p>Files: {{ progress.files }}</p>
      <p>Books: {{ progress.books }}</p>
      <p>Time: {{ progress.time }} ms</p>
    </div>
  </div>
</template>
