package domain

import "time"

// Topic is a system-defined classification entity for vocabularies, scoped to a category.
type Topic struct {
	ID         TopicID
	CategoryID CategoryID
	Slug       string
	Names      map[string]string
	Offset     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewTopic(categoryID CategoryID, slug string, names map[string]string, offset int) (*Topic, error) {
	if slug == "" {
		return nil, ErrTopicSlugRequired
	}

	return &Topic{
		ID:         NewTopicID(),
		CategoryID: categoryID,
		Slug:       slug,
		Names:      names,
		Offset:     offset,
	}, nil
}
