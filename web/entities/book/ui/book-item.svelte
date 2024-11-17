<script lang="ts">
  import type { Book } from 'shared/api/fb'
  import { FileIcon, Item } from 'shared/ui'

  export let book: Book

  const ext = 'fb2'
  const size = Math.round(book.size / 1024 || 0).toLocaleString('ru') + 'Kb'
  const subtitle = [size, book.series, +book.serno, book.lang].filter(Boolean).join(', ')
</script>

<Item suptitle={book.authors} title={book.title} {subtitle}>
  <FileIcon slot="icon" text={ext} />
</Item>
