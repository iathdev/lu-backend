CREATE TABLE languages (
    id          UUID PRIMARY KEY,
    code        VARCHAR(10) NOT NULL UNIQUE,
    name_en     VARCHAR(100) NOT NULL,
    name_native VARCHAR(100) NOT NULL,
    is_active   BOOLEAN DEFAULT true,
    config      JSONB DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);
