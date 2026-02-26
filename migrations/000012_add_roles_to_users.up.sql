BEGIN;

-- 1. Create User Role Enum Type
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('super_admin', 'admin_content', 'user');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- 2. Add Role column to users table with default 'user'
ALTER TABLE public.users ADD COLUMN IF NOT EXISTS role user_role NOT NULL DEFAULT 'user';

-- 3. Update the first registered user to be super_admin (Optional/Convenience)
-- UPDATE public.users SET role = 'super_admin' WHERE id = (SELECT id FROM public.users ORDER BY id ASC LIMIT 1);

COMMIT;
