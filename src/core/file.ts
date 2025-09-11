import type { Subprocess } from 'bun'
import fs from 'node:fs'
import { basename } from 'path'

export async function getFile(file: string, path: string) {
  const z7Path = `files/${file}`.replace(/\.zip$/, '.7z')

  if (!fs.existsSync(z7Path)) {
    return extractFile(`files/${file}`, path)
  }

  const stdout = await extractFile(z7Path, path)

  if (stdout) {
    const id = path.replace('.fb2', '')
    const withCover = await hasCover(file, id)
    const images = await getImages(file, id)

    const binaries = [
      { id: 'cover', text: await getCoverBase64(file, id) },
      ...(await Promise.all(images.map(async (img) => ({ id: basename(img), text: await getImageBase64(file, img) })))),
    ]

    if (withCover || images.length) {
      const text = await new Response(stdout).text()
      const b = binaries.map((bin) => `<binary id="${bin.id}" content-type="image/jpeg">${bin.text}</binary>`).join('')

      return text.replace(/<\/\s*FictionBook>/i, (match) => `${b}${match}`)
    }

    return stdout
  }

  return null
}

function spawnExtract(file: string, include: string) {
  return Bun.spawn(['7z', 'x', '-so', file, include])
}

async function procToStream<T>(proc: Subprocess<ReadableStream<T>>) {
  await proc.exited

  return proc.stdout as ReadableStream<T>
}

async function extractFile(file: string, include: string) {
  const proc = spawnExtract(file, include)

  await proc.exited

  return proc.stdout
}

function extractJpegXl(file: string, include: string) {
  const proc = Bun.spawn(['djxl', '--quiet', '--output_format', 'jpeg', '-', '-'], {
    stdin: spawnExtract(file, include).stdout,
  })

  return procToStream(proc)
}

export function getCover(file: string, id: string) {
  return extractJpegXl(`files/covers/${file}`, id)
}

function getImage(file: string, path: string) {
  return extractJpegXl(`files/images/${file}`, path)
}

async function getCoverBase64(file: string, id: string) {
  return getBase64(await getCover(file, id))
}

async function getImageBase64(file: string, path: string) {
  return getBase64(await getImage(file, path))
}

export async function getText(file: string, include: string) {
  const stdout = await extractFile(file, include)

  if (stdout) {
    return await new Response(stdout).text()
  }

  return ''
}

async function getBase64(from: ReadableStream) {
  const arrBuf = await Bun.readableStreamToArrayBuffer(from)
  const buf = Buffer.from(arrBuf)

  return buf.toString('base64')
}

export async function listFiles(file: string, include: string) {
  const proc = Bun.spawn(['7z', 'l', '-ba', file, include])

  await proc.exited

  const text = await new Response(proc.stdout).text()
  const files = text
    .split('\n')
    .filter(Boolean)
    .map((file) => file.match(/(\S+)$/)?.['1'])
    .filter(Boolean)

  return files as string[]
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
