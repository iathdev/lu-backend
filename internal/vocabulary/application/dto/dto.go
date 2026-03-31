package dto

import "time"

// --- Language DTOs ---

type LanguageResponse struct {
	ID         string         `json:"id"`
	Code       string         `json:"code"`
	NameEN     string         `json:"name_en"`
	NameNative string         `json:"name_native"`
	IsActive   bool           `json:"is_active"`
	Config     map[string]any `json:"config,omitempty"`
}

// --- Category DTOs ---

type CategoryResponse struct {
	ID         string `json:"id"`
	LanguageID string `json:"language_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	IsPublic   bool   `json:"is_public"`
}

// --- Proficiency Level DTOs ---

type ProficiencyLevelResponse struct {
	ID            string  `json:"id"`
	CategoryID    string  `json:"category_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Target        float64 `json:"target,omitempty"`
	DisplayTarget string  `json:"display_target,omitempty"`
	Offset        int     `json:"offset"`
}

// --- Vocabulary DTOs ---

type MeaningExampleDTO struct {
	Sentence     string            `json:"sentence" binding:"required"`
	Phonetic     string            `json:"phonetic,omitempty"`
	Translations map[string]string `json:"translations,omitempty"`
	AudioURL     string            `json:"audio_url,omitempty"`
}

type MeaningDTO struct {
	LanguageID string              `json:"language_id" binding:"required"`
	Meaning    string              `json:"meaning" binding:"required"`
	WordType   string              `json:"word_type,omitempty"`
	IsPrimary  bool                `json:"is_primary,omitempty"`
	Examples   []MeaningExampleDTO `json:"examples,omitempty"`
}

type CreateVocabularyRequest struct {
	LanguageID         string         `json:"language_id" binding:"required"`
	ProficiencyLevelID string         `json:"proficiency_level_id,omitempty"`
	Word               string         `json:"word" binding:"required"`
	Phonetic           string         `json:"phonetic,omitempty"`
	AudioURL           string         `json:"audio_url,omitempty"`
	ImageURL           string         `json:"image_url,omitempty"`
	FrequencyRank      int            `json:"frequency_rank,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
	Meanings           []MeaningDTO   `json:"meanings" binding:"required,min=1"`
}

type UpdateVocabularyRequest struct {
	LanguageID         string         `json:"language_id" binding:"required"`
	ProficiencyLevelID string         `json:"proficiency_level_id,omitempty"`
	Word               string         `json:"word" binding:"required"`
	Phonetic           string         `json:"phonetic,omitempty"`
	AudioURL           string         `json:"audio_url,omitempty"`
	ImageURL           string         `json:"image_url,omitempty"`
	FrequencyRank      int            `json:"frequency_rank,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
	Meanings           []MeaningDTO   `json:"meanings" binding:"required,min=1"`
	TopicIDs           []string       `json:"topic_ids,omitempty"`
	GrammarPointIDs    []string       `json:"grammar_point_ids,omitempty"`
}

type MeaningExampleResponse struct {
	ID           string            `json:"id"`
	Sentence     string            `json:"sentence"`
	Phonetic     string            `json:"phonetic,omitempty"`
	Translations map[string]string `json:"translations,omitempty"`
	AudioURL     string            `json:"audio_url,omitempty"`
}

type MeaningResponse struct {
	ID         string                   `json:"id"`
	LanguageID string                   `json:"language_id"`
	Meaning    string                   `json:"meaning"`
	WordType   string                   `json:"word_type,omitempty"`
	IsPrimary  bool                     `json:"is_primary,omitempty"`
	Offset     int                      `json:"offset"`
	Examples   []MeaningExampleResponse `json:"examples,omitempty"`
}

