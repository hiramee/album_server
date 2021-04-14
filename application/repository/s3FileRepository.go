package repository

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3FileRepository struct {
	cli        *s3.S3
	bucketName string
}

func NewS3FileRepository() *S3FileRepository {
	sess := session.Must(session.NewSession())
	cli := s3.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))
	repo := new(S3FileRepository)
	repo.cli = cli
	repo.bucketName = "album-file-bucket" // TODO 切り出す
	return repo
}

func (repo *S3FileRepository) Upload(objectKey string, data []byte) error {
	req := new(s3.PutObjectInput)
	_, err := repo.cli.PutObject(req.SetBucket(repo.bucketName).SetKey(objectKey).SetBody(bytes.NewReader(data)))
	if err != nil {
		return err
	}
	return nil
}
