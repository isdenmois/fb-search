import { http } from './client'

export interface Book {
  id: number
  lang: string
  authors?: string
  title: string
  size?: number
  series?: string
  serno?: string
}

export const search = (query: string) =>
  http.url('/v2/search').query({ q: query.trim() }).get().json() as Promise<Book[]>
