CREATE TABLE folder_vocabularies (
    folder_id     UUID NOT NULL,
    vocabulary_id UUID NOT NULL,
    added_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (folder_id, vocabulary_id)
);

CREATE INDEX idx_fv_vocabulary ON folder_vocabularies(vocabulary_id);
