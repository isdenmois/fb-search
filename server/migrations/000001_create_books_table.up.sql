CREATE TABLE IF NOT EXISTS "books" (
	"id" text PRIMARY KEY,
	"title" text NOT NULL,
	"search" text NOT NULL,
	"authors" text,
	"series" text,
	"serno" text,
	"lang" text,
    "size" integer
);
CREATE INDEX IF NOT EXISTS "search_ru_idx" ON "books" USING gin (to_tsvector('russian', "search"));
CREATE INDEX IF NOT EXISTS "search_simple_idx" ON "books" USING gin (to_tsvector('simple', "search"));
