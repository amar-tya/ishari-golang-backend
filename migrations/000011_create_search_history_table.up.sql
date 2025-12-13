BEGIN;

-- Sequence
CREATE SEQUENCE IF NOT EXISTS public.search_history_id_seq;

-- Table
CREATE TABLE IF NOT EXISTS public.search_history (
    id integer NOT NULL DEFAULT nextval('search_history_id_seq'::regclass),
    search_query text NOT NULL,
    search_type varchar(20) DEFAULT 'text',
    results_count integer DEFAULT 0,
    searched_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT search_history_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.search_history_id_seq OWNED BY public.search_history.id;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_search_history_searched_at ON public.search_history USING btree (searched_at DESC);
CREATE INDEX IF NOT EXISTS idx_search_history_deleted_at ON public.search_history USING btree (deleted_at);

COMMIT;
