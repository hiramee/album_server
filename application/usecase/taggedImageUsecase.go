package usecase

import (
	"album-server/application/repository"
	"album-server/consts"
	"album-server/domain"
	"album-server/openapi"
	"album-server/util"
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
)

type TaggedImageUsecase struct {
	repo   *repository.TaggedImageRepository
	s3Repo *repository.S3FileRepository
}

// constructor
func NewTaggedImageUsecase() *TaggedImageUsecase {
	usecase := new(TaggedImageUsecase)
	usecase.repo = repository.NewTaggedImageRepository()
	usecase.s3Repo = repository.NewS3FileRepository(consts.Region, consts.AlbumFileBucket)
	return usecase
}

func (usecase *TaggedImageUsecase) SaveTaggedImage(userName string, tagNames []string, ext string, data []byte) error {
	uuid, err := usecase.uuid()
	if err != nil {
		return err
	}
	objectKey := userName + "/" + uuid + "." + ext
	uuid, err = usecase.uuid()
	if err != nil {
		return err
	}
	thumbNailObjectKey := userName + "/" + uuid + "." + "png"
	slice := domain.NewTaggedImageSlice(userName, tagNames, uuid, objectKey, thumbNailObjectKey)
	if err := usecase.repo.BatchUpdate(slice); err != nil {
		return err
	}
	if err := usecase.s3Repo.Upload(objectKey, data); err != nil {
		return err
	}
	var thumbNail *[]byte
	if thumbNail, err = createThumbNailData(&data); err != nil {
		return err
	}
	if err := usecase.s3Repo.Upload(thumbNailObjectKey, *thumbNail); err != nil {
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

func (usecase *TaggedImageUsecase) ListTaggedImageByTagNames(userName string, tagNameSlice []string) (*openapi.GetPicturesResponse, error) {
	taggedImages, err := usecase.repo.BatchGet(userName, tagNameSlice)
	if err != nil {
		return nil, err
	}
	picturesResponseItem := make([]openapi.PicturesResponseItem, len(taggedImages))
	for i, e := range taggedImages {
		fileName := util.TrimPrefixFromString(&e.ObjectKey, &userName) // userName/fileName => fileName
		if err != nil {
			return nil, err
		}
		tags := make([]string, len(e.Tags))
		for itag, etag := range e.Tags {
			tags[itag] = etag
		}
		id := e.ID
		picturesResponseItem[i].Id = &id
		picturesResponseItem[i].FileName = fileName
		picturesResponseItem[i].Tags = &tags
	}
	response := new(openapi.GetPicturesResponse)
	response.Pictures = &picturesResponseItem
	return response, nil
}

func (usecase *TaggedImageUsecase) UpdateTaggedImage(userName string, id string, tagNames []string) error {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return err
	}
	if len(current) == 0 {
		return errors.New("file already deleted")
	}
	var objectKey string
	var thumbNailObjectKey string
	currentTagNames := make([]string, len(current))
	for i, e := range current {
		objectKey = e.ObjectKey
		thumbNailObjectKey = e.ThumbNailObjectKey
		currentTagNames[i] = *util.TrimPrefixFromString(&e.UserTagName, &userName) // userNameTagName => TagName
	}
	deleteTarget := util.GetTwoSliceDiff(currentTagNames, tagNames)
	if err := usecase.repo.BatchDelete(id, userName, deleteTarget); err != nil {
		return err
	}
	if err := usecase.repo.BatchUpdate(domain.NewTaggedImageSlice(userName, tagNames, id, objectKey, thumbNailObjectKey)); err != nil {
		return err
	}
	return nil
}

func (usecase *TaggedImageUsecase) DeleteTaggedImage(userName string, id string) error {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return err
	}
	if len(current) == 0 {
		return errors.New("file is already deleted")
	}
	tags := current[0].Tags
	if err := usecase.repo.BatchDelete(id, userName, tags); err != nil {
		return err
	}
	objectKey := current[0].ObjectKey
	if err := usecase.s3Repo.Delete(objectKey); err != nil {
		return err
	}
	thumbNailObjectKey := current[0].ThumbNailObjectKey
	if err := usecase.s3Repo.Delete(thumbNailObjectKey); err != nil {
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
	uniqueUsedTags := util.GetUniqueSlice(usedTags)

	return uniqueUsedTags, nil
}

func (usecase *TaggedImageUsecase) GetImageById(id string, userName string) (*openapi.GetPictureResponse, error) {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return nil, err
	}
	if len(current) == 0 {
		return nil, errors.New("file already deleted")
	}
	objectKey := current[0].ObjectKey

	data, err := usecase.s3Repo.Download(objectKey)
	if err != nil {
		return nil, err
	}
	dataStr := base64.StdEncoding.EncodeToString(data)
	response := new(openapi.GetPictureResponse)
	response.Picture = &dataStr
	fileName := util.TrimPrefixFromString(&current[0].ObjectKey, &userName)
	response.FileName = fileName

	return response, nil
}

func (usecase *TaggedImageUsecase) GetThumbNailImageById(id string, userName string) (*openapi.GetPictureResponse, error) {
	current, err := usecase.repo.ListAllById(id, userName)
	if err != nil {
		return nil, err
	}
	if len(current) == 0 {
		return nil, errors.New("file already deleted")
	}
	objectKey := current[0].ThumbNailObjectKey

	data, err := usecase.s3Repo.Download(objectKey)
	if err != nil {
		return nil, err
	}
	dataStr := base64.StdEncoding.EncodeToString(data)
	response := new(openapi.GetPictureResponse)
	response.Picture = &dataStr
	fileName := util.TrimPrefixFromString(&objectKey, &userName)
	response.FileName = fileName
	return response, nil
}

// Create png image whose width is 256
func createThumbNailData(src *[]byte) (*[]byte, error) {
	img, _, err := image.Decode(bytes.NewBuffer(*src))

	if err != nil {
		println(err)
		return nil, errors.New("unsupported image type")
	}
	width := 256
	aspectRatio := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	height := int(float64(width) * aspectRatio)
	imgDst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buff []byte
	output := bytes.NewBuffer(buff)
	if err := png.Encode(output, imgDst); err != nil {
		println(err)
		return nil, errors.New("saving image failed")
	}
	thumbNail := output.Bytes()
	return &thumbNail, nil
}
