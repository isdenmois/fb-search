import { like, eq, sql } from 'drizzle-orm'

import { db } from './db'
import { books, search, type Book } from './schema'

const limit = 1024

export async function insertBooks(values: Book[]) {
  for (let i = 0; i < values.length; i += limit) {
    const slice = values.slice(i, Math.min(i + limit, values.length))

    try {
      await db.insert(books).values(slice).execute()
    } catch (e) {
      console.log(slice)
      throw e
    }
  }
}

export function searchBooks(query: string) {
  return (
    db
      .select()
      .from(books)
      // .where(like(books.search, `%${query.replaceAll('*', '%').toLowerCase()}%`))
      .where(sql`to_tsvector('simple', ${books.search}) @@ websearch_to_tsquery('simple', ${query})`)
      .limit(100)
      // .orderBy(sql`websearch_to_tsquery(${query})`)
      .execute()
  )
}

export function searchBooksRu(query: string) {
  return (
    db
      .select()
      .from(books)
      // .where(like(books.search, `%${query.replaceAll('*', '%').toLowerCase()}%`))
      .where(sql`to_tsvector('russian', ${books.search}) @@ websearch_to_tsquery('russian', ${query})`)
      .limit(100)
      // .orderBy(sql`websearch_to_tsquery(${query})`)
      .execute()
  )
}

export function findAllByFile(file: string) {
  return db
    .select({
      id: books.id,
      title: books.title,
      authors: books.authors,
      lang: books.lang,
    })
    .from(books)
    .where(eq(books.file, file))
    .execute()
}

export async function findFileById(id: number) {
  const result = await db.select().from(books).where(eq(books.id, id)).execute()

  return result[0]
}
