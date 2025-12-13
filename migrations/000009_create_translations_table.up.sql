BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.translations_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.translations (
    id integer NOT NULL DEFAULT nextval('translations_id_seq'::regclass),
    verse_id integer NOT NULL,
    language_code varchar(10) NOT NULL,
    translation_text text NOT NULL,
    translator_name varchar(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT translations_pkey PRIMARY KEY (id),
    CONSTRAINT unique_translation UNIQUE (verse_id, language_code)
);

ALTER SEQUENCE public.translations_id_seq OWNED BY public.translations.id;

-- Foreign Keys
ALTER TABLE public.translations
    ADD CONSTRAINT translations_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_translations_language ON public.translations USING btree (language_code);
CREATE INDEX IF NOT EXISTS idx_translations_text ON public.translations USING gin (to_tsvector('simple'::regconfig, translation_text));
CREATE INDEX IF NOT EXISTS idx_translations_verse_id ON public.translations USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_translations_deleted_at ON public.translations USING btree (deleted_at);

COMMIT;
