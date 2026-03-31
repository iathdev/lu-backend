package usecase

import (
	"context"

	apperr "learning-go/internal/shared/error"
	vdto "learning-go/internal/vocabulary/application/dto"
	"learning-go/internal/vocabulary/application/mapper"
	"learning-go/internal/vocabulary/application/port"
	"learning-go/internal/vocabulary/domain"
)

type VocabularyCommand struct {
	vocabRepo   port.VocabularyRepositoryPort
	topicRepo   port.TopicRepositoryPort
	grammarRepo port.GrammarPointRepositoryPort
}

func NewVocabularyCommand(
	vocabRepo port.VocabularyRepositoryPort,
	topicRepo port.TopicRepositoryPort,
	grammarRepo port.GrammarPointRepositoryPort,
) port.VocabularyCommandPort {
	return &VocabularyCommand{vocabRepo: vocabRepo, topicRepo: topicRepo, grammarRepo: grammarRepo}
}

func (useCase *VocabularyCommand) CreateVocabulary(ctx context.Context, req vdto.CreateVocabularyRequest) (*vdto.VocabularyResponse, error) {
	params := toVocabularyParams(req.LanguageID, req.ProficiencyLevelID, req.Word, req.Phonetic, req.AudioURL, req.ImageURL, req.FrequencyRank, req.Metadata, req.Meanings)

	vocab, err := domain.NewVocabularyFromParams(params)
	if err != nil {
		return nil, mapVocabEntityError(err)
	}

	if err := useCase.vocabRepo.Save(ctx, vocab); err != nil {
		return nil, apperr.InternalServerError("vocabulary.save_failed", err)
	}

	return mapper.ToVocabularyResponse(vocab), nil
}

func (useCase *VocabularyCommand) UpdateVocabulary(ctx context.Context, id string, req vdto.UpdateVocabularyRequest) (*vdto.VocabularyResponse, error) {
	vocabID, err := domain.ParseVocabularyID(id)
	if err != nil {
		return nil, apperr.BadRequest("vocabulary.invalid_id")
	}

	existing, err := useCase.vocabRepo.FindByID(ctx, vocabID)
	if err != nil {
		return nil, apperr.InternalServerError("vocabulary.query_failed", err)
	}
	if existing == nil {
		return nil, apperr.NotFound("vocabulary.not_found")
	}

	params := toVocabularyParams(req.LanguageID, req.ProficiencyLevelID, req.Word, req.Phonetic, req.AudioURL, req.ImageURL, req.FrequencyRank, req.Metadata, req.Meanings)

	vocab, err := domain.NewVocabularyFromParams(params)
	if err != nil {
		return nil, mapVocabEntityError(err)
	}

	// Preserve original ID and timestamps
	vocab.ID = existing.ID
	vocab.CreatedAt = existing.CreatedAt
	vocab.UpdatedAt = existing.UpdatedAt

	if err := useCase.vocabRepo.Update(ctx, vocab); err != nil {
		return nil, apperr.InternalServerError("vocabulary.update_failed", err)
	}

	// Set topics if provided
	if req.TopicIDs != nil {
		topicIDs, parseErr := parseTopicIDs(req.TopicIDs)
		if parseErr != nil {
			return nil, apperr.BadRequest("vocabulary.invalid_topic_id")
		}
		found, findErr := useCase.topicRepo.FindByIDs(ctx, topicIDs)
		if findErr != nil {
			return nil, apperr.InternalServerError("topic.query_failed", findErr)
		}
		if len(found) != len(topicIDs) {
			return nil, apperr.BadRequest("vocabulary.invalid_topic_id")
		}
		if err := useCase.vocabRepo.SetTopics(ctx, vocab.ID, topicIDs); err != nil {
			return nil, apperr.InternalServerError("vocabulary.set_topics_failed", err)
		}
	}

	// Set grammar points if provided
	if req.GrammarPointIDs != nil {
		gpIDs, parseErr := parseGrammarPointIDs(req.GrammarPointIDs)
		if parseErr != nil {
			return nil, apperr.BadRequest("vocabulary.invalid_grammar_point_id")
		}
		found, findErr := useCase.grammarRepo.FindByIDs(ctx, gpIDs)
		if findErr != nil {
			return nil, apperr.InternalServerError("grammar_point.query_failed", findErr)
		}
		if len(found) != len(gpIDs) {
			return nil, apperr.BadRequest("vocabulary.invalid_grammar_point_id")
		}
		if err := useCase.vocabRepo.SetGrammarPoints(ctx, vocab.ID, gpIDs); err != nil {
			return nil, apperr.InternalServerError("vocabulary.set_grammar_points_failed", err)
		}
	}

	return mapper.ToVocabularyResponse(vocab), nil
}

