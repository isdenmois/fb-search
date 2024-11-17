import { migrate } from 'drizzle-orm/pglite/migrator'
import { db } from '../src/core/db/db'

await migrate(db, { migrationsFolder: './drizzle' })
