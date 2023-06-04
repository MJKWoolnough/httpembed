package httpembed

import (
	"errors"
	"io"
	"io/fs"
	"time"
)

var ErrNotFound = errors.New("file not found")

type decompressedFS struct {
	files map[string]file
}

func (d *decompressedFS) Open(name string) (fs.File, error) {
	f, ok := d.files[name]
	if !ok {
		return nil, ErrNotFound
	}
	return &f, nil
}

type file struct {
	name    string
	pos     int
	data    []byte
	modTime time.Time
}

func (f *file) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *file) Read(p []byte) (int, error) {
	n := copy(p, f.data[f.pos:])
	f.pos += n
	return n, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.pos = whence
	case io.SeekCurrent:
		f.pos += whence
	case io.SeekEnd:
		f.pos = len(f.data) + whence
	}
	return int64(f.pos), nil
}

func (f *file) Close() error {
	f.data = nil
	return nil
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Size() int64 {
	return int64(len(f.data))
}

func (f *file) Mode() fs.FileMode {
	return fs.ModePerm
}

func (f *file) ModTime() time.Time {
	return f.modTime
}

func (f *file) IsDir() bool {
	return false
}

func (f *file) Sys() any {
	return f
}
