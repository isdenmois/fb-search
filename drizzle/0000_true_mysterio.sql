CREATE TABLE IF NOT EXISTS "books" (
	"id" integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY (sequence name "books_id_seq" INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START WITH 1 CACHE 1),
	"title" text NOT NULL,
	"search" text NOT NULL,
	"authors" text,
	"series" text,
	"serno" text,
	"file" text,
	"path" text,
	"lang" text
);
--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "search_ru_idx" ON "books" USING gin (to_tsvector('russian', "search"));--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "search_simple_idx" ON "books" USING gin (to_tsvector('simple', "search"));