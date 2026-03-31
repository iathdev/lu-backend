package port

import (
	"context"
	"learning-go/internal/shared/dto"
	vdto "learning-go/internal/vocabulary/application/dto"
)

// Input ports (driving) — used by handlers to call usecases

// --- Language Ports ---

type LanguageQueryPort interface {
	ListLanguages(ctx context.Context, activeOnly bool) ([]*vdto.LanguageResponse, error)
	GetLanguage(ctx context.Context, id string) (*vdto.LanguageResponse, error)
}

// --- Category Ports ---

type CategoryQueryPort interface {
	ListCategories(ctx context.Context, languageID string, isPublic *bool) ([]*vdto.CategoryResponse, error)
	GetCategory(ctx context.Context, id string) (*vdto.CategoryResponse, error)
}

// --- Proficiency Level Ports ---

type ProficiencyLevelQueryPort interface {
	ListProficiencyLevels(ctx context.Context, categoryID string) ([]*vdto.ProficiencyLevelResponse, error)
	GetProficiencyLevel(ctx context.Context, id string) (*vdto.ProficiencyLevelResponse, error)
}

// --- Vocabulary Ports ---

type VocabularyCommandPort interface {
	CreateVocabulary(ctx context.Context, req vdto.CreateVocabularyRequest) (*vdto.VocabularyResponse, error)
	UpdateVocabulary(ctx context.Context, id string, req vdto.UpdateVocabularyRequest) (*vdto.VocabularyResponse, error)
	DeleteVocabulary(ctx context.Context, id string) error
	SetTopics(ctx context.Context, id string, topicIDs []string) error
	SetGrammarPoints(ctx context.Context, id string, grammarPointIDs []string) error
}

type VocabularyQueryPort interface {
	GetVocabulary(ctx context.Context, id string, meaningLang string) (*vdto.VocabularyResponse, error)
	GetVocabularyDetail(ctx context.Context, id string, meaningLang string) (*vdto.VocabularyDetailResponse, error)
	ListVocabularies(ctx context.Context, filter vdto.VocabularyFilter, pagination dto.PaginationRequest) (*dto.ListResult[*vdto.VocabularyListResponse], error)
	SearchVocabulary(ctx context.Context, query string, languageID string, meaningLang string, pagination dto.PaginationRequest) (*dto.ListResult[*vdto.VocabularyListResponse], error)
}

// --- Topic Ports ---

type TopicQueryPort interface {
	ListTopics(ctx context.Context, categoryID string) ([]*vdto.TopicResponse, error)
	GetTopic(ctx context.Context, id string) (*vdto.TopicResponse, error)
}

// --- Grammar Point Ports ---

type GrammarPointQueryPort interface {
	ListGrammarPoints(ctx context.Context, categoryID string, proficiencyLevelID string, pagination dto.PaginationRequest) (*dto.ListResult[*vdto.GrammarPointResponse], error)
	GetGrammarPoint(ctx context.Context, id string) (*vdto.GrammarPointResponse, error)
}

// --- OCR Ports ---

type OCRScanCommandPort interface {
	ProcessOCRScan(ctx context.Context, req OCRScanInput) (*vdto.OCRScanResponse, error)
}

// --- Import Ports ---

type ImportCommandPort interface {
	ImportVocabularies(ctx context.Context, req vdto.BulkImportRequest) (*vdto.BulkImportResponse, error)
}

// --- Folder Ports ---

type FolderCommandPort interface {
	CreateFolder(ctx context.Context, userID string, req vdto.CreateFolderRequest) (*vdto.FolderResponse, error)
	UpdateFolder(ctx context.Context, id string, userID string, req vdto.UpdateFolderRequest) (*vdto.FolderResponse, error)
	DeleteFolder(ctx context.Context, id string, userID string) error
	AddVocabulary(ctx context.Context, folderID string, vocabID string, userID string) error
	RemoveVocabulary(ctx context.Context, folderID string, vocabID string, userID string) error
}

type FolderQueryPort interface {
	ListFolders(ctx context.Context, userID string, languageID string) ([]*vdto.FolderResponse, error)
	ListVocabularies(ctx context.Context, folderID string, userID string, meaningLang string, pagination dto.PaginationRequest) (*dto.ListResult[*vdto.VocabularyListResponse], error)
}
