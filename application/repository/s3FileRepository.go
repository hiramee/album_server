package repository

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type S3FileRepository struct {
	cli        *s3.S3
	bucketName string
	ctx        context.Context
}

func NewS3FileRepository(region string, bucketName string, ctx context.Context) *S3FileRepository {
	sess := session.Must(session.NewSession())
	cli := s3.New(sess, aws.NewConfig().WithRegion(region))
	repo := new(S3FileRepository)
	repo.cli = cli
	repo.bucketName = bucketName
	repo.ctx = ctx
	xray.AWS(cli.Client)
	return repo
}

func (repo *S3FileRepository) Upload(objectKey string, data []byte) error {
	req := new(s3.PutObjectInput)
	_, err := repo.cli.PutObjectWithContext(repo.ctx, req.SetBucket(repo.bucketName).SetKey(objectKey).SetBody(bytes.NewReader(data)))
	if err != nil {
		return err
	}
	return nil
}

func (repo *S3FileRepository) Download(objectKey string) ([]byte, error) {
	req := new(s3.GetObjectInput)
	res, err := repo.cli.GetObjectWithContext(repo.ctx, req.SetBucket(repo.bucketName).SetKey(objectKey))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	return buf.Bytes(), nil
}

func (repo *S3FileRepository) Delete(objectKey string) error {
	req := new(s3.DeleteObjectInput)
	_, err := repo.cli.DeleteObjectWithContext(repo.ctx, req.SetBucket(repo.bucketName).SetKey(objectKey))
	if err != nil {
		return err
	}
	return nil
}
