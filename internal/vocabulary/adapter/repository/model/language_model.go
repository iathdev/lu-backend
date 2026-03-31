package model

import (
	"learning-go/internal/shared/common"
	"learning-go/internal/vocabulary/domain"
	"time"

	"github.com/google/uuid"
)

type LanguageModel struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key"`
	Code       string      `gorm:"not null;uniqueIndex"`
	NameEN     string      `gorm:"column:name_en;not null"`
	NameNative string      `gorm:"column:name_native;not null"`
	IsActive   bool        `gorm:"default:true"`
	Config     GenericJSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (LanguageModel) TableName() string { return "languages" }

func (model *LanguageModel) ToEntity() *domain.Language {
	config := make(map[string]any)
	for key, val := range model.Config.Data {
		config[key] = val
	}

	return &domain.Language{
		ID:         domain.LanguageIDFromUUID(model.ID),
		Code:       model.Code,
		NameEN:     model.NameEN,
		NameNative: model.NameNative,
		IsActive:   model.IsActive,
		Config:     config,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func FromLanguageEntity(lang *domain.Language) *LanguageModel {
	return &LanguageModel{
		ID:         lang.ID.UUID(),
		Code:       lang.Code,
		NameEN:     lang.NameEN,
		NameNative: lang.NameNative,
		IsActive:   lang.IsActive,
		Config:     common.NewJSONB(lang.Config),
		CreatedAt:  lang.CreatedAt,
		UpdatedAt:  lang.UpdatedAt,
	}
}
