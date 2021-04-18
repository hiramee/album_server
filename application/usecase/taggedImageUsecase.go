package usecase

import (
	"album-server/application/repository"
	"album-server/domain"
	"album-server/openapi"
	"encoding/base64"
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
)

type TaggedImageUsecase struct {
	repo   *repository.TaggedImageRepository
	s3Repo *repository.S3FileRepository
}

// constructor
func NewTaggedImageUsecase() *TaggedImageUsecase {
	usecase := new(TaggedImageUsecase)
	usecase.repo = repository.NewTaggedImageRepository()
	usecase.s3Repo = repository.NewS3FileRepository()
	return usecase
}

func (usecase *TaggedImageUsecase) Create(userName string, tagNames []string, ext string, data []byte) error {
	uuid, err := usecase.uuid()
	if err != nil {
		return err
	}
	objectKey := userName + "/" + uuid + "." + ext
	slice := domain.NewTaggedImageSlice(userName, tagNames, uuid, objectKey)
	fmt.Print(objectKey)
	if err := usecase.repo.BatchUpdate(slice); err != nil {
		return err
	}
	if err := usecase.s3Repo.Upload(objectKey, data); err != nil {
		return err
	}
	return nil
}

func (usecase *TaggedImageUsecase) uuid() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), err
}

func (usecase *TaggedImageUsecase) ListByTagNames(userName string, tagNameSlice []string) (*openapi.GetPicturesResponse, error) {
	response := new(openapi.GetPicturesResponse)
	taggedImages, err := usecase.repo.BatchGet(userName, tagNameSlice)
	if err != nil {
		return nil, err
	}
	picturesResponseItem := make([]openapi.PicturesResponseItem, len(taggedImages))
	for i, e := range taggedImages {
		prefixCount := utf8.RuneCountInString(userName) + 1 // userName/fileName => fileName
		fileName := string([]rune(e.ObjectKey)[prefixCount:])
		data, err := usecase.s3Repo.Download(e.ObjectKey)
		if err != nil {
			return nil, err
		}
		tags := make([]string, len(e.Tags))
		for itag, etag := range e.Tags {
			tags[itag] = etag
		}
		id := e.ID
		dataStr := base64.StdEncoding.EncodeToString(data)
		picturesResponseItem[i].Id = &id
		picturesResponseItem[i].FileName = &fileName
		picturesResponseItem[i].Tags = &tags
		picturesResponseItem[i].Picture = &dataStr
	}
	response.Pictures = &picturesResponseItem

	return response, nil
}
func (usecase *TaggedImageUsecase) Update(userName string, id string, tagNames []string) error {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return err
	}
	var objectKey string
	currentTagNames := make([]string, len(current))
	for i, e := range current {
		objectKey = e.ObjectKey
		prefixCount := utf8.RuneCountInString(userName) + 1 // userName/fileName => fileName
		currentTagNames[i] = string([]rune(e.UserTagName)[prefixCount:])
	}
	deleteTarget := findDiff(currentTagNames, tagNames)
	if err := usecase.repo.BatchDelete(id, userName, deleteTarget); err != nil {
		return err
	}
	if err := usecase.repo.BatchUpdate(domain.NewTaggedImageSlice(userName, tagNames, id, objectKey)); err != nil {
		return err
	}
	return nil
}

func findDiff(a []string, b []string) []string {
	var amap map[string]bool
	for _, e := range a {
		if !amap[e] {
			amap[e] = true
		}
	}
	var diff []string
	for _, e := range b {
		if !amap[e] {
			diff = append(diff, e)
		}
	}
	return diff
}
