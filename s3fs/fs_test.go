package s3fs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testBucket = "test-bucket"
	testFile   = "test.txt"
)

func TestNew(t *testing.T) {
	client := &S3ClientMock{}
	fs := New(client, testBucket)

	assert.Equal(t, client, fs.(*Fs).client)
	assert.Equal(t, testBucket, fs.(*Fs).bucket)
}

func TestCreate(t *testing.T) {
	client := &S3ClientMock{}
	client.PutObjectMock = func(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		require.Equal(t, testBucket, *input.Bucket)
		require.Equal(t, testFile, *input.Key)
		return &s3.PutObjectOutput{}, nil
	}
	fs := New(client, testBucket)

	t.Run("success", func(t *testing.T) {
		f, err := fs.Create(testFile)
		require.NoError(t, err)
		assert.IsType(t, &S3File{}, f)
	})

	t.Run("error", func(t *testing.T) {
		client.PutObjectMock = func(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
			return nil, errors.New("test error")
		}

		_, err := fs.Create(testFile)
		require.Error(t, err)
	})
}

func TestOpen(t *testing.T) {
	client := &S3ClientMock{}
	client.GetObjectMock = func(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		content := []byte("hello world")
		return &s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader(content)),
		}, nil
	}
	fs := New(client, testBucket)

	t.Run("success", func(t *testing.T) {
		f, err := fs.Open(testFile)
		require.NoError(t, err)

		content := []byte("hello world")
		buf := make([]byte, len(content))
		n, err := f.Read(buf)
		require.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, buf)
	})

	t.Run("error", func(t *testing.T) {
		client.GetObjectMock = func(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			return nil, errors.New("test error")
		}

		_, err := fs.Open(testFile)
		require.Error(t, err)
	})
}
