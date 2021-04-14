package repository

import (
	"album-server/domain"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type TaggedImageRepository struct {
	table *dynamo.Table
}

// constructor
func NewTaggedImageRepository() *TaggedImageRepository {
	sess, err := session.NewSession()
	if err != nil {
		// retry once
		fmt.Print(err)
		sess, err = session.NewSession()
		fmt.Print(err)
	}
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
	table := db.Table("TaggedImage")
	repo := new(TaggedImageRepository)
	repo.table = &table
	return repo
}

func (repo *TaggedImageRepository) ListAll(category string) ([]domain.TaggedImage, error) {
	var results []domain.TaggedImage
	if err := repo.table.Scan().Filter("'Category' = ?", category).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *TaggedImageRepository) DeleteByIdAndCategory(id int64, category string) error {
	if err := repo.table.Delete("ID", id).Range("Category", category).Run(); err != nil {
		return err
	}
	return nil
}

func (repo *TaggedImageRepository) Update(domain *domain.TaggedImage) error {
	if err := repo.table.Put(domain).Run(); err != nil {
		return err
	}
	return nil
}

func (repo *TaggedImageRepository) BatchUpdate(domains []domain.TaggedImage) error {
	sliceSize := len(domains)
	for i := 0; i < sliceSize; i += 25 {
		end := i + 25
		if sliceSize < end {
			end = sliceSize
		}
		current := domains[i:end]
		items := make([]interface{}, len(current))
		for i, e := range current {
			items[i] = e
		}
		if _, err := repo.table.Batch("UserTagName", "ID").Write().Put(items...).Run(); err != nil {
			return err
		}
	}
	return nil
}
