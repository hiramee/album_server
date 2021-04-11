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
