<script lang="ts">
  import type { Book } from 'shared/api/fb'
  import { FileIcon, Item } from 'shared/ui'

  export let book: Book

  const ext = 'fb2'
</script>

<Item suptitle={book.authors} title={book.title} subtitle={book.lang}>
  <FileIcon slot="icon" text={ext} />
</Item>
