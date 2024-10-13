import StreamZip, { type StreamZipAsync, type ZipEntry } from "node-stream-zip";
import * as CSV from "csv-parse";
import { insertBooks } from "./db";

export async function parseInpx(inpxPath: string) {
  console.log("Start parsing INPX file", inpxPath);
  const zip = new StreamZip.async({ file: inpxPath });
  const entriesCount = await zip.entriesCount;

  console.log("Entries Count", entriesCount);

  const entries = await zip.entries();
  let count = 0;

  for (const entry of Object.values(entries)) {
    if (!entry.isDirectory && entry.name.endsWith(".inp")) {
      count += await parseInp(zip, entry);
    }
  }
  zip.close();

  console.log("parsed", count, "entries");
}

function createCsvParser() {
  return CSV.parse({
    delimiter: "\x04",
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
  });
}

async function parseInp(zip: StreamZipAsync, entry: ZipEntry) {
  console.log("Parse INP file", entry.name);
  const stream = await zip.stream(entry);
  let count = 0;
  const books: any[] = [];
  const file = entry.name.replace(".inp", ".zip");

  // biome-ignore lint/complexity/noForEach: <explanation>
  await stream.pipe(createCsvParser()).forEach((data) => {
    const [
      author,
      genre,
      title,
      series,
      serno,
      filename,
      size,
      libid,
      del,
      ext,
      date,
      lang,
      librate,
      keywords,
    ] = data;
    // console.log(data);
    const authors = author
      .split(":")
      .map((s: string) => s.replace(/,/g, " ").trim())
      .filter((s: string) => s)
      .join(",");
    const path = `${filename}.${ext}`;

    books.push({
      authors,
      title,
      file,
      lang,
      path,
      search: [authors, title].join("|").toLocaleLowerCase(),
    });
    count++;
  });

  await insertBooks(books);

  return count;
}
