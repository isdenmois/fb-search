import { sql } from 'drizzle-orm'
import { pgTable, integer, text, index } from 'drizzle-orm/pg-core'

export interface Book {
  title: string
  search: string
  authors?: string
  series?: string
  serno?: string
  file: string
  path: string
  lang?: string
  size: number
}

export const books = pgTable(
  'books',
  {
    id: integer().primaryKey().generatedByDefaultAsIdentity(),
    title: text().notNull(),
    search: text().notNull(),
    authors: text(),
    series: text(),
    serno: text(),
    file: text(),
    path: text(),
    lang: text(),
    size: integer(),
  },
  (table) => ({
    searchRuIdx: index('search_ru_idx').using('gin', sql`to_tsvector('russian', ${table.search})`),
    searchSimpleIdx: index('search_simple_idx').using('gin', sql`to_tsvector('simple', ${table.search})`),
  }),
)

// export const search = sql`to_tsvector(${books.title} || ${books.authors})`
