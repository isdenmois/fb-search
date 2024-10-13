import { sql } from "drizzle-orm";
import { db } from "./db";
import * as schema from "./schema";

export * from "./schema";

export { db, schema };

const limit = 10_000;

export async function insertBooks(
  values: Array<{
    authors: string;
    title: string;
    file: string;
    path: string;
    lang: string;
    search: string;
  }>
) {
  for (let i = 0; i < values.length; i += limit) {
    db.insert(schema.books)
      .values(values.slice(i, Math.min(i + limit, values.length)))
      .execute();
  }
}
