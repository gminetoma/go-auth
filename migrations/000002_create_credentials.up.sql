CREATE TABLE IF NOT EXISTS
    credentials (
        id TEXT PRIMARY KEY,
        owner_ID TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        email TEXT UNIQUE NOT NULL,
        password_hash BYTEA NOT NULL,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
    );