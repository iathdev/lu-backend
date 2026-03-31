package repository

import (
	"context"
	"errors"
	"learning-go/internal/vocabulary/adapter/repository/model"
	"learning-go/internal/vocabulary/application/port"
	"learning-go/internal/vocabulary/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) port.FolderRepositoryPort {
	return &FolderRepository{db: db}
}

func (repo *FolderRepository) Save(ctx context.Context, folder *domain.Folder) error {
	m := model.FromFolderEntity(folder)
	if err := repo.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	folder.CreatedAt = m.CreatedAt
	folder.UpdatedAt = m.UpdatedAt
	return nil
}

func (repo *FolderRepository) FindByID(ctx context.Context, id domain.FolderID) (*domain.Folder, error) {
	var m model.FolderModel
	if err := repo.db.WithContext(ctx).First(&m, "id = ?", id.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

func (repo *FolderRepository) FindByUserID(ctx context.Context, userID domain.UserID, languageID *domain.LanguageID) ([]*domain.Folder, error) {
	query := repo.db.WithContext(ctx).Where("user_id = ?", userID.UUID())
	if languageID != nil {
		query = query.Where("language_id = ?", languageID.UUID())
	}

	var models []model.FolderModel
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Folder, 0, len(models))
	for _, m := range models {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (repo *FolderRepository) CountVocabulariesByFolderIDs(ctx context.Context, folderIDs []domain.FolderID) (map[domain.FolderID]int, error) {
	if len(folderIDs) == 0 {
		return map[domain.FolderID]int{}, nil
	}

	uuids := make([]uuid.UUID, 0, len(folderIDs))
	for _, id := range folderIDs {
		uuids = append(uuids, id.UUID())
	}

	type countRow struct {
		FolderID uuid.UUID
		Count    int
	}

	var rows []countRow
	if err := repo.db.WithContext(ctx).
		Model(&model.FolderVocabularyModel{}).
		Select("folder_id, COUNT(*) as count").
		Where("folder_id IN ?", uuids).
		Group("folder_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[domain.FolderID]int, len(rows))
	for _, row := range rows {
		result[domain.FolderIDFromUUID(row.FolderID)] = row.Count
	}
	return result, nil
}

func (repo *FolderRepository) Update(ctx context.Context, folder *domain.Folder) error {
	m := model.FromFolderEntity(folder)
	if err := repo.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}
	folder.UpdatedAt = m.UpdatedAt
	return nil
}

func (repo *FolderRepository) Delete(ctx context.Context, id domain.FolderID) error {
	return repo.db.WithContext(ctx).Delete(&model.FolderModel{}, "id = ?", id.UUID()).Error
}

func (repo *FolderRepository) AddVocabulary(ctx context.Context, folderID domain.FolderID, vocabID domain.VocabularyID) error {
	fv := model.FolderVocabularyModel{
		FolderID:     folderID.UUID(),
		VocabularyID: vocabID.UUID(),
		AddedAt:      time.Now(),
	}
	return repo.db.WithContext(ctx).Create(&fv).Error
}

func (repo *FolderRepository) RemoveVocabulary(ctx context.Context, folderID domain.FolderID, vocabID domain.VocabularyID) error {
	return repo.db.WithContext(ctx).
		Where("folder_id = ? AND vocabulary_id = ?", folderID.UUID(), vocabID.UUID()).
		Delete(&model.FolderVocabularyModel{}).Error
}

func (repo *FolderRepository) FindVocabularies(ctx context.Context, folderID domain.FolderID, offset, limit int) ([]*domain.Vocabulary, error) {
	var models []model.VocabularyModel
	if err := repo.db.WithContext(ctx).
		Joins("JOIN folder_vocabularies fv ON fv.vocabulary_id = vocabularies.id").
		Where("fv.folder_id = ?", folderID.UUID()).
		Offset(offset).Limit(limit).Order("fv.added_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	vocabs := toVocabEntities(models)

	// Use a temporary VocabularyRepository to reuse the loadMeanings helper.
	vocabRepo := &VocabularyRepository{db: repo.db}
	if err := vocabRepo.loadMeanings(ctx, vocabs); err != nil {
		return nil, err
	}

	return vocabs, nil
}

func (repo *FolderRepository) CountVocabularies(ctx context.Context, folderID domain.FolderID) (int64, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&model.FolderVocabularyModel{}).
		Where("folder_id = ?", folderID.UUID()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
