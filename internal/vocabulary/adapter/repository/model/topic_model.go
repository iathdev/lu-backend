package model

import (
	"learning-go/internal/shared/common"
	"learning-go/internal/vocabulary/domain"
	"time"

	"github.com/google/uuid"
)

type NamesJSON = common.JSONB[map[string]string]

type TopicModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	Slug       string    `gorm:"not null"`
	Names      NamesJSON `gorm:"type:jsonb;not null"`
	Offset     int       `gorm:"column:offset;not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TopicModel) TableName() string { return "topics" }

func (model *TopicModel) ToEntity() *domain.Topic {
	names := make(map[string]string)
	for key, val := range model.Names.Data {
		names[key] = val
	}

	return &domain.Topic{
		ID:         domain.TopicIDFromUUID(model.ID),
		CategoryID: domain.CategoryIDFromUUID(model.CategoryID),
		Slug:       model.Slug,
		Names:      names,
		Offset:     model.Offset,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func FromTopicEntity(topic *domain.Topic) *TopicModel {
	return &TopicModel{
		ID:         topic.ID.UUID(),
		CategoryID: topic.CategoryID.UUID(),
		Slug:       topic.Slug,
		Names:      common.NewJSONB(topic.Names),
		Offset:     topic.Offset,
		CreatedAt:  topic.CreatedAt,
		UpdatedAt:  topic.UpdatedAt,
	}
}
