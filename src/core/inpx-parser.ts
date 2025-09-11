import { spawn } from 'node:child_process'
import { existsSync } from 'node:fs'
import { basename } from 'node:path'
import * as CSV from 'csv-parse'
import { db, insertBooks } from './db'
import { books, type Book } from './db/schema'
import { sql } from 'drizzle-orm'
import { listFiles } from './file'

const MAX = 2048

export const parser = {
  currentFile: 0,
  totalFiles: 1,
  booksImported: 0,
  took: 0,
  async parse() {
    parser.booksImported = 0
    parser.currentFile = 0
    parser.took = 0

    const start = Date.now()

    await db.execute(sql`TRUNCATE TABLE ${books} RESTART IDENTITY`)

    await parseInpx('./files/fb2.Flibusta.Net.7z.inpx')

    parser.took = (Date.now() - start) / 1000
  },
}

export async function parseInpx(inpxPath: string) {
  parser.booksImported = 0
  parser.currentFile = 0

  console.log('Start parsing INPX file', inpxPath)
  const files = await listFiles(inpxPath, '*.inp')

  const filesCount = files.length

  console.log('Files count', filesCount)

  for (const file of files) {
    parser.currentFile++

    parser.booksImported += await parseInp(inpxPath, file)
  }
}

function createCsvParser() {
  return CSV.parse({
    delimiter: '\x04',
    relax_quotes: true,
    relax_column_count: true,
    // columns: [
    //   "author",
    //   "genre",
    //   "title",
    //   "series",
    //   "serno",
    //   "file",
    //   "size",
    //   "libid",
    //   "del",
    //   "ext",
    //   "date",
    //   "lang",
    //   "librate",
    //   "keywords",
    // ],
  })
}

function trunc(s: string) {
  return s.length > MAX ? s.slice(0, MAX) : s
}

async function parseInp(inpx: string, entry: string) {
  const base = entry.substring(0, entry.lastIndexOf('.'))
  const zip = `./files/${base}.7z`

  if (!existsSync(zip)) {
    return 0
  }

  console.log('Parse INP file', entry)
  const parser = createCsvParser()
  const proc = spawn('7z', ['x', '-so', inpx, entry])

  let count = 0
  const books: Book[] = []
  const file = `${base}.zip`

  // biome-ignore lint/complexity/noForEach: <explanation>
  await proc.stdout.pipe(parser).forEach((data) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [author, genre, title, series, serno, filename, size, libid, del, ext, date, lang, librate, keywords] = data
    // console.log(data);
    const authors = author
      .split(':')
      .map((s: string) => s.replace(/,/g, ' ').trim())
      .filter((s: string) => s)
      .join(',')
    const path = `${filename}.${ext}`

    books.push({
      authors,
      title: trunc(title),
      file,
      lang,
      // keywords,
      // genre,
      series,
      serno,
      size: +size || 0,
      // librate,
      path,
      search: trunc([authors, title, series].filter(Boolean).join(' ').toLocaleLowerCase()),
    })
    count++
  })

  await insertBooks(books)

  return count
}
