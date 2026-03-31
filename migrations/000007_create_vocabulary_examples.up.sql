CREATE TABLE vocabulary_examples (
    id           UUID PRIMARY KEY,
    meaning_id   UUID NOT NULL,
    sentence     TEXT NOT NULL,
    phonetic     TEXT,
    translations JSONB DEFAULT '{}'::jsonb,
    audio_url    VARCHAR(500),
    "offset"     INTEGER DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ve_meaning ON vocabulary_examples(meaning_id);
