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

type TaggedImageRepository struct {
	table *dynamo.Table
	ctx   context.Context
}

// constructor
func NewTaggedImageRepository(ctx context.Context) *TaggedImageRepository {
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
	table := db.Table("TaggedImage")
	repo := new(TaggedImageRepository)
	repo.table = &table
	repo.ctx = ctx
	return repo
}

func (repo *TaggedImageRepository) ListAllById(id string, userName string) ([]domain.TaggedImage, error) {
	var results []domain.TaggedImage
	if err := repo.table.Get("ID", id).Range("UserTagName", dynamo.Greater, userName).AllWithContext(repo.ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *TaggedImageRepository) DeleteByIdAndCategory(id int64, category string) error {
	if err := repo.table.Delete("ID", id).Range("Category", category).RunWithContext(repo.ctx); err != nil {
		return err
	}
	return nil
}

func (repo *TaggedImageRepository) Update(domain *domain.TaggedImage) error {
	if err := repo.table.Put(domain).RunWithContext(repo.ctx); err != nil {
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
		if _, err := repo.table.Batch("ID", "UserTagName").Write().Put(items...).RunWithContext(repo.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (repo *TaggedImageRepository) BatchGet(userName string, tagNames []string) ([]domain.TaggedImage, error) {
	// TODO: ??????????????????????????????
	idKeyImageMap := make(map[string]domain.TaggedImage)
	for _, e := range tagNames {
		var results []domain.TaggedImage
		if err := repo.table.Get("UserTagName", userName+e).Index("GSI-UserTagName").AllWithContext(repo.ctx, &results); err != nil {
			return nil, err
		}
		for _, e := range results {
			if hasAllTags(tagNames, e.Tags) {
				if _, ok := idKeyImageMap[e.ID]; !ok {
					idKeyImageMap[e.ID] = e
				}
			}
		}
	}
	var response []domain.TaggedImage
	for _, e := range idKeyImageMap {
		response = append(response, e)
	}
	return response, nil
}

func hasAllTags(tags []string, tested []string) bool {
	tagKeyMap := make(map[string]bool)
	for _, e := range tested {
		if _, has := tagKeyMap[e]; !has {
			tagKeyMap[e] = true
		}
	}
	for _, e := range tags {
		if _, has := tagKeyMap[e]; !has {
			return false
		}
	}
	return true
}

func (repo *TaggedImageRepository) BatchDelete(id string, userName string, tags []string) error {
	sliceSize := len(tags)
	for i := 0; i < sliceSize; i += 25 {
		end := i + 25
		if sliceSize < end {
			end = sliceSize
		}
		current := tags[i:end]
		items := make([]dynamo.Keyed, len(current))
		for i, e := range current {
			items[i] = dynamo.Keys{id, userName + e}
		}
		if _, err := repo.table.Batch("ID", "UserTagName").Write().Delete(items...).RunWithContext(repo.ctx); err != nil {
			return err
		}
	}
	return nil
}
