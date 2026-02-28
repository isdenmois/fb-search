import { http } from './client'

export interface ParseProgress {
  files: number
  books: number
  time: string
}

export const getProgress = () => http.get('/parse').json() as Promise<ParseProgress>

export const rebuild = () => http.url('/parse/rebuild').post({}).json() as Promise<ParseProgress>
