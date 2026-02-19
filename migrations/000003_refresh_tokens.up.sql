CREATE TABLE IF NOT EXISTS
    refresh_tokens (
        id TEXT PRIMARY KEY,
        owner_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        token_hash BYTEA UNIQUE NOT NULL,
        expires_at TIMESTAMPTZ NOT NULL,
        revoked_at TIMESTAMPTZ DEFAULT NULL,
        created_at TIMESTAMPTZ NOT NULL
    );