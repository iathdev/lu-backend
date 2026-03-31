package domain

import "time"

// Folder is the aggregate root for user-created vocabulary decks.
// Scoped per language (1 folder = 1 language).
type Folder struct {
	ID          FolderID
	UserID      UserID
	LanguageID  LanguageID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewFolder(userID UserID, languageID LanguageID, name, description string) (*Folder, error) {
	if name == "" {
		return nil, ErrFolderNameRequired
	}

	return &Folder{
		ID:          NewFolderID(),
		UserID:      userID,
		LanguageID:  languageID,
		Name:        name,
		Description: description,
	}, nil
}

func (folder *Folder) Update(name, description string) error {
	if name == "" {
		return ErrFolderNameRequired
	}
	folder.Name = name
	folder.Description = description
	return nil
}
