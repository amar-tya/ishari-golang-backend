BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.chapters_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.chapters (
    id integer NOT NULL DEFAULT nextval('chapters_id_seq'::regclass),
    book_id integer NOT NULL,
    chapter_number integer NOT NULL,
    title varchar(255) NOT NULL,
    category varchar(50) NOT NULL,
    description text,
    total_verses integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT chapters_pkey PRIMARY KEY (id),
    CONSTRAINT unique_chapter_number UNIQUE (book_id, chapter_number),
    CONSTRAINT chapters_category_check CHECK (((category)::text = ANY ((ARRAY['Diwan'::character varying, 'Syaraful Anam'::character varying])::text[])))
);

ALTER SEQUENCE public.chapters_id_seq OWNED BY public.chapters.id;

-- Foreign Keys
ALTER TABLE public.chapters
    ADD CONSTRAINT chapters_book_id_fkey
    FOREIGN KEY (book_id) REFERENCES public.books (id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON public.chapters USING btree (book_id);
CREATE INDEX IF NOT EXISTS idx_chapters_category ON public.chapters USING btree (category);
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON public.chapters USING btree (deleted_at);

COMMIT;
