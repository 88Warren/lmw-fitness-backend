-- Create auth_tokens table manually
CREATE TABLE IF NOT EXISTS auth_tokens (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id INTEGER,
    token VARCHAR(64) UNIQUE,
    program_name VARCHAR(255),
    day_number INTEGER,
    is_used BOOLEAN DEFAULT false,
    session_id VARCHAR(255)
);

-- Create index on user_id
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_auth_tokens_deleted_at ON auth_tokens(deleted_at);