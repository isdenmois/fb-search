import { migrate } from 'drizzle-orm/node-postgres/migrator'
import { db } from '../src/core/db/db'

await migrate(db, { migrationsFolder: './drizzle' })
