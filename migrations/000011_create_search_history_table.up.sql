BEGIN;

-- Table
CREATE TABLE IF NOT EXISTS public.search_history (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    search_query text NOT NULL,
    search_type varchar(20) DEFAULT 'text',
    results_count integer DEFAULT 0,
    searched_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

-- Foreign Keys
ALTER TABLE public.search_history
    ADD CONSTRAINT search_history_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES public.users (id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_search_history_user_id ON public.search_history USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_search_history_searched_at ON public.search_history USING btree (searched_at DESC);
CREATE INDEX IF NOT EXISTS idx_search_history_deleted_at ON public.search_history USING btree (deleted_at);

COMMIT;
