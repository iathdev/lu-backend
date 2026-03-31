package repository

import (
	"context"
	"errors"
	"learning-go/internal/vocabulary/adapter/repository/model"
	"learning-go/internal/vocabulary/application/port"
	"learning-go/internal/vocabulary/domain"

	"gorm.io/gorm"
)

type LanguageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) port.LanguageRepositoryPort {
	return &LanguageRepository{db: db}
}

func (repo *LanguageRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Language, error) {
	query := repo.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var models []model.LanguageModel
	if err := query.Order("code ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Language, 0, len(models))
	for _, m := range models {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (repo *LanguageRepository) FindByID(ctx context.Context, id domain.LanguageID) (*domain.Language, error) {
	var m model.LanguageModel
	if err := repo.db.WithContext(ctx).First(&m, "id = ?", id.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToEntity(), nil
}
