import { Elysia, t } from "elysia";
import { swagger } from "@elysiajs/swagger";
import { db, books } from "./core/db";
import { like, eq } from "drizzle-orm";
import { getFile } from "./core/file";
import { slugify } from "transliteration";
import { staticPlugin } from "@elysiajs/static";

const port = +(process.env.PORT || 3000);

const app = new Elysia({
  serve: {
    port,
    maxRequestBodySize: 1024 * 1024 * 256, // 256MB
  },
})
  .use(swagger())
  .use(staticPlugin({ prefix: "/" }))
  .get(
    "/api/search",
    ({ query }) => {
      return db
        .select({
          id: books.id,
          title: books.title,
          authors: books.authors,
          lang: books.lang,
        })
        .from(books)
        .where(
          like(books.search, `%${query.q.replaceAll("*", "%").toLowerCase()}%`)
        )
        .limit(100)
        .execute();
    },
    {
      query: t.Object({
        q: t.String(),
      }),
    }
  )
  .get("/api/by-file/:file", ({ params: { file } }) =>
    db
      .select({
        id: books.id,
        title: books.title,
        authors: books.authors,
        lang: books.lang,
      })
      .from(books)
      .where(eq(books.file, file))
      .execute()
  )
  .get(
    "/dl/:id",
    async ({ params: { id }, set }) => {
      const file = db.select().from(books).where(eq(books.id, id)).get();

      if (file?.file && file.path) {
        const { stream, size } = (await getFile(file.file, file.path)) ?? {};

        const filename = [
          file.authors ? slugify(file.authors) : null,
          file.title ? slugify(file.title) : null,
          file.path,
        ]
          .filter(Boolean)
          .join(".");

        set.headers[
          "content-disposition"
        ] = `attachment; filename="${filename}"`;
        set.headers["content-type"] = "text/fb2+xml";

        if (size) {
          set.headers["content-length"] = String(size);
        }

        return stream;
      }

      set.status = 404;
      return { ok: false };
    },
    { params: t.Object({ id: t.Number() }) }
  )
  .listen(port);
