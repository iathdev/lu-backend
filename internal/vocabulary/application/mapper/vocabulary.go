package mapper

import (
	"math"

	"learning-go/internal/shared/dto"
	vdto "learning-go/internal/vocabulary/application/dto"
	"learning-go/internal/vocabulary/domain"
)

// ToVocabularyResponse maps domain.Vocabulary to VocabularyResponse with Meanings[].Examples[].
func ToVocabularyResponse(vocab *domain.Vocabulary) *vdto.VocabularyResponse {
	meanings := make([]vdto.MeaningResponse, 0, len(vocab.Meanings))
	for _, meaning := range vocab.Meanings {
		examples := make([]vdto.MeaningExampleResponse, 0, len(meaning.Examples))
		for _, example := range meaning.Examples {
			examples = append(examples, vdto.MeaningExampleResponse{
				ID:           example.ID.String(),
				Sentence:     example.Sentence,
				Phonetic:     example.Phonetic,
				Translations: example.Translations,
				AudioURL:     example.AudioURL,
			})
		}

		meanings = append(meanings, vdto.MeaningResponse{
			ID:         meaning.ID.String(),
			LanguageID: meaning.LanguageID.String(),
			Meaning:    meaning.Meaning,
			WordType:   meaning.WordType,
			IsPrimary:  meaning.IsPrimary,
			Offset:     meaning.Offset,
			Examples:   examples,
		})
	}

	return &vdto.VocabularyResponse{
		ID:                 vocab.ID.String(),
		LanguageID:         vocab.LanguageID.String(),
		ProficiencyLevelID: vocab.ProficiencyLevelID.String(),
		Word:               vocab.Word,
		Phonetic:           vocab.Phonetic,
		AudioURL:           vocab.AudioURL,
		ImageURL:           vocab.ImageURL,
		FrequencyRank:      vocab.FrequencyRank,
		Metadata:           vocab.Metadata,
		Meanings:           meanings,
		CreatedAt:          vocab.CreatedAt,
	}
}

// ToVocabularyListResponse maps domain.Vocabulary to lightweight VocabularyListResponse (no examples).
func ToVocabularyListResponse(vocab *domain.Vocabulary) *vdto.VocabularyListResponse {
	meanings := make([]vdto.MeaningListResponse, 0, len(vocab.Meanings))
	for _, meaning := range vocab.Meanings {
		meanings = append(meanings, vdto.MeaningListResponse{
			Meaning:   meaning.Meaning,
			WordType:  meaning.WordType,
			IsPrimary: meaning.IsPrimary,
		})
	}

	return &vdto.VocabularyListResponse{
		ID:                 vocab.ID.String(),
		Word:               vocab.Word,
		Phonetic:           vocab.Phonetic,
		Meanings:           meanings,
		ProficiencyLevelID: vocab.ProficiencyLevelID.String(),
		FrequencyRank:      vocab.FrequencyRank,
	}
}

// ToPaginatedListResult maps a slice of domain vocabularies to a paginated ListResult of VocabularyListResponse.
func ToPaginatedListResult(vocabs []*domain.Vocabulary, total int64, pagination dto.PaginationRequest) *dto.ListResult[*vdto.VocabularyListResponse] {
	items := make([]*vdto.VocabularyListResponse, 0, len(vocabs))
	for _, vocab := range vocabs {
		items = append(items, ToVocabularyListResponse(vocab))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &dto.ListResult[*vdto.VocabularyListResponse]{
		Items:      items,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}
}
