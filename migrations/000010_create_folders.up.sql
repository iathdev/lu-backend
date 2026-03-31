CREATE TABLE folders (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL,
    language_id UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_folders_user ON folders(user_id, language_id);
