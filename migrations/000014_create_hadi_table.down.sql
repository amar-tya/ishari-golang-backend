BEGIN;

ALTER TABLE public.verse_media DROP CONSTRAINT IF EXISTS verse_media_hadi_id_fkey;
DROP INDEX IF EXISTS public.idx_verse_media_hadi_id;
ALTER TABLE public.verse_media DROP COLUMN IF EXISTS hadi_id;

DROP INDEX IF EXISTS public.idx_hadi_deleted_at;
DROP TABLE IF EXISTS public.hadi;

COMMIT;
