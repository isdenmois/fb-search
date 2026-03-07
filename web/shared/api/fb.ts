import { http } from './client'

export interface Book {
  id: string
  lang?: string
  authors?: string
  title: string
  size?: number
  series?: string
  serno?: string
}

export const search = (query: string) => http.url('/search').query({ q: query.trim() }).get().json() as Promise<Book[]>
