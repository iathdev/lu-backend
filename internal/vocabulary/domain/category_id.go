package domain

import "github.com/google/uuid"

// CategoryID uniquely identifies a Category entity.
type CategoryID uuid.UUID

func NewCategoryID() CategoryID {
	return CategoryID(uuid.Must(uuid.NewV7()))
}

func ParseCategoryID(raw string) (CategoryID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return CategoryID{}, err
	}
	return CategoryID(parsed), nil
}

func CategoryIDFromUUID(id uuid.UUID) CategoryID { return CategoryID(id) }
func (id CategoryID) UUID() uuid.UUID            { return uuid.UUID(id) }
func (id CategoryID) String() string              { return uuid.UUID(id).String() }
func (id CategoryID) IsZero() bool                { return uuid.UUID(id) == uuid.Nil }
