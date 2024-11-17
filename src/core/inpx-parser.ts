import StreamZip, { type StreamZipAsync, type ZipEntry } from 'node-stream-zip'
import * as CSV from 'csv-parse'
import { db, insertBooks } from './db'
import { books, type Book } from './db/schema'
import { sql } from 'drizzle-orm'

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

    await parseInpx('./files/flibusta_fb2_local.inpx')

    parser.took = (Date.now() - start) / 1000
  },
}

export async function parseInpx(inpxPath: string) {
  parser.booksImported = 0
  parser.currentFile = 0

  console.log('Start parsing INPX file', inpxPath)
  const zip = new StreamZip.async({ file: inpxPath })
  const entriesCount = await zip.entriesCount

  console.log('Entries Count', entriesCount)

  const entries = Object.values(await zip.entries())
  parser.totalFiles = entries.length

  for (const entry of entries) {
    parser.currentFile++

    if (!entry.isDirectory && entry.name.endsWith('.inp')) {
      parser.booksImported += await parseInp(zip, entry)
    }
  }
  zip.close()
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

async function parseInp(zip: StreamZipAsync, entry: ZipEntry) {
  console.log('Parse INP file', entry.name)
  const stream = await zip.stream(entry)
  let count = 0
  const books: Book[] = []
  const file = entry.name.replace('.inp', '.zip')

  // biome-ignore lint/complexity/noForEach: <explanation>
  await stream.pipe(createCsvParser()).forEach((data) => {
    const [author, genre, title, series, serno, filename, size, libid, del, ext, date, lang, librate, keywords] = data
    // console.log(data);
    const authors = author
      .split(':')
      .map((s: string) => s.replace(/,/g, ' ').trim())
      .filter((s: string) => s)
      .join(',')
    const path = `${filename}.${ext}`

    // if (title.toLowerCase().includes('пепел и сталь')) {
    //   console.log(data)
    // }

    books.push({
      authors,
      title: trunc(title),
      file,
      lang,
      // keywords,
      // genre,
      series,
      serno,
      // librate,
      path,
      search: trunc([authors, title, series].filter(Boolean).join(' ').toLocaleLowerCase()),
    })
    count++
  })

  await insertBooks(books)

  return count
}
