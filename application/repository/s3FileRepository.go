package repository

import (
	"bytes"
	"io"

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

func (repo *S3FileRepository) Download(objectKey string) ([]byte, error) {
	req := new(s3.GetObjectInput)
	res, err := repo.cli.GetObject(req.SetBucket(repo.bucketName).SetKey(objectKey))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	return buf.Bytes(), nil
}