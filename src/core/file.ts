import StreamZip from 'node-stream-zip'
import { XMLParser, XMLBuilder } from 'fast-xml-parser'
import { basename } from 'path'

const parser = new XMLParser({
  ignoreAttributes: false,
  allowBooleanAttributes: true,
})
const builder = new XMLBuilder({
  ignoreAttributes: false,
})

export async function getFile(file: string, path: string) {
  const zPath = `files/${file}`.replace(/\.zip$/, '.7z')
  const proc = Bun.spawn(['7z', 'x', '-so', zPath, path])

  await proc.exited

  if (proc.exitCode === 0) {
    const id = path.replace('.fb2', '')
    const withCover = await hasCover(file, id)
    const images = await getImages(file, id)

    console.log({ id, withCover, images })

    if (withCover || images.length) {
      const text = await new Response(proc.stdout).text()
      const jObj = parser.parse(text)
      let binary: any[]

      if ('binary' in jObj.FictionBook) {
        if (!Array.isArray(jObj.FictionBook.binary)) {
          jObj.FictionBook.binary = [jObj.FictionBook.binary]
        }
      } else {
        jObj.FictionBook.binary = []
      }
      binary = jObj.FictionBook.binary

      if (withCover) {
        const cover = await getCoverBase64(file, id)

        if (cover) {
          binary.push({ '@_id': 'cover', '@_content-type': 'image/jpeg', '#text': cover })
        }
      }

      for (const imgFile of images) {
        const img = await getImageBase64(file, imgFile)

        if (img) {
          const imgId = basename(imgFile)

          binary.push({ '@_id': imgId, '@_content-type': 'image/jpeg', '#text': img })
        }
      }

      return builder.build(jObj)
    }

    return {
      stream: proc.stdout,
      size: proc.stdout.length,
    }
  }

  return null
}

export async function getCover(file: string, id: string) {
  const zip = new StreamZip.async({ file: `files/covers/${file}` })

  const entry = await zip.entry(id)

  return entry
    ? {
        stream: await zip.stream(entry),
        size: entry.size,
      }
    : null
}

function getCoverBase64(file: string, id: string) {
  return getFileBase64(`files/covers/${file}`, id)
}

function getImageBase64(file: string, path: string) {
  return getFileBase64(`files/images/${file}`, path)
}

async function getFileBase64(file: string, include: string) {
  const proc = Bun.spawn(['7z', 'x', '-so', file, include])

  await proc.exited

  if (proc.exitCode === 0) {
    const arrBuf = await Bun.readableStreamToArrayBuffer(proc.stdout)
    const buf = Buffer.from(arrBuf)

    return buf.toString('base64')
  }

  return null
}

async function listFiles(file: string, include: string) {
  const proc = Bun.spawn(['7z', 'l', '-ba', file, include])

  await proc.exited

  if (proc.exitCode === 0) {
    const text = await new Response(proc.stdout).text()
    const files = text
      .split('\n')
      .filter(Boolean)
      .map((file) => file.match(/(\S+)$/)?.['1'])
      .filter(Boolean)

    return files as string[]
  }

  return []
}

async function hasCover(file: string, id: string) {
  const zPath = `files/covers/${file}`
  const files = await listFiles(zPath, id)

  return files.length > 0
}

async function getImages(file: string, id: string) {
  const zPath = `files/images/${file}`
  const files = await listFiles(zPath, `${id}/*`)

  return files
}
