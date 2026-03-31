package domain

import "github.com/google/uuid"

// ExampleID uniquely identifies a VocabularyExample entity.
type ExampleID uuid.UUID

func NewExampleID() ExampleID {
	return ExampleID(uuid.Must(uuid.NewV7()))
}

func ParseExampleID(raw string) (ExampleID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ExampleID{}, err
	}
	return ExampleID(parsed), nil
}

func ExampleIDFromUUID(id uuid.UUID) ExampleID { return ExampleID(id) }
func (id ExampleID) UUID() uuid.UUID           { return uuid.UUID(id) }
func (id ExampleID) String() string             { return uuid.UUID(id).String() }
func (id ExampleID) IsZero() bool               { return uuid.UUID(id) == uuid.Nil }
