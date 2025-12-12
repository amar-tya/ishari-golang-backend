BEGIN;

-- Sequences
CREATE SEQUENCE IF NOT EXISTS public.books_id_seq;
CREATE SEQUENCE IF NOT EXISTS public.chapters_id_seq;
CREATE SEQUENCE IF NOT EXISTS public.verses_id_seq;
CREATE SEQUENCE IF NOT EXISTS public.bookmarks_id_seq;
CREATE SEQUENCE IF NOT EXISTS public.translations_id_seq;
CREATE SEQUENCE IF NOT EXISTS public.verse_media_id_seq;
CREATE SEQUENCE IF NOT EXISTS public.search_history_id_seq;

-- Tables
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

CREATE TABLE IF NOT EXISTS public.verses (
    id integer NOT NULL DEFAULT nextval('verses_id_seq'::regclass),
    chapter_id integer NOT NULL,
    verse_number integer NOT NULL,
    arabic_text text NOT NULL,
    transliteration text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT verses_pkey PRIMARY KEY (id),
    CONSTRAINT unique_verse_number UNIQUE (chapter_id, verse_number)
);
ALTER SEQUENCE public.verses_id_seq OWNED BY public.verses.id;

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

CREATE TABLE IF NOT EXISTS public.verse_media (
    id integer NOT NULL DEFAULT nextval('verse_media_id_seq'::regclass),
    verse_id integer NOT NULL,
    media_type varchar(20) NOT NULL,
    media_url text NOT NULL,
    file_size integer,
    duration integer,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT verse_media_pkey PRIMARY KEY (id),
    CONSTRAINT unique_verse_media UNIQUE (verse_id, media_type, media_url),
    CONSTRAINT verse_media_media_type_check CHECK (((media_type)::text = ANY ((ARRAY['audio'::character varying, 'image'::character varying])::text[])))
);
ALTER SEQUENCE public.verse_media_id_seq OWNED BY public.verse_media.id;

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

-- Foreign Keys
ALTER TABLE public.chapters
    ADD CONSTRAINT chapters_book_id_fkey
    FOREIGN KEY (book_id) REFERENCES public.books (id)
    ON DELETE CASCADE;

ALTER TABLE public.verses
    ADD CONSTRAINT verses_chapter_id_fkey
    FOREIGN KEY (chapter_id) REFERENCES public.chapters (id)
    ON DELETE CASCADE;

ALTER TABLE public.bookmarks
    ADD CONSTRAINT bookmarks_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

ALTER TABLE public.translations
    ADD CONSTRAINT translations_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

ALTER TABLE public.verse_media
    ADD CONSTRAINT verse_media_verse_id_fkey
    FOREIGN KEY (verse_id) REFERENCES public.verses (id)
    ON DELETE CASCADE;

-- Non-unique Indexes
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON public.bookmarks USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookmarks_verse_id ON public.bookmarks USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_deleted_at ON public.bookmarks USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON public.chapters USING btree (book_id);
CREATE INDEX IF NOT EXISTS idx_chapters_category ON public.chapters USING btree (category);
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON public.chapters USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_search_history_searched_at ON public.search_history USING btree (searched_at DESC);
CREATE INDEX IF NOT EXISTS idx_search_history_deleted_at ON public.search_history USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_translations_language ON public.translations USING btree (language_code);
CREATE INDEX IF NOT EXISTS idx_translations_text ON public.translations USING gin (to_tsvector('simple'::regconfig, translation_text));
CREATE INDEX IF NOT EXISTS idx_translations_verse_id ON public.translations USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_translations_deleted_at ON public.translations USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_verse_media_type ON public.verse_media USING btree (media_type);
CREATE INDEX IF NOT EXISTS idx_verse_media_verse_id ON public.verse_media USING btree (verse_id);
CREATE INDEX IF NOT EXISTS idx_verse_media_deleted_at ON public.verse_media USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_verses_arabic_text ON public.verses USING gin (to_tsvector('simple'::regconfig, arabic_text));
CREATE INDEX IF NOT EXISTS idx_verses_chapter_id ON public.verses USING btree (chapter_id);
CREATE INDEX IF NOT EXISTS idx_verses_deleted_at ON public.verses USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON public.books USING btree (deleted_at);

COMMIT;
