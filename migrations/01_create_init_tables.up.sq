-- Drop the table if it exists
DROP TABLE IF EXISTS users CASCADE;


-- Create the users table
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL CHECK (last_name <> ''),
    email VARCHAR(64) UNIQUE NOT NULL CHECK (email <> ''),
    password VARCHAR(250) NOT NULL CHECK (octet_length(password) <> 0),
    role VARCHAR(20) DEFAULT 'client',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
