<script setup lang="ts">
import { ref } from 'vue'
import { api } from 'shared/api'
import type { Book } from 'shared/api/fb'
import { InputField, SearchIcon, LoadingSpinner } from 'shared/ui'
import { BookItem } from 'entities/book'
import { getErrorMessage } from 'shared/lib'

const query = ref('')
const books = ref<Book[] | null>(null)
const disabled = ref(false)
const error = ref<string | null>(null)

const search = async () => {
  disabled.value = true
  books.value = null
  error.value = null

  try {
    books.value = await api.fb.search(query.value)
  } catch (err) {
    error.value = getErrorMessage(err) || 'Unexpected error'
  } finally {
    disabled.value = false
  }
}
</script>

<template>
  <div class="pt-4 flex flex-col flex-1 overflow-hidden gap-4">
    <form class="px-4" @submit.prevent="search">
      <InputField v-model="query" :disabled="disabled">
        <div class="flex items-center px-1">
          <SearchIcon />
        </div>
      </InputField>
    </form>

    <LoadingSpinner v-if="disabled" />

    <ul v-if="books?.length" class="flex flex-col gap-3 flex-1 p-2 overflow-y-auto">
      <li v-for="book in books" :key="book.id">
        <a :href="`/dl/${book.id}`">
          <BookItem :book="book" />
        </a>
      </li>
    </ul>

    <div v-if="books && !books.length" class="text-white text-align-center">Nothing has found</div>

    <div v-if="error" class="text-red text-align-center">{{ error }}</div>
  </div>
</template>
