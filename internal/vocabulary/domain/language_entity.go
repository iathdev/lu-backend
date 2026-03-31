package domain

import "time"

// Language is a top-level entity representing a supported language.
type Language struct {
	ID         LanguageID
	Code       string
	NameEN     string
	NameNative string
	IsActive   bool
	Config     map[string]any
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
