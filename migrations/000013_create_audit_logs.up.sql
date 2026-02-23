BEGIN;

-- 1. Create the audit table
CREATE TABLE IF NOT EXISTS public.audit_logs (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL,
    record_id INTEGER NOT NULL,
    operation VARCHAR(10) NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    old_data JSONB,
    new_data JSONB,
    changed_by_user_id INTEGER REFERENCES public.users(id) ON DELETE SET NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for querying audit logs quickly
CREATE INDEX IF NOT EXISTS idx_audit_logs_table_name ON public.audit_logs (table_name);
CREATE INDEX IF NOT EXISTS idx_audit_logs_changed_at ON public.audit_logs (changed_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON public.audit_logs (changed_by_user_id);

-- 2. Create the Generic Trigger Function
CREATE OR REPLACE FUNCTION public.log_table_changes()
RETURNS TRIGGER AS $$
DECLARE
    aktor_id TEXT;
    r_id INTEGER;
BEGIN
    -- Get user id set by the application through SET LOCAL app.current_user_id
    aktor_id := current_setting('app.current_user_id', true);

    IF (TG_OP = 'INSERT') THEN
        EXECUTE 'SELECT ($1).id' INTO r_id USING NEW;
        INSERT INTO public.audit_logs (table_name, record_id, operation, new_data, changed_by_user_id)
        VALUES (TG_TABLE_NAME, r_id, 'INSERT', row_to_json(NEW)::jsonb, NULLIF(aktor_id, '')::integer);
        RETURN NEW;
        
    ELSIF (TG_OP = 'UPDATE') THEN
        EXECUTE 'SELECT ($1).id' INTO r_id USING NEW;
        -- Optional: Only insert if there is actual change
        IF row_to_json(OLD)::jsonb != row_to_json(NEW)::jsonb THEN
            INSERT INTO public.audit_logs (table_name, record_id, operation, old_data, new_data, changed_by_user_id)
            VALUES (TG_TABLE_NAME, r_id, 'UPDATE', row_to_json(OLD)::jsonb, row_to_json(NEW)::jsonb, NULLIF(aktor_id, '')::integer);
        END IF;
        RETURN NEW;
        
    ELSIF (TG_OP = 'DELETE') THEN
        EXECUTE 'SELECT ($1).id' INTO r_id USING OLD;
        INSERT INTO public.audit_logs (table_name, record_id, operation, old_data, changed_by_user_id)
        VALUES (TG_TABLE_NAME, r_id, 'DELETE', row_to_json(OLD)::jsonb, NULLIF(aktor_id, '')::integer);
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 3. Attach empty triggers to core tables
CREATE TRIGGER audit_users_changes
    AFTER INSERT OR UPDATE OR DELETE ON public.users
    FOR EACH ROW EXECUTE FUNCTION public.log_table_changes();

CREATE TRIGGER audit_books_changes
    AFTER INSERT OR UPDATE OR DELETE ON public.books
    FOR EACH ROW EXECUTE FUNCTION public.log_table_changes();

CREATE TRIGGER audit_chapters_changes
    AFTER INSERT OR UPDATE OR DELETE ON public.chapters
    FOR EACH ROW EXECUTE FUNCTION public.log_table_changes();

CREATE TRIGGER audit_verses_changes
    AFTER INSERT OR UPDATE OR DELETE ON public.verses
    FOR EACH ROW EXECUTE FUNCTION public.log_table_changes();

CREATE TRIGGER audit_translations_changes
    AFTER INSERT OR UPDATE OR DELETE ON public.translations
    FOR EACH ROW EXECUTE FUNCTION public.log_table_changes();

COMMIT;
