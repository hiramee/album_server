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

func (repo *TagRepository) BatchUpdate(domains []domain.Tag) error {
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
		if _, err := repo.table.Batch("UserName", "TagName").Write().Put(items...).Run(); err != nil {
			return err
		}
	}
	return nil
}
