package model

import (
	"learning-go/internal/shared/common"
	"learning-go/internal/vocabulary/domain"
	"time"

	"github.com/google/uuid"
)

type MeaningModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	VocabularyID uuid.UUID `gorm:"type:uuid;not null"`
	LanguageID   uuid.UUID `gorm:"type:uuid;not null"`
	Meaning      string    `gorm:"not null"`
	WordType     string
	IsPrimary    bool `gorm:"default:false"`
	Offset       int  `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (MeaningModel) TableName() string { return "vocabulary_meanings" }

func (model *MeaningModel) ToEntity() *domain.VocabularyMeaning {
	return &domain.VocabularyMeaning{
		ID:           domain.MeaningIDFromUUID(model.ID),
		VocabularyID: domain.VocabularyIDFromUUID(model.VocabularyID),
		LanguageID:   domain.LanguageIDFromUUID(model.LanguageID),
		Meaning:      model.Meaning,
		WordType:     model.WordType,
		IsPrimary:    model.IsPrimary,
		Offset:       model.Offset,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func FromMeaningEntity(meaning *domain.VocabularyMeaning) *MeaningModel {
	return &MeaningModel{
		ID:           meaning.ID.UUID(),
		VocabularyID: meaning.VocabularyID.UUID(),
		LanguageID:   meaning.LanguageID.UUID(),
		Meaning:      meaning.Meaning,
		WordType:     meaning.WordType,
		IsPrimary:    meaning.IsPrimary,
		Offset:       meaning.Offset,
		CreatedAt:    meaning.CreatedAt,
		UpdatedAt:    meaning.UpdatedAt,
	}
}

type TranslationsJSON = common.JSONB[map[string]string]

type ExampleModel struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key"`
	MeaningID    uuid.UUID        `gorm:"type:uuid;not null"`
	Sentence     string           `gorm:"not null"`
	Phonetic     string
	Translations TranslationsJSON `gorm:"type:jsonb;default:'{}'"`
	AudioURL     string
	Offset       int `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (ExampleModel) TableName() string { return "vocabulary_examples" }

func (model *ExampleModel) ToEntity() *domain.VocabularyExample {
	translations := make(map[string]string)
	for key, val := range model.Translations.Data {
		translations[key] = val
	}

	return &domain.VocabularyExample{
		ID:           domain.ExampleIDFromUUID(model.ID),
		MeaningID:    domain.MeaningIDFromUUID(model.MeaningID),
		Sentence:     model.Sentence,
		Phonetic:     model.Phonetic,
		Translations: translations,
		AudioURL:     model.AudioURL,
		Offset:       model.Offset,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func FromExampleEntity(example *domain.VocabularyExample) *ExampleModel {
	return &ExampleModel{
		ID:           example.ID.UUID(),
		MeaningID:    example.MeaningID.UUID(),
		Sentence:     example.Sentence,
		Phonetic:     example.Phonetic,
		Translations: common.NewJSONB(example.Translations),
		AudioURL:     example.AudioURL,
		Offset:       example.Offset,
		CreatedAt:    example.CreatedAt,
		UpdatedAt:    example.UpdatedAt,
	}
}
