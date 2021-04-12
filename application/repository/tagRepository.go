package repository

import (
	"album-server/domain"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type TagRepository struct {
	table *dynamo.Table
}

// constructor
func NewTagRepository() *TagRepository {
	sess, err := session.NewSession()
	if err != nil {
		// retry once
		fmt.Print(err)
		sess, err = session.NewSession()
		fmt.Print(err)
	}
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
	table := db.Table("Tag")
	repo := new(TagRepository)
	repo.table = &table
	return repo
}

func (repo *TagRepository) ListAll(userName string) ([]domain.Tag, error) {
	var results []domain.Tag
	if err := repo.table.Get("UserName", userName).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *TagRepository) DeleteByIdAndCategory(id int64, category string) error {
	if err := repo.table.Delete("ID", id).Range("Category", category).Run(); err != nil {
		return err
	}
	return nil
}

func (repo *TagRepository) Update(domain *domain.Tag) error {
	if err := repo.table.Put(domain).Run(); err != nil {
		return err
	}
	return nil
}

func (repo *TagRepository) BatchUpdate(domains []domain.Tag) error {
	sliceSize := len(domains)
	for i := 0; i < sliceSize; i += 25 {
		end := i + 25
		if sliceSize < end {
			end = sliceSize
		}
		if _, err := repo.table.Batch("UserName", "TagName").Write().Put(domains[i:end]).Run(); err != nil {
			return err
		}
	}
	return nil
}
