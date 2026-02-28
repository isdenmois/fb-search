<script setup lang="ts">
import type { Book } from 'shared/api/fb'
import { FileIcon, ItemCard } from 'shared/ui'

const props = defineProps<{
  book: Book
}>()

const ext = 'fb2'
const size = Math.round(props.book.size / 1024 || 0).toLocaleString('ru') + 'Kb'
const subtitle = [size, props.book.series, props.book.serno, props.book.lang].filter(Boolean).join(', ')
</script>

<template>
  <ItemCard :suptitle="book.authors" :title="book.title" :subtitle="subtitle">
    <template #icon>
      <FileIcon :text="ext" />
    </template>
  </ItemCard>
</template>
