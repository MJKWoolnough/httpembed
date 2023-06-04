package httpembed

import (
	"errors"
	"io"
	"testing"
	"time"
)

func TestFiles(t *testing.T) {
	aTime := time.Now()
	bTime := time.Now().Add(-5 * time.Minute)

	d := decompressedFS{
		files: map[string]file{
			"a.txt": {
				name:    "a.txt",
				data:    []byte("Hello, World!"),
				modTime: aTime,
			},
			"dir/b.zip": {
				name:    "b.zip",
				data:    []byte("12345ABC"),
				modTime: bTime,
			},
		},
	}

	a, err := d.Open("a.txt")
	if err != nil {
		t.Errorf("test 1: unexpected error: %s", err)
		return
	}

	s, err := a.Stat()
	if err != nil {
		t.Errorf("test 2: unexpected error: %s", err)
		return
	}

	if modTime := s.ModTime(); !modTime.Equal(aTime) {
		t.Errorf("test 3: expecting modtime %s, got %s", aTime, modTime)
		return
	}

	if size := s.Size(); size != 13 {
		t.Errorf("test 4: expecting size 14, got %d", size)
		return
	}

	buf := make([]byte, 5)

	n, err := a.Read(buf)

	if err != nil {
		t.Errorf("test 5: unexpected error: %s", err)
		return
	} else if n != 5 {
		t.Errorf("test 6: expecting to read 5 bytes, read %d", n)
		return
	} else if string(buf) != "Hello" {
		t.Errorf("test 7: expecting to read \"Hello\", read %q", buf)
		return
	}

	n, err = a.Read(buf)

	if err != nil {
		t.Errorf("test 8: unexpected error: %s", err)
		return
	} else if n != 5 {
		t.Errorf("test 9: expecting to read 5 bytes, read %d", n)
		return
	} else if string(buf) != ", Wor" {
		t.Errorf("test 10: expecting to read \", Wor\", read %q", buf)
		return
	}

	n, err = a.Read(buf)

	if err != nil {
		t.Errorf("test 11: unexpected error: %s", err)
		return
	} else if n != 3 {
		t.Errorf("test 12: expecting to read 3 bytes, read %d", n)
		return
	} else if string(buf[:n]) != "ld!" {
		t.Errorf("test 13: expecting to read \"ld!\", read %q", buf[:n])
		return
	}

	n, err = a.Read(buf)

	if !errors.Is(err, io.EOF) {
		t.Errorf("test 14: expected error EOF, got: %s", err)
		return
	} else if n != 0 {
		t.Errorf("test 15: expecting to read 0 bytes, read %d", n)
		return
	}

	p, err := a.(io.Seeker).Seek(-6, io.SeekEnd)

	if err != nil {
		t.Errorf("test 16: unexpected error: %s", err)
		return
	} else if p == 8 {
		t.Errorf("test 17: expecting to seek to 8, got %d", p)
		return
	}

	n, err = a.Read(buf)

	if err != nil {
		t.Errorf("test 18: unexpected error: %s", err)
		return
	} else if n != 5 {
		t.Errorf("test 19: expecting to read 5 bytes, read %d", n)
		return
	} else if string(buf) != "World" {
		t.Errorf("test 20: expecting to read \"World\", read %q", buf)
		return
	}

	a, err = d.Open("a.txt")
	if err != nil {
		t.Errorf("test 21: unexpected error: %s", err)
		return
	}

	n, err = a.Read(buf)

	if err != nil {
		t.Errorf("test 22: unexpected error: %s", err)
		return
	} else if n != 5 {
		t.Errorf("test 23: expecting to read 5 bytes, read %d", n)
		return
	} else if string(buf) != "Hello" {
		t.Errorf("test 24: expecting to read \"Hello\", read %q", buf)
		return
	}

	b, err := d.Open("dir/b.zip")
	if err != nil {
		t.Errorf("test 26: unexpected error: %s", err)
		return
	}

	s, err = b.Stat()
	if err != nil {
		t.Errorf("test 27: unexpected error: %s", err)
		return
	}

	if modTime := s.ModTime(); !modTime.Equal(bTime) {
		t.Errorf("test 28: expecting modtime %s, got %s", bTime, modTime)
		return
	}

	if size := s.Size(); size != 8 {
		t.Errorf("test 29: expecting size 8, got %d", size)
		return
	}

	_, err = d.Open("unknownfile")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("test 30: expecting ErrNotFound, got %s", err)
	}
}