func (useCase *VocabularyCommand) DeleteVocabulary(ctx context.Context, id string) error {
	vocabID, err := domain.ParseVocabularyID(id)
	if err != nil {
		return apperr.BadRequest("vocabulary.invalid_id")
	}

	vocab, err := useCase.vocabRepo.FindByID(ctx, vocabID)
	if err != nil {
		return apperr.InternalServerError("vocabulary.query_failed", err)
	}
	if vocab == nil {
		return apperr.NotFound("vocabulary.not_found")
	}

	if err := useCase.vocabRepo.Delete(ctx, vocabID); err != nil {
		return apperr.InternalServerError("vocabulary.delete_failed", err)
	}

	return nil
}

func (useCase *VocabularyCommand) SetTopics(ctx context.Context, id string, topicIDs []string) error {
	vocabID, err := domain.ParseVocabularyID(id)
	if err != nil {
		return apperr.BadRequest("vocabulary.invalid_id")
	}

	vocab, err := useCase.vocabRepo.FindByID(ctx, vocabID)
	if err != nil {
		return apperr.InternalServerError("vocabulary.query_failed", err)
	}
	if vocab == nil {
		return apperr.NotFound("vocabulary.not_found")
	}

	parsed, err := parseTopicIDs(topicIDs)
	if err != nil {
		return apperr.BadRequest("vocabulary.invalid_topic_id")
	}

	if len(parsed) > 0 {
		found, findErr := useCase.topicRepo.FindByIDs(ctx, parsed)
		if findErr != nil {
			return apperr.InternalServerError("topic.query_failed", findErr)
		}
		if len(found) != len(parsed) {
			return apperr.BadRequest("vocabulary.invalid_topic_id")
		}
	}

	if err := useCase.vocabRepo.SetTopics(ctx, vocabID, parsed); err != nil {
		return apperr.InternalServerError("vocabulary.set_topics_failed", err)
	}

	return nil
}

func (useCase *VocabularyCommand) SetGrammarPoints(ctx context.Context, id string, grammarPointIDs []string) error {
	vocabID, err := domain.ParseVocabularyID(id)
	if err != nil {
		return apperr.BadRequest("vocabulary.invalid_id")
	}

	vocab, err := useCase.vocabRepo.FindByID(ctx, vocabID)
	if err != nil {
		return apperr.InternalServerError("vocabulary.query_failed", err)
	}
	if vocab == nil {
		return apperr.NotFound("vocabulary.not_found")
	}

	parsed, err := parseGrammarPointIDs(grammarPointIDs)
	if err != nil {
		return apperr.BadRequest("vocabulary.invalid_grammar_point_id")
	}

	if len(parsed) > 0 {
		found, findErr := useCase.grammarRepo.FindByIDs(ctx, parsed)
		if findErr != nil {
			return apperr.InternalServerError("grammar_point.query_failed", findErr)
		}
		if len(found) != len(parsed) {
			return apperr.BadRequest("vocabulary.invalid_grammar_point_id")
		}
	}

	if err := useCase.vocabRepo.SetGrammarPoints(ctx, vocabID, parsed); err != nil {
		return apperr.InternalServerError("vocabulary.set_grammar_points_failed", err)
	}

	return nil
}

// toVocabularyParams converts DTO fields to domain VocabularyParams.
func toVocabularyParams(
	languageID, proficiencyLevelID, word, phonetic, audioURL, imageURL string,
	frequencyRank int,
	metadata map[string]any,
	meanings []vdto.MeaningDTO,
) domain.VocabularyParams {
	meaningParams := make([]domain.MeaningParams, 0, len(meanings))
	for _, meaning := range meanings {
		exampleParams := make([]domain.ExampleParams, 0, len(meaning.Examples))
		for _, example := range meaning.Examples {
			exampleParams = append(exampleParams, domain.ExampleParams{
				Sentence:     example.Sentence,
				Phonetic:     example.Phonetic,
				Translations: example.Translations,
				AudioURL:     example.AudioURL,
			})
		}

		meaningParams = append(meaningParams, domain.MeaningParams{
			LanguageID: meaning.LanguageID,
			Meaning:    meaning.Meaning,
			WordType:   meaning.WordType,
			IsPrimary:  meaning.IsPrimary,
			Examples:   exampleParams,
		})
	}

	return domain.VocabularyParams{
		LanguageID:         languageID,
		ProficiencyLevelID: proficiencyLevelID,
		Word:               word,
		Phonetic:           phonetic,
		AudioURL:           audioURL,
		ImageURL:           imageURL,
		FrequencyRank:      frequencyRank,
		Metadata:           metadata,
		Meanings:           meaningParams,
	}
}
