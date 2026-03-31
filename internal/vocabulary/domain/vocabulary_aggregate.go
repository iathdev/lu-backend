package domain

import "time"

// Vocabulary is the aggregate root for vocabulary content.
// Language-agnostic: uses generic field names (word, phonetic).
// Language-specific data goes into Metadata JSONB.
type Vocabulary struct {
	ID                 VocabularyID
	LanguageID         LanguageID
	ProficiencyLevelID ProficiencyLevelID
	Word               string
	Phonetic           string
	AudioURL           string
	ImageURL           string
	FrequencyRank      int
	Metadata           map[string]any
	Meanings           []VocabularyMeaning
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// VocabularyMeaning represents a meaning in a target language.
type VocabularyMeaning struct {
	ID           MeaningID
	VocabularyID VocabularyID
	LanguageID   LanguageID
	Meaning      string
	WordType     string
	IsPrimary    bool
	Offset       int
	Examples     []VocabularyExample
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// VocabularyExample is an example sentence for a meaning.
type VocabularyExample struct {
	ID           ExampleID
	MeaningID    MeaningID
	Sentence     string
	Phonetic     string
	Translations map[string]string
	AudioURL     string
	Offset       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// VocabularyParams carries raw input for creating or updating a Vocabulary.
type VocabularyParams struct {
	LanguageID         string
	ProficiencyLevelID string
	Word               string
	Phonetic           string
	AudioURL           string
	ImageURL           string
	FrequencyRank      int
	Metadata           map[string]any
	Meanings           []MeaningParams
}

// MeaningParams carries raw input for a meaning.
type MeaningParams struct {
	LanguageID string
	Meaning    string
	WordType   string
	IsPrimary  bool
	Examples   []ExampleParams
}

// ExampleParams carries raw input for an example.
type ExampleParams struct {
	Sentence     string
	Phonetic     string
	Translations map[string]string
	AudioURL     string
}

func NewVocabularyFromParams(params VocabularyParams) (*Vocabulary, error) {
	if params.Word == "" {
		return nil, ErrWordRequired
	}
	if len(params.Meanings) == 0 {
		return nil, ErrMeaningRequired
	}

	langID, err := ParseLanguageID(params.LanguageID)
	if err != nil {
		return nil, ErrWordRequired
	}

	var plID ProficiencyLevelID
	if params.ProficiencyLevelID != "" {
		plID, err = ParseProficiencyLevelID(params.ProficiencyLevelID)
		if err != nil {
			return nil, err
		}
	}

	vocabID := NewVocabularyID()
	meanings := make([]VocabularyMeaning, 0, len(params.Meanings))
	for idx, mp := range params.Meanings {
		mLangID, parseErr := ParseLanguageID(mp.LanguageID)
		if parseErr != nil {
			return nil, parseErr
		}

		meaningID := NewMeaningID()
		examples := make([]VocabularyExample, 0, len(mp.Examples))
		for exIdx, ep := range mp.Examples {
			examples = append(examples, VocabularyExample{
				ID:           NewExampleID(),
				MeaningID:    meaningID,
				Sentence:     ep.Sentence,
				Phonetic:     ep.Phonetic,
				Translations: ep.Translations,
				AudioURL:     ep.AudioURL,
				Offset:       exIdx,
			})
		}

		meanings = append(meanings, VocabularyMeaning{
			ID:           meaningID,
			VocabularyID: vocabID,
			LanguageID:   mLangID,
			Meaning:      mp.Meaning,
			WordType:     mp.WordType,
			IsPrimary:    mp.IsPrimary,
			Offset:       idx,
			Examples:     examples,
		})
	}

	return &Vocabulary{
		ID:                 vocabID,
		LanguageID:         langID,
		ProficiencyLevelID: plID,
		Word:               params.Word,
		Phonetic:           params.Phonetic,
		AudioURL:           params.AudioURL,
		ImageURL:           params.ImageURL,
		FrequencyRank:      params.FrequencyRank,
		Metadata:           params.Metadata,
		Meanings:           meanings,
	}, nil
}
