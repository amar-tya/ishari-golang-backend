BEGIN;

DROP TRIGGER IF EXISTS audit_users_changes ON public.users;
DROP TRIGGER IF EXISTS audit_books_changes ON public.books;
DROP TRIGGER IF EXISTS audit_chapters_changes ON public.chapters;
DROP TRIGGER IF EXISTS audit_verses_changes ON public.verses;
DROP TRIGGER IF EXISTS audit_translations_changes ON public.translations;

DROP FUNCTION IF EXISTS public.log_table_changes();
DROP TABLE IF EXISTS public.audit_logs;

COMMIT;
