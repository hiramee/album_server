package usecase

import (
	"album-server/application/repository"
	"album-server/domain"
)

type TagUsecase struct {
	repo *repository.TagRepository
}

// constructor
func NewTagUsecase() *TagUsecase {
	usecase := new(TagUsecase)
	usecase.repo = repository.NewTagRepository()
	return usecase
}

func (usecase *TagUsecase) ListAll(userName string) ([]domain.Tag, error) {
	return usecase.repo.ListAll(userName)
}

func (usecase *TagUsecase) DeleteByIdAndCategory(id int64, category string) error {
	return usecase.repo.DeleteByIdAndCategory(id, category)
}

func (usecase *TagUsecase) Update(domain *domain.Tag) error {
	return usecase.repo.Update(domain)
}

func (usecase *TagUsecase) CreateIfAbsent(userName string, tags []string) error {
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
