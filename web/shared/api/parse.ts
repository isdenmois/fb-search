import { http } from './client'

export interface ParseProgress {
  files: number
  books: number
  time: string
}

export const getProgress = () => http.get('/parse').json() as Promise<ParseProgress>

export const rebuild = (key: string) =>
  http.url('/parse/rebuild').headers({ 'X-API-Key': key }).post({}).json() as Promise<ParseProgress>
