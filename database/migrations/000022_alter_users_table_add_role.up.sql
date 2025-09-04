-- Add role_id column to users table
ALTER TABLE users ADD COLUMN role_id INTEGER REFERENCES roles(id) ON DELETE SET NULL;

-- Update existing users based on is_admin field
-- Set admin role (ID 1) for admin users, customer role (ID 4) for regular users
UPDATE users SET role_id = CASE 
    WHEN is_admin = true THEN 1 
    ELSE 4 
END;

-- Remove is_admin column
ALTER TABLE users DROP COLUMN is_admin;

-- Add index for performance
CREATE INDEX idx_users_role_id ON users(role_id);