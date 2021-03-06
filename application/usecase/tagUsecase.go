package usecase

import (
	"album-server/application/repository"
	"album-server/domain"
	"context"
)

type TagUsecase struct {
	repo *repository.TagRepository
}

// constructor
func NewTagUsecase(ctx context.Context) *TagUsecase {
	usecase := new(TagUsecase)
	usecase.repo = repository.NewTagRepository(ctx)
	return usecase
}

func (usecase *TagUsecase) ListAll(userName string) ([]domain.Tag, error) {
	return usecase.repo.ListAll(userName)
}

func (usecase *TagUsecase) SaveTagIfAbsent(userName string, tags []string) error {
	oldTags, err := usecase.repo.ListAll(userName)
	if err != nil {
		return err
	}
	newTagStrs := findNewTags(oldTags, tags)
	var newTags []domain.Tag
	for _, el := range newTagStrs {
		newTag := new(domain.Tag)
		newTag.UserName = userName
		newTag.TagName = el
		newTags = append(newTags, *newTag)
	}
	return usecase.repo.BatchUpdate(newTags)
}

func findNewTags(oldTags []domain.Tag, tags []string) []string {
	m := make(map[string]struct{})
	for _, el := range oldTags {
		if _, ok := m[el.TagName]; !ok {
			m[el.TagName] = struct{}{}
		}
	}
	var newTags []string
	for _, el := range tags {
		if _, ok := m[el]; !ok {
			newTags = append(newTags, el)
		}
	}
	return newTags
}

func (usecase *TagUsecase) DeleteTag(userName string, tags []string) error {
	return usecase.repo.BatchDelete(userName, tags)
}
