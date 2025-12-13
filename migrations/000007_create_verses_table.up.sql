BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.verses_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.verses (
    id integer NOT NULL DEFAULT nextval('verses_id_seq'::regclass),
    chapter_id integer NOT NULL,
    verse_number integer NOT NULL,
    arabic_text text NOT NULL,
    transliteration text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT verses_pkey PRIMARY KEY (id),
    CONSTRAINT unique_verse_number UNIQUE (chapter_id, verse_number)
);

ALTER SEQUENCE public.verses_id_seq OWNED BY public.verses.id;

-- Foreign Keys
ALTER TABLE public.verses
    ADD CONSTRAINT verses_chapter_id_fkey
    FOREIGN KEY (chapter_id) REFERENCES public.chapters (id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_verses_arabic_text ON public.verses USING gin (to_tsvector('simple'::regconfig, arabic_text));
CREATE INDEX IF NOT EXISTS idx_verses_chapter_id ON public.verses USING btree (chapter_id);
CREATE INDEX IF NOT EXISTS idx_verses_deleted_at ON public.verses USING btree (deleted_at);

COMMIT;
