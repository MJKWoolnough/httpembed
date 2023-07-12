package httpembed

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestDecompressFS(t *testing.T) {
	c := new(compressor)

	const (
		valA = "Hello, World!"
		valB = "12345ABC"
	)
	dir := t.TempDir()

	bufA, err := c.Compress(valA)
	if err != nil {
		t.Errorf("test 1: unexpected error: %s", err)
		return
	}

	if err = writeFile(filepath.Join(dir, "a.txt.gz"), bufA); err != nil {
		t.Errorf("test 3: unexpected error: %s", err)
		return
	}

	bufB, err := c.Compress(valB)
	if err != nil {
		t.Errorf("test 2: unexpected error: %s", err)
		return
	}

	if err = os.Mkdir(filepath.Join(dir, "dir"), os.ModePerm); err != nil {
		t.Errorf("test 4: unexpected error: %s", err)
		return
	}

	if err = writeFile(filepath.Join(dir, "dir", "b.zip.gz"), bufB); err != nil {
		t.Errorf("test 5: unexpected error: %s", err)
		return
	}

	dd, err := DecompressFS(os.DirFS(dir))
	if err != nil {
		t.Errorf("test 6: unexpected error: %s", err)
		return
	}

	f, err := dd.Open("a.txt")
	if err != nil {
		t.Errorf("test 7: unexpected error: %s", err)
		return
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("test 8: unexpected error: %s", err)
		return
	}

	if string(buf) != valA {
		t.Errorf("test 9: expected to read %s, read %s", valA, buf)
		return
	}

	f, err = dd.Open("dir/b.zip")
	if err != nil {
		t.Errorf("test 10: unexpected error: %s", err)
		return
	}

	buf, err = io.ReadAll(f)
	if err != nil {
		t.Errorf("test 11: unexpected error: %s", err)
		return
	}

	if string(buf) != valB {
		t.Errorf("test 12: expected to read %s, read %s", valB, buf)
		return
	}
}

func writeFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
