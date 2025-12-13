BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.bookmarks_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.bookmarks (
    id integer NOT NULL DEFAULT nextval('bookmarks_id_seq'::regclass),
    verse_id integer NOT NULL,
    note text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT bookmarks_pkey PRIMARY KEY (id),
    CONSTRAINT unique_bookmark UNIQUE (verse_id)
);

ALTER SEQUENCE public.bookmarks_id_seq OWNED BY public.bookmarks.id;

-- Foreign Keys
ALTER TABLE public.bookmarks
    ADD CONSTRAINT bookmarks_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON public.bookmarks USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookmarks_verse_id ON public.bookmarks USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_deleted_at ON public.bookmarks USING btree (deleted_at);

COMMIT;
