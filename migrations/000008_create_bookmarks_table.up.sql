BEGIN;

-- Table
CREATE TABLE IF NOT EXISTS public.bookmarks (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    verse_id integer NOT NULL,
    note text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

-- Foreign Keys
ALTER TABLE public.bookmarks
    ADD CONSTRAINT bookmarks_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES public.users (id)
    ON DELETE CASCADE;

ALTER TABLE public.bookmarks
    ADD CONSTRAINT bookmarks_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS unique_user_bookmark ON public.bookmarks (user_id, verse_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON public.bookmarks USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON public.bookmarks USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookmarks_verse_id ON public.bookmarks USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_deleted_at ON public.bookmarks USING btree (deleted_at);

COMMIT;
