<script setup lang="ts">
import { ref } from 'vue'
import { api } from 'shared/api'
import type { Book } from 'shared/api/fb'
import { InputField, SearchIcon, LoadingSpinner } from 'shared/ui'
import { BookItem } from 'entities/book'

const query = ref('')
const disabled = ref(false)
const books = ref<Book[]>([])

const search = async () => {
  disabled.value = true
  books.value = []

  try {
    books.value = await api.fb.search(query.value)
  } finally {
    disabled.value = false
  }
}
</script>

<template>
  <div class="pt-4 flex flex-col flex-1 overflow-hidden">
    <form class="px-4" @submit.prevent="search">
      <InputField v-model="query" :disabled="disabled">
        <div class="flex items-center px-1">
          <SearchIcon />
        </div>
      </InputField>
    </form>

    <LoadingSpinner v-if="disabled" />

    <ul class="flex flex-col gap-3 flex-1 mt-4 p-2 overflow-y-auto">
      <li v-for="book in books" :key="book.id">
        <a :href="`/dl/${book.id}`">
          <BookItem :book="book" />
        </a>
      </li>
    </ul>
  </div>
</template>
