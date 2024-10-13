import { SQL, sql } from "drizzle-orm";
import {
  sqliteTable,
  sqliteView,
  text,
  index,
  integer,
  type AnySQLiteColumn,
} from "drizzle-orm/sqlite-core";

export const books = sqliteTable(
  "books",
  {
    id: integer().primaryKey({ autoIncrement: true }),
    title: text().notNull(),
    search: text().notNull(),
    authors: text(),
    file: text(),
    path: text(),
    lang: text(),
  },
  (table) => ({
    searchIdx: index("search_idx").on(table.search),
  })
);
