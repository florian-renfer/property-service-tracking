BEGIN;

CREATE TABLE roles (
    id UUID PRIMARY KEY,
    label VARCHAR(100) NOT NULL UNIQUE CHECK (LENGTH(TRIM(label)) > 0 AND label = UPPER(label)),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD COLUMN role_id UUID NOT NULL REFERENCES roles (id) ON DELETE RESTRICT;

CREATE INDEX users_role_id_idx ON users (role_id);

INSERT INTO roles (id, label, description) VALUES
    ('00000000-0000-0000-0000-000000000001', 'GLOBAL_ADMIN', 'Global administrator with full system access.'),
    ('00000000-0000-0000-0000-000000000002', 'PROPERTY_MANAGER', 'Property manager with access to assigned properties.');

COMMIT;
