CREATE TABLE vocabulary_topics (
    vocabulary_id UUID NOT NULL,
    topic_id      UUID NOT NULL,
    PRIMARY KEY (vocabulary_id, topic_id)
);

CREATE INDEX idx_vt_topic ON vocabulary_topics(topic_id);
