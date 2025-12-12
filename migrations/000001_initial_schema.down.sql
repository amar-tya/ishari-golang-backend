-- Reversible down migration: drops tables and sequences in reverse dependency order.

BEGIN;

DROP TABLE IF EXISTS public.verse_media;
DROP TABLE IF EXISTS public.translations;
DROP TABLE IF EXISTS public.bookmarks;
DROP TABLE IF EXISTS public.verses;
DROP TABLE IF EXISTS public.chapters;
DROP TABLE IF EXISTS public.books;
DROP TABLE IF EXISTS public.search_history;

DROP SEQUENCE IF EXISTS public.verse_media_id_seq;
DROP SEQUENCE IF EXISTS public.translations_id_seq;
DROP SEQUENCE IF EXISTS public.bookmarks_id_seq;
DROP SEQUENCE IF EXISTS public.verses_id_seq;
DROP SEQUENCE IF EXISTS public.chapters_id_seq;
DROP SEQUENCE IF EXISTS public.books_id_seq;
DROP SEQUENCE IF EXISTS public.search_history_id_seq;

COMMIT;
