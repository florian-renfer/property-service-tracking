BEGIN;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL CHECK (
        LENGTH(TRIM(email)) > 0 AND email = LOWER(email)
    ),
    password_hash TEXT NOT NULL CHECK (LENGTH(TRIM(password_hash)) > 0),
    given_name VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(given_name)) > 0),
    family_name VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(family_name)) > 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (email)
);

CREATE INDEX users_active_created_at_idx ON users (active, created_at);

COMMIT;
