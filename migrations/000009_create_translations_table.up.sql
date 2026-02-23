BEGIN;

-- Table
CREATE TABLE IF NOT EXISTS public.translations (
    id SERIAL PRIMARY KEY,
    verse_id integer NOT NULL,
    language_code varchar(10) NOT NULL,
    translation_text text NOT NULL,
    translator_name varchar(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

-- Foreign Keys
ALTER TABLE public.translations
    ADD CONSTRAINT translations_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS unique_translation ON public.translations (verse_id, language_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_translations_language ON public.translations USING btree (language_code);
CREATE INDEX IF NOT EXISTS idx_translations_text ON public.translations USING gin (to_tsvector('simple'::regconfig, translation_text));
CREATE INDEX IF NOT EXISTS idx_translations_verse_id ON public.translations USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_translations_deleted_at ON public.translations USING btree (deleted_at);

COMMIT;