type VocabularyResponse struct {
	ID                 string            `json:"id"`
	LanguageID         string            `json:"language_id"`
	ProficiencyLevelID string            `json:"proficiency_level_id,omitempty"`
	Word               string            `json:"word"`
	Phonetic           string            `json:"phonetic,omitempty"`
	AudioURL           string            `json:"audio_url,omitempty"`
	ImageURL           string            `json:"image_url,omitempty"`
	FrequencyRank      int               `json:"frequency_rank,omitempty"`
	Metadata           map[string]any    `json:"metadata,omitempty"`
	Meanings           []MeaningResponse `json:"meanings,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

type VocabularyDetailResponse struct {
	VocabularyResponse
	Topics        []TopicResponse        `json:"topics"`
	GrammarPoints []GrammarPointResponse `json:"grammar_points"`
}

// VocabularyListResponse is a lightweight version for list endpoints.
type VocabularyListResponse struct {
	ID                 string                `json:"id"`
	Word               string                `json:"word"`
	Phonetic           string                `json:"phonetic,omitempty"`
	Meanings           []MeaningListResponse `json:"meanings,omitempty"`
	ProficiencyLevelID string                `json:"proficiency_level_id,omitempty"`
	FrequencyRank      int                   `json:"frequency_rank,omitempty"`
}

// MeaningListResponse is lightweight — no examples.
type MeaningListResponse struct {
	Meaning   string `json:"meaning"`
	WordType  string `json:"word_type,omitempty"`
	IsPrimary bool   `json:"is_primary,omitempty"`
}

// VocabularyFilter carries query params for listing vocabularies.
type VocabularyFilter struct {
	LanguageID         string `form:"language_id"`
	ProficiencyLevelID string `form:"proficiency_level_id"`
	TopicID            string `form:"topic_id"`
	MeaningLang        string `form:"meaning_lang"`
}

// --- Topic DTOs ---

type TopicResponse struct {
	ID         string            `json:"id"`
	CategoryID string            `json:"category_id"`
	Slug       string            `json:"slug"`
	Names      map[string]string `json:"names"`
	Offset     int               `json:"offset"`
}

// --- Grammar Point DTOs ---

type GrammarPointResponse struct {
	ID                 string         `json:"id"`
	CategoryID         string         `json:"category_id"`
	ProficiencyLevelID string         `json:"proficiency_level_id,omitempty"`
	Code               string         `json:"code"`
	Pattern            string         `json:"pattern"`
	Examples           map[string]any `json:"examples,omitempty"`
	Rule               map[string]any `json:"rule,omitempty"`
	CommonMistakes     map[string]any `json:"common_mistakes,omitempty"`
}

// --- OCR DTOs ---

type OCRScanHTTPRequest struct {
	Type     string `form:"type" binding:"omitempty,oneof=printed handwritten auto"`
	Language string `form:"language" binding:"omitempty"`
	Engine   string `form:"engine" binding:"omitempty,oneof=paddleocr tesseract google_vision baidu_ocr"`
}

type OCRScanCharacterItem struct {
	Text          string   `json:"text"`
	Pronunciation string   `json:"pronunciation,omitempty"`
	Confidence    float64  `json:"confidence"`
	Candidates    []string `json:"candidates,omitempty"`
}

type OCRScanExistingItem struct {
	VocabularyListResponse
	Confidence float64 `json:"confidence"`
}

type OCRScanMetadata struct {
	EngineUsed       string `json:"engine_used"`
	TotalDetected    int    `json:"total_detected"`
	ProcessingTimeMs int64  `json:"processing_time_ms"`
}

type OCRScanResponse struct {
	NewItems      []OCRScanCharacterItem `json:"new_items"`
	ExistingItems []OCRScanExistingItem  `json:"existing_items"`
	LowConfidence []OCRScanCharacterItem `json:"low_confidence"`
	Metadata      OCRScanMetadata        `json:"metadata"`
}

// --- Bulk Import DTOs ---

type BulkImportRequest struct {
	Vocabularies []CreateVocabularyRequest `json:"vocabularies" binding:"required,min=1"`
}

type BulkImportResponse struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
	Total    int `json:"total"`
}

// --- Folder DTOs ---

type CreateFolderRequest struct {
	LanguageID  string `json:"language_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateFolderRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type FolderResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	LanguageID      string    `json:"language_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	VocabularyCount int       `json:"vocabulary_count,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type FolderVocabularyRequest struct {
	VocabularyID string `json:"vocabulary_id" binding:"required"`
}

// --- Association DTOs ---

type SetTopicsRequest struct {
	TopicIDs []string `json:"topic_ids" binding:"required"`
}

type SetGrammarPointsRequest struct {
	GrammarPointIDs []string `json:"grammar_point_ids" binding:"required"`
}
