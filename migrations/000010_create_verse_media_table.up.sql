BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.verse_media_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.verse_media (
    id integer NOT NULL DEFAULT nextval('verse_media_id_seq'::regclass),
    verse_id integer NOT NULL,
    media_type varchar(20) NOT NULL,
    media_url text NOT NULL,
    file_size integer,
    duration integer,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT verse_media_pkey PRIMARY KEY (id),
    CONSTRAINT unique_verse_media UNIQUE (verse_id, media_type, media_url),
    CONSTRAINT verse_media_media_type_check CHECK (((media_type)::text = ANY ((ARRAY['audio'::character varying, 'image'::character varying])::text[])))
);

ALTER SEQUENCE public.verse_media_id_seq OWNED BY public.verse_media.id;

-- Foreign Keys
ALTER TABLE public.verse_media
    ADD CONSTRAINT verse_media_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_verse_media_type ON public.verse_media USING btree (media_type);
CREATE INDEX IF NOT EXISTS idx_verse_media_verse_id ON public.verse_media USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_verse_media_deleted_at ON public.verse_media USING btree (deleted_at);

COMMIT;
