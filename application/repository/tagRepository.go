package repository

import (
	"album-server/domain"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/guregu/dynamo"
)

type TagRepository struct {
	table *dynamo.Table
	ctx   context.Context
}

// constructor
func NewTagRepository(ctx context.Context) *TagRepository {
	sess, err := session.NewSession()
	if err != nil {
		// retry once
		fmt.Print(err)
		sess, err = session.NewSession()
		fmt.Print(err)
	}
	dynamodb := dynamodb.New(sess)
	xray.AWS(dynamodb.Client)
	db := dynamo.NewFromIface(dynamodb)

	table := db.Table("Tag")
	repo := new(TagRepository)
	repo.table = &table
	repo.ctx = ctx
	return repo
}

func (repo *TagRepository) ListAll(userName string) ([]domain.Tag, error) {
	var results []domain.Tag
	if err := repo.table.Get("UserName", userName).AllWithContext(repo.ctx, &results); err != nil {
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
		if _, err := repo.table.Batch("UserName", "TagName").Write().Put(items...).RunWithContext(repo.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (repo *TagRepository) BatchDelete(userName string, tags []string) error {
	sliceSize := len(tags)
	for i := 0; i < sliceSize; i += 25 {
		end := i + 25
		if sliceSize < end {
			end = sliceSize
		}
		current := tags[i:end]
		items := make([]dynamo.Keyed, len(current))
		for i, e := range current {
			items[i] = dynamo.Keys{userName, e}
		}
		if _, err := repo.table.Batch("UserName", "TagName").Write().Delete(items...).RunWithContext(repo.ctx); err != nil {
			return err
		}
	}
	return nil
}
