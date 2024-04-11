package s3fs

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/afero"
)

type Fs struct {
	client S3Client
	bucket string // the name of the bucket
}

func New(client S3Client, bucket string) afero.Fs {
	return &Fs{client: client, bucket: bucket}
}

func (fs Fs) Name() string {
	return "S3"
}

func (fs *Fs) Create(name string) (afero.File, error) {
	// Create an empty object in S3 to simulate file creation
	_, err := fs.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &fs.bucket,
		Key:    &name,
		Body:   bytes.NewReader([]byte{}),
	})

	if err != nil {
		return nil, err
	}

	// Return an S3File instance prepared for read and write operations
	return NewS3File(fs.client, fs.bucket, name), nil
}

func (fs Fs) Mkdir(name string, perm os.FileMode) error {
	return nil
}

func (fs Fs) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

// Open a file from S3
func (fs *Fs) Open(name string) (afero.File, error) {
	resp, err := fs.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &fs.bucket,
		Key:    &name,
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	file := NewS3File(fs.client, fs.bucket, name)
	file.buffer = bytes.NewBuffer(content)
	return file, nil
}

func (fs Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return nil, nil
}

// Remove a file from S3
func (fs *Fs) Remove(name string) error {
	_, err := fs.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &fs.bucket,
		Key:    &name,
	})

	return err
}

func (fs Fs) RemoveAll(path string) error {
	return nil
}

func (fs Fs) Rename(oldname, newname string) error {
	return nil
}

// Stat returns a FileInfo describing the named file from S3
func (fs *Fs) Stat(name string) (os.FileInfo, error) {
	resp, err := fs.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &fs.bucket,
		Key:    &name,
	})

	if err != nil {
		return nil, err
	}

	return &FileInfo{
		name:    name,
		size:    *resp.ContentLength,
		mode:    os.FileMode(0644), // Default mode, as S3 does not store file mode
		modTime: *resp.LastModified,
		isDir:   false,
	}, nil
}

func (fs Fs) Chmod(name string, mode os.FileMode) error {
	return nil
}

func (fs Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}

func (fs Fs) Chown(name string, uid, gid int) error {
	return nil
}
