CREATE TABLE vocabulary_meanings (
    id            UUID PRIMARY KEY,
    vocabulary_id UUID NOT NULL,
    language_id   UUID NOT NULL,
    meaning       TEXT NOT NULL,
    word_type     VARCHAR(20),
    is_primary    BOOLEAN DEFAULT false,
    "offset"      INTEGER DEFAULT 0,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(vocabulary_id, language_id, "offset")
);

CREATE INDEX idx_vm_vocab ON vocabulary_meanings(vocabulary_id);
CREATE INDEX idx_vm_lang ON vocabulary_meanings(vocabulary_id, language_id);
