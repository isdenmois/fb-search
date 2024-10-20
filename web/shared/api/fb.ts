import { http } from './client'

export interface Book {
  id: number
  lang: string
  authors: string
  title: string
}

export const search = (query: string) => http.url('/search').query({ q: query.trim() }).get().json() as Promise<Book[]>
