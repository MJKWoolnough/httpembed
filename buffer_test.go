package httpembed

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type compressor struct {
	buf bytes.Buffer
	gz  gzip.Writer
}

func (c *compressor) Compress(str string) ([]byte, error) {
	c.buf.Reset()
	c.gz.Reset(&c.buf)

	if _, err := io.WriteString(&c.gz, str); err != nil {
		return nil, err
	}

	if err := c.gz.Close(); err != nil {
		return nil, err
	}

	return c.buf.Bytes(), nil
}

func test(t *testing.T, fn func(test string, buf *bytes.Buffer) http.Handler) {
	t.Helper()

	c := new(compressor)

	for n, test := range [...]string{
		"HELLO",
		"Hello, World!",
		"Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!",
	} {
		bs, err := c.Compress(test)
		if err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}

		h := fn(test, &c.buf)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/file", nil)

		r.Header.Set("Accept-Encoding", "identity")
		h.ServeHTTP(w, r)

		if res := w.Result(); res.StatusCode != http.StatusOK {
			t.Errorf("test %d.1: expecting code %v, got %v", n+1, http.StatusOK, res.StatusCode)

			continue
		} else if b, err := io.ReadAll(res.Body); err != nil {
			t.Errorf("test %d.1: unexpected error: %s", n+1, err)

			continue
		} else if res.ContentLength != int64(len(test)) {
			fmt.Println(res.ContentLength)
			t.Errorf("test %d.1: expecting to read %d bytes, read %d", n+1, len(test), res.ContentLength)

			continue
		} else if string(b) != test {
			t.Errorf("test %d.1: expecting to read %q, got %q", n+1, test, b)

			continue
		}

		r.Header.Set("Accept-Encoding", "gzip")

		w = httptest.NewRecorder()

		h.ServeHTTP(w, r)

		if b, err := io.ReadAll(w.Result().Body); err != nil {
			t.Errorf("test %d.2: unexpected error: %s", n+1, err)
		} else if !bytes.Equal(b, bs) {
			t.Errorf("test %d.2: expecting to read %v, got %v", n+1, bs, b)
		}
	}
}

func TestBuffer(t *testing.T) {
	test(t, func(test string, buf *bytes.Buffer) http.Handler {
		return HandleBuffer("data.txt", buf.Bytes(), len(test), time.Now())
	})
}

func TestBufferNoSize(t *testing.T) {
	test(t, func(test string, buf *bytes.Buffer) http.Handler {
		return HandleBuffer("data.txt", buf.Bytes(), 0, time.Now())
	})
}

func TestReader(t *testing.T) {
	test(t, func(test string, buf *bytes.Buffer) http.Handler {
		return HandleReader("data.txt", buf, buf.Len(), len(test), time.Now())
	})
}

func TestReaderNoSize(t *testing.T) {
	test(t, func(_ string, buf *bytes.Buffer) http.Handler {
		return HandleReader("data.txt", buf, 0, 0, time.Now())
	})
}
