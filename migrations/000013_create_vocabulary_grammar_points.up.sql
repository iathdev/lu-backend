CREATE TABLE vocabulary_grammar_points (
    vocabulary_id    UUID NOT NULL,
    grammar_point_id UUID NOT NULL,
    PRIMARY KEY (vocabulary_id, grammar_point_id)
);

CREATE INDEX idx_vgp_gp ON vocabulary_grammar_points(grammar_point_id);
