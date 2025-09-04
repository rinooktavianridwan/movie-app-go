-- Add back is_admin column
ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;

-- Update is_admin based on role_id (assuming admin role has ID 1)
UPDATE users SET is_admin = CASE 
    WHEN role_id = 1 THEN true 
    ELSE false 
END;

-- Drop role_id column
DROP INDEX IF EXISTS idx_users_role_id;
ALTER TABLE users DROP COLUMN role_id;