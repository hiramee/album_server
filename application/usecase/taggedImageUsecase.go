package usecase

import (
	"album-server/application/repository"
	"album-server/domain"
	"album-server/openapi"
	"encoding/base64"
	"errors"
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
		if err != nil {
			return nil, err
		}
		tags := make([]string, len(e.Tags))
		for itag, etag := range e.Tags {
			tags[itag] = etag
		}
		id := e.ID
		picturesResponseItem[i].Id = &id
		picturesResponseItem[i].FileName = &fileName
		picturesResponseItem[i].Tags = &tags
	}
	response.Pictures = &picturesResponseItem

	return response, nil
}
func (usecase *TaggedImageUsecase) Update(userName string, id string, tagNames []string) error {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return err
	}
	if len(current) == 0 {
		return errors.New("File already deleted")
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
	amap := make(map[string]bool)
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

func (usecase *TaggedImageUsecase) Delete(userName string, id string) error {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return err
	}
	if len(current) == 0 {
		return errors.New("File is already deleted")
	}
	tags := current[0].Tags
	if err := usecase.repo.BatchDelete(id, userName, tags); err != nil {
		return err
	}
	objectKey := current[0].ObjectKey
	if err := usecase.s3Repo.Delete(objectKey); err != nil {
		return err
	}
	return nil
}

func (usecase *TaggedImageUsecase) ValidateDeleteTags(userName string, tagNameSlice []string) ([]string, error) {
	taggedImages, err := usecase.repo.BatchGet(userName, tagNameSlice)
	if err != nil {
		return nil, err
	}
	if len(taggedImages) == 0 {
		return nil, nil
	}
	var usedTags []string
	for _, e := range taggedImages {
		usedTags = append(usedTags, e.Tags...)
	}
	uniqueUsedTags := unique(usedTags)

	return uniqueUsedTags, nil
}

func unique(a []string) []string {
	amap := make(map[string]bool)
	var uniqueSlice []string
	for _, e := range a {
		if !amap[e] {
			uniqueSlice = append(uniqueSlice, e)
			amap[e] = true
		}
	}
	return uniqueSlice
}

func (usecase *TaggedImageUsecase) GetById(id string, userName string) (*openapi.GetPictureResponse, error) {
	response := new(openapi.GetPictureResponse)
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return nil, err
	}
	if len(current) == 0 {
		return nil, errors.New("File already deleted")
	}
	objectKey := current[0].ObjectKey

	data, err := usecase.s3Repo.Download(objectKey)
	if err != nil {
		return nil, err
	}
	dataStr := base64.StdEncoding.EncodeToString(data)
	response.Picture = &dataStr

	return response, nil
}
