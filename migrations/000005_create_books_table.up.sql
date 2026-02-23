BEGIN;

-- Table
CREATE TABLE IF NOT EXISTS public.books (
    id SERIAL PRIMARY KEY,
    title varchar(255) NOT NULL,
    author varchar(255),
    description text,
    published_year integer,
    cover_image_url text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON public.books USING btree (deleted_at);

COMMIT;
