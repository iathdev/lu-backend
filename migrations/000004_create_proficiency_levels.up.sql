CREATE TABLE proficiency_levels (
    id             UUID PRIMARY KEY,
    category_id    UUID NOT NULL,
    code           VARCHAR(20) NOT NULL,
    name           VARCHAR(100) NOT NULL,
    target         DECIMAL(8,2),
    display_target VARCHAR(255),
    "offset"       INTEGER NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(category_id, code)
);

CREATE INDEX idx_pl_category ON proficiency_levels(category_id, "offset");
