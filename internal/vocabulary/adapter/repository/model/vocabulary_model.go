package model

import (
	"learning-go/internal/shared/common"
	"learning-go/internal/vocabulary/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetadataJSON = common.JSONB[map[string]any]

type VocabularyModel struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key"`
	LanguageID         uuid.UUID      `gorm:"type:uuid;not null"`
	ProficiencyLevelID uuid.UUID      `gorm:"type:uuid"`
	Word               string         `gorm:"not null"`
	Phonetic           string
	AudioURL           string
	ImageURL           string
	FrequencyRank      int
	Metadata           MetadataJSON   `gorm:"type:jsonb;default:'{}'"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (VocabularyModel) TableName() string { return "vocabularies" }

func (model *VocabularyModel) ToEntity() *domain.Vocabulary {
	metadata := make(map[string]any)
	for key, val := range model.Metadata.Data {
		metadata[key] = val
	}

	return &domain.Vocabulary{
		ID:                 domain.VocabularyIDFromUUID(model.ID),
		LanguageID:         domain.LanguageIDFromUUID(model.LanguageID),
		ProficiencyLevelID: domain.ProficiencyLevelIDFromUUID(model.ProficiencyLevelID),
		Word:               model.Word,
		Phonetic:           model.Phonetic,
		AudioURL:           model.AudioURL,
		ImageURL:           model.ImageURL,
		FrequencyRank:      model.FrequencyRank,
		Metadata:           metadata,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}

func FromVocabularyEntity(vocab *domain.Vocabulary) *VocabularyModel {
	return &VocabularyModel{
		ID:                 vocab.ID.UUID(),
		LanguageID:         vocab.LanguageID.UUID(),
		ProficiencyLevelID: vocab.ProficiencyLevelID.UUID(),
		Word:               vocab.Word,
		Phonetic:           vocab.Phonetic,
		AudioURL:           vocab.AudioURL,
		ImageURL:           vocab.ImageURL,
		FrequencyRank:      vocab.FrequencyRank,
		Metadata:           common.NewJSONB(vocab.Metadata),
		CreatedAt:          vocab.CreatedAt,
		UpdatedAt:          vocab.UpdatedAt,
	}
}
