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

async function extractFile(file: string, include: string) {
  const proc = Bun.spawn(['7z', 'x', '-so', file, include])

  await proc.exited

  if (proc.exitCode === 0) {
    return proc.stdout
  }

  return null
}

export async function getCover(file: string, id: string) {
  return extractFile(`files/covers/${file}`, id)
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
