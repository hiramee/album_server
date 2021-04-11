package usecase

import (
	"album-server/application/repository"
	"album-server/domain"
)

type TaggedImageUsecase struct {
	repo *repository.TaggedImageRepository
}

// constructor
func NewTaggedImageUsecase() *TaggedImageUsecase {
	usecase := new(TaggedImageUsecase)
	usecase.repo = repository.NewTaggedImageRepository()
	return usecase
}

func (usecase *TaggedImageUsecase) ListAll(category string) ([]domain.TaggedImage, error) {
	return usecase.repo.ListAll(category)
}

func (usecase *TaggedImageUsecase) DeleteByIdAndCategory(id int64, category string) error {
	return usecase.repo.DeleteByIdAndCategory(id, category)
}

func (usecase *TaggedImageUsecase) Update(domain *domain.TaggedImage) error {
	return usecase.repo.Update(domain)
}
