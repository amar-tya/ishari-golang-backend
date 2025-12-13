BEGIN;

-- Add deleted_at to all tables (idempotent)
ALTER TABLE public.books ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON public.books USING btree (deleted_at);

ALTER TABLE public.chapters ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON public.chapters USING btree (deleted_at);

ALTER TABLE public.verses ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_verses_deleted_at ON public.verses USING btree (deleted_at);

ALTER TABLE public.bookmarks ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_bookmarks_deleted_at ON public.bookmarks USING btree (deleted_at);

ALTER TABLE public.translations ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_translations_deleted_at ON public.translations USING btree (deleted_at);

ALTER TABLE public.verse_media ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_verse_media_deleted_at ON public.verse_media USING btree (deleted_at);

ALTER TABLE public.search_history ADD COLUMN IF NOT EXISTS deleted_at timestamp without time zone;
CREATE INDEX IF NOT EXISTS idx_search_history_deleted_at ON public.search_history USING btree (deleted_at);

COMMIT;
