package s3fs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
S3Client retroactively builds an interface for the S3 Client.
*/
type S3Client interface {
	HeadObject(ctx context.Context, input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

/*
S3ClientMock is a mock of the S3 Client.
*/
type S3ClientMock struct{}

func (m *S3ClientMock) HeadObject(ctx context.Context, input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	// Implement mock logic here
	return &s3.HeadObjectOutput{}, nil
}

func (m *S3ClientMock) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	// Implement mock logic here
	return &s3.ListObjectsV2Output{}, nil
}

func (m *S3ClientMock) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	// Implement mock logic here
	return &s3.PutObjectOutput{}, nil
}

func (m *S3ClientMock) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	// Implement mock logic here
	return &s3.GetObjectOutput{}, nil
}

func (m *S3ClientMock) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	// Implement mock logic here
	return &s3.DeleteObjectOutput{}, nil
}
