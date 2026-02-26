BEGIN;

-- Table
CREATE TABLE IF NOT EXISTS public.chapters (
    id SERIAL PRIMARY KEY,
    book_id integer NOT NULL,
    chapter_number integer NOT NULL,
    title varchar(255) NOT NULL,
    category varchar(50) NOT NULL,
    description text,
    total_verses integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT chapters_category_check CHECK (((category)::text = ANY ((ARRAY['Diwan'::character varying, 'Syaraful Anam'::character varying, 'Muhud'::character varying, 'Rowi'::character varying, 'Diba'::character varying, 'Muradah'::character varying])::text[])))
);

-- Foreign Keys
ALTER TABLE public.chapters
    ADD CONSTRAINT chapters_book_id_fkey
    FOREIGN KEY (book_id) REFERENCES public.books (id)
    ON DELETE CASCADE;

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS unique_chapter_number ON public.chapters (book_id, chapter_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON public.chapters USING btree (book_id);
CREATE INDEX IF NOT EXISTS idx_chapters_category ON public.chapters USING btree (category);
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON public.chapters USING btree (deleted_at);

COMMIT;
