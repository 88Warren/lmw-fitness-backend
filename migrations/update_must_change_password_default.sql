-- Update the default value for MustChangePassword column
ALTER TABLE users ALTER COLUMN must_change_password SET DEFAULT false;

-- Update existing users who were created through manual registration to not require password change
-- (This assumes payment users will have auth_tokens, manual registration users won't)
UPDATE users 
SET must_change_password = false 
WHERE id NOT IN (
    SELECT DISTINCT user_id 
    FROM auth_tokens 
    WHERE user_id IS NOT NULL
);