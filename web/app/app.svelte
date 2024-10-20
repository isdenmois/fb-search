<script lang="ts">
  import { api } from 'shared/api'
  import type { Book } from 'shared/api/fb'
  import { Input, Item, SearchIcon, Spinner } from 'shared/ui'
  import { BookItem } from 'entities/book'

  let query = ''
  let disabled = false
  let books: Book[] = []

  const search = async () => {
    disabled = true
    books = []

    try {
      books = await api.fb.search(query)
    } finally {
      disabled = false
    }
  }
</script>

<div class="pt-4 flex flex-col flex-1 overflow-hidden">
  <form class="px-4" on:submit|preventDefault={search}>
    <Input bind:value={query} {disabled}>
      <SearchIcon slot="icon" />
    </Input>
  </form>

  {#if disabled}
    <Spinner />
  {/if}

  <ul class="flex flex-col gap-3 flex-1 mt-4 p-2 overflow-y-auto">
    {#each books as book (book.id)}
      <a href={`/dl/${book.id}`}>
        <BookItem {book} />
      </a>
    {/each}
  </ul>
</div>
