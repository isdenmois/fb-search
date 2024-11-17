<script lang="ts">
  import type { Book } from 'shared/api/fb'
  import { FileIcon, Item } from 'shared/ui'

  export let book: Book

  const ext = 'fb2'
  const subtitle = [book.lang, book.series, +book.serno].filter(Boolean).join(' ')
</script>

<Item suptitle={book.authors} title={book.title} {subtitle}>
  <FileIcon slot="icon" text={ext} />
</Item>
