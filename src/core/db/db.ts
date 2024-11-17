import { drizzle } from 'drizzle-orm/pglite'

export const db = drizzle('data/pg')
