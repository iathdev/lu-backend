package domain

import "github.com/google/uuid"

// ProficiencyLevelID uniquely identifies a ProficiencyLevel entity.
type ProficiencyLevelID uuid.UUID

func NewProficiencyLevelID() ProficiencyLevelID {
	return ProficiencyLevelID(uuid.Must(uuid.NewV7()))
}

func ParseProficiencyLevelID(raw string) (ProficiencyLevelID, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ProficiencyLevelID{}, err
	}
	return ProficiencyLevelID(parsed), nil
}

func ProficiencyLevelIDFromUUID(id uuid.UUID) ProficiencyLevelID { return ProficiencyLevelID(id) }
func (id ProficiencyLevelID) UUID() uuid.UUID                    { return uuid.UUID(id) }
func (id ProficiencyLevelID) String() string                      { return uuid.UUID(id).String() }
func (id ProficiencyLevelID) IsZero() bool                        { return uuid.UUID(id) == uuid.Nil }
