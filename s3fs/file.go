package s3fs

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/spf13/afero"
)

// The FileInfo type implements os.FileInfo interface, and represents a file in S3
type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fi *FileInfo) Name() string       { return fi.name }
func (fi *FileInfo) Size() int64        { return fi.size }
func (fi *FileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *FileInfo) ModTime() time.Time { return fi.modTime }
func (fi *FileInfo) IsDir() bool        { return fi.isDir }
func (fi *FileInfo) Sys() interface{}   { return nil }

type S3File struct {
	client S3Client
	bucket string
	key    string
	offset int64
	buffer *bytes.Buffer
}

// Ensure S3File implements afero.File
var _ afero.File = (*S3File)(nil)

func NewS3File(client S3Client, bucket, key string) *S3File {
	return &S3File{
		client: client,
		bucket: bucket,
		key:    key,
		buffer: bytes.NewBuffer(nil),
	}
}

func (f *S3File) Close() error {
	// Implement the necessary cleanup; for read-only files, this might just be a no-op
	return nil
}

func (f *S3File) Read(p []byte) (int, error) {
	if f.buffer == nil {
		return 0, io.EOF
	}
	return f.buffer.Read(p)
}

func (f *S3File) ReadAt(p []byte, off int64) (int, error) {
	// ReadAt is not simple to implement directly using S3's API as it's a full object store,
	// Consider fetching full content or using range requests
	if off != f.offset {
		return 0, io.EOF
	}
	return f.Read(p)
}

func (f *S3File) Seek(offset int64, whence int) (int64, error) {
	size := int64(f.buffer.Len())
	switch whence {
	case io.SeekStart:
		f.offset = offset
	case io.SeekCurrent:
		f.offset += offset
	case io.SeekEnd:
		f.offset = size + offset
	default:
		return 0, io.ErrUnexpectedEOF
	}
	if f.offset > size {
		return 0, io.EOF
	}
	return f.offset, nil
}

func (f *S3File) Write(p []byte) (int, error) {
	// To keep things simple, assume writing is not supported
	return 0, io.ErrUnexpectedEOF
}

func (f *S3File) WriteAt(p []byte, off int64) (int, error) {
	// As above, assume writing is not supported
	return 0, io.ErrUnexpectedEOF
}

func (f *S3File) Name() string {
	return f.key
}

func (f *S3File) Readdir(count int) ([]os.FileInfo, error) {
	// Not applicable for S3 as it's not a true file system
	return nil, io.ErrUnexpectedEOF
}

func (f *S3File) Readdirnames(n int) ([]string, error) {
	// As above, not applicable
	return nil, io.ErrUnexpectedEOF
}

func (f *S3File) Stat() (os.FileInfo, error) {
	return &FileInfo{
		name:    f.key,
		size:    int64(f.buffer.Len()),
		mode:    os.FileMode(0644),
		modTime: time.Now(),
		isDir:   false,
	}, nil
}

func (f *S3File) Sync() error {
	// For S3, this might be a no-op
	return nil
}

func (f *S3File) Truncate(size int64) error {
	// Not applicable for S3 objects
	return io.ErrUnexpectedEOF
}

func (f *S3File) WriteString(s string) (ret int, err error) {
	// Assume not supported
	return 0, io.ErrUnexpectedEOF
}
