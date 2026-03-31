CREATE TABLE topics (
    id          UUID PRIMARY KEY,
    category_id UUID NOT NULL,
    slug        VARCHAR(100) NOT NULL,
    names       JSONB NOT NULL,
    "offset"    INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(category_id, slug)
);

CREATE INDEX idx_topics_category ON topics(category_id);
