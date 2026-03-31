package mapper

import (
	vdto "learning-go/internal/vocabulary/application/dto"
	"learning-go/internal/vocabulary/domain"
)

// ToFolderResponse maps domain.Folder to FolderResponse, including vocab count.
func ToFolderResponse(folder *domain.Folder, vocabCount int) *vdto.FolderResponse {
	return &vdto.FolderResponse{
		ID:              folder.ID.String(),
		UserID:          folder.UserID.String(),
		LanguageID:      folder.LanguageID.String(),
		Name:            folder.Name,
		Description:     folder.Description,
		VocabularyCount: vocabCount,
		CreatedAt:       folder.CreatedAt,
	}
}
