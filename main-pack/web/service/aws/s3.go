package aws

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/glog"
	"io"
)

type BucketClient interface {
	Put(bucket, fileName string, content []byte, ctx context.Context) (string, error)
	Delete(bucket, fileName string, ctx context.Context) error
	Get(bucket, fileName string, ctx context.Context) ([]byte, error)
	Exist(bucket string, fileName string) (bool, error)
}

type S3 struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func (s S3) Exist(bucket string, fileName string) (bool, error) {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func NewS3(config Config) BucketClient {
	session, err := newAWSSession(config)
	if err != nil {
		return nil
	}

	s3manager.NewUploader(session)
	return S3{
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

func (s S3) Put(bucket, fileName string, content []byte, ctx context.Context) (string, error) {
	res, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   bytes.NewReader(content),
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		glog.Errorf("upload: %s", err)
		return "", err
	}

	return res.Location, nil
}

func (s S3) Delete(bucket, fileName string, ctx context.Context) error {
	if _, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		glog.Errorf("delete: %s", err)
		return err
	}

	if err := s.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		glog.Errorf("wait: %s", err)
		return err
	}

	return nil
}

func (s S3) Get(bucket, fileName string, ctx context.Context) ([]byte, error) {
	res, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
