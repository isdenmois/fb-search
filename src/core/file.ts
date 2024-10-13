import StreamZip, { type StreamZipAsync, type ZipEntry } from "node-stream-zip";

export async function getFile(file: string, path: string) {
  const zip = new StreamZip.async({ file: `files/${file}` });
  const entry = await zip.entry(path);

  return entry
    ? {
        stream: zip.stream(entry),
        size: entry.size,
      }
    : null;
}
