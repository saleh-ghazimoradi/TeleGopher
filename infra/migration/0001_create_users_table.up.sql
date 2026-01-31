CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email citext UNIQUE NOT NULL,
    password TEXT NOT NULL,
    refresh_token_web TEXT,
    refresh_token_web_at TIMESTAMPTZ,
    refresh_token_mobile TEXT,
    refresh_token_mobile_at TIMESTAMPTZ,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);