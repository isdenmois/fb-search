import { Elysia, t } from 'elysia'
import { swagger } from '@elysiajs/swagger'
import { searchBooks, searchBooksRu, findAllByFile, findFileById } from './core/db'
import { getFile } from './core/file'
import { parser } from './core/inpx-parser'
import { slugify } from 'transliteration'
import { staticPlugin } from '@elysiajs/static'

import { migrate } from 'drizzle-orm/node-postgres/migrator'
import { db } from './core/db/db'

await migrate(db, { migrationsFolder: './drizzle' })

const port = +(process.env.PORT || 3000)

function containsCyrillic(str: string) {
  return /[\u0400-\u04FF]/.test(str)
}

new Elysia({
  serve: {
    port,
    maxRequestBodySize: 1024 * 1024 * 256, // 256MB
  },
})
  .use(swagger())
  .use(staticPlugin({ prefix: '/' }))
  .get(
    '/api/search',
    ({ query }) => {
      const q = query.q.toLowerCase()

      return containsCyrillic(q) ? searchBooksRu(q) : searchBooks(q)
    },
    {
      query: t.Object({
        q: t.String(),
      }),
    },
  )
  .get('/api/by-file/:file', ({ params: { file } }) => findAllByFile(file))
  .get(
    '/dl/:id',
    async ({ params: { id }, set }) => {
      const file = findFileById(id)

      if (file?.file && file.path) {
        const { stream, size } = (await getFile(file.file, file.path)) ?? {}

        const filename = [
          file.authors ? slugify(file.authors) : null,
          file.title ? slugify(file.title) : null,
          file.path,
        ]
          .filter(Boolean)
          .join('.')

        set.headers['content-disposition'] = `attachment; filename="${filename}"`
        set.headers['content-type'] = 'text/fb2+xml'

        if (size) {
          set.headers['content-length'] = String(size)
        }

        return new Response(stream!).bytes()
      }

      set.status = 404
      return { ok: false }
    },
    { params: t.Object({ id: t.Number() }) },
  )
  .post('/api/parser', () => {
    parser.parse()

    return parser
  })
  .get('/api/parser', () => parser)
  .listen(port)
