import { sql } from "drizzle-orm";
import { db, schema } from "./src/core/db";
import { parseInpx } from "./src/core/inpx-parser";

await db.delete(schema.books).execute();

await parseInpx("./files/flibusta_fb2_local.inpx");

const result = await db
  .select()
  .from(schema.books)
  .where(sql`${schema.books.search} LIKE '%игра престолов%'`)
  .execute();

console.log(result);
