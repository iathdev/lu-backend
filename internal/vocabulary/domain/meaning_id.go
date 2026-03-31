package domain

import "github.com/google/uuid"

// MeaningID uniquely identifies a VocabularyMeaning entity.
type MeaningID uuid.UUID

func NewMeaningID() MeaningID {
	return MeaningID(uuid.Must(uuid.NewV7()))
}

func ParseMeaningID(raw string) (MeaningID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return MeaningID{}, err
	}
	return MeaningID(parsed), nil
}

func MeaningIDFromUUID(id uuid.UUID) MeaningID { return MeaningID(id) }
func (id MeaningID) UUID() uuid.UUID           { return uuid.UUID(id) }
func (id MeaningID) String() string             { return uuid.UUID(id).String() }
func (id MeaningID) IsZero() bool               { return uuid.UUID(id) == uuid.Nil }
