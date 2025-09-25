package httpembed

import (
	"compress/gzip"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"vimagination.zapto.org/memfs"
)

// DecompressFS takes a FS with compressed (.gz) files and returns a new FS with
// those files decompressed and stored under the same name with the .gz suffix
// removed.
//
// The output of this is intended to be use with httpgzip.FileServer.
func DecompressFS(files fs.FS) (fs.FS, error) {
	g := new(gzip.Reader)
	mfs := memfs.New()

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

		info, err := f.Stat()
		if err != nil {
			return err
		}

		dir := filepath.Dir(path)

		if dir != "." && dir != "" && dir != "/" {
			if err = mfs.MkdirAll(dir, fs.ModePerm); err != nil {
				return err
			}
		}

		return writeFileToFS(mfs, strings.TrimSuffix(path, ".gz"), g, info)
	}); err != nil {
		return nil, err
	}

	return mfs.Seal(), nil
}

func writeFileToFS(mfs *memfs.FS, gpath string, g io.Reader, info fs.FileInfo) error {
	gf, err := mfs.Create(gpath)
	if err != nil {
		return err
	}

	_, err = io.Copy(gf, g)
	if err != nil {
		return err
	}

	if err = gf.Close(); err != nil {
		return err
	}

	mtime := info.ModTime()

	if err = mfs.Chtimes(gpath, mtime, mtime); err != nil {
		return err
	}

	if err = mfs.Chmod(gpath, info.Mode()); err != nil {
		return err
	}

	return nil
}
