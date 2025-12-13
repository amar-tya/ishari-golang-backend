BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.books_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.books (
    id integer NOT NULL DEFAULT nextval('books_id_seq'::regclass),
    title varchar(255) NOT NULL,
    author varchar(255),
    description text,
    published_year integer,
    cover_image_url text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT books_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON public.books USING btree (deleted_at);

COMMIT;
