BEGIN;

CREATE TABLE IF NOT EXISTS public.hadi (
    id SERIAL PRIMARY KEY,
    name character varying(255) NOT NULL,
    description text,
    image_url text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

ALTER TABLE public.verse_media 
    ADD COLUMN IF NOT EXISTS hadi_id integer;

ALTER TABLE public.verse_media
    ADD CONSTRAINT verse_media_hadi_id_fkey
    FOREIGN KEY (hadi_id) REFERENCES public.hadi (id)
    ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_hadi_deleted_at ON public.hadi USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_verse_media_hadi_id ON public.verse_media USING btree (hadi_id);

COMMIT;
