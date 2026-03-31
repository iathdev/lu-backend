package repository

import (
	"context"
	"errors"
	"learning-go/internal/vocabulary/adapter/repository/model"
	"learning-go/internal/vocabulary/application/port"
	"learning-go/internal/vocabulary/domain"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) port.CategoryRepositoryPort {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) FindAll(ctx context.Context, languageID *domain.LanguageID, isPublic *bool) ([]*domain.Category, error) {
	query := repo.db.WithContext(ctx)
	if languageID != nil {
		query = query.Where("language_id = ?", languageID.UUID())
	}
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	var models []model.CategoryModel
	if err := query.Order("code ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Category, 0, len(models))
	for _, m := range models {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (repo *CategoryRepository) FindByID(ctx context.Context, id domain.CategoryID) (*domain.Category, error) {
	var m model.CategoryModel
	if err := repo.db.WithContext(ctx).First(&m, "id = ?", id.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToEntity(), nil
}
