BEGIN;

-- Table
CREATE TABLE IF NOT EXISTS public.verses (
    id SERIAL PRIMARY KEY,
    chapter_id integer NOT NULL,
    verse_number integer NOT NULL,
    arabic_text text NOT NULL,
    transliteration text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

-- Foreign Keys
ALTER TABLE public.verses
    ADD CONSTRAINT verses_chapter_id_fkey
    FOREIGN KEY (chapter_id) REFERENCES public.chapters (id)
    ON DELETE CASCADE;

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS unique_verse_number ON public.verses (chapter_id, verse_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_verses_arabic_text ON public.verses USING gin (to_tsvector('simple'::regconfig, arabic_text));
CREATE INDEX IF NOT EXISTS idx_verses_chapter_id ON public.verses USING btree (chapter_id);
CREATE INDEX IF NOT EXISTS idx_verses_deleted_at ON public.verses USING btree (deleted_at);

COMMIT;
