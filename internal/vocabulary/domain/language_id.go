package domain

import "github.com/google/uuid"

// LanguageID uniquely identifies a Language entity.
type LanguageID uuid.UUID

func NewLanguageID() LanguageID {
	return LanguageID(uuid.Must(uuid.NewV7()))
}

func ParseLanguageID(raw string) (LanguageID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return LanguageID{}, err
	}
	return LanguageID(parsed), nil
}

func LanguageIDFromUUID(id uuid.UUID) LanguageID { return LanguageID(id) }
func (id LanguageID) UUID() uuid.UUID            { return uuid.UUID(id) }
func (id LanguageID) String() string              { return uuid.UUID(id).String() }
func (id LanguageID) IsZero() bool                { return uuid.UUID(id) == uuid.Nil }
