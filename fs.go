package httpembed

import (
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"strings"
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
	if f.pos >= len(f.data) {
		return 0, io.EOF
	}
	n := copy(p, f.data[f.pos:])
	f.pos += n
	return n, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.pos = int(offset)
	case io.SeekCurrent:
		f.pos += int(offset)
	case io.SeekEnd:
		f.pos = len(f.data) + int(offset)
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

// DecompressFS takes a FS with compressed (.gz) files and returns a new FS with
// those files decompressed and store under the same name with the .gz suffix
// removed.
func DecompressFS(files fs.FS) (fs.FS, error) {
	g := new(gzip.Reader)
	h := make(map[string]file)

	if err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if !d.Type().IsRegular() || !strings.HasSuffix(name, ".gz") {
			return nil
		}

		f, err := files.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()

		if err = g.Reset(f); err != nil {
			return err
		}

		buf, err := io.ReadAll(g)
		if err != nil {
			return err
		}

		info, err := f.Stat()
		if err != nil {
			return err
		}

		h[strings.TrimSuffix(path, ".gz")] = file{
			name:    strings.TrimSuffix(name, ".gz"),
			data:    buf,
			modTime: info.ModTime(),
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &decompressedFS{files: h}, nil
}
