package httpembed

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	var (
		buf bytes.Buffer
		gz  = gzip.NewWriter(&buf)
	)
	for n, test := range [...]string{
		"HELLO",
		"Hello, World!",
		"Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!",
	} {
		buf.Reset()
		gz.Reset(&buf)
		if _, err := io.WriteString(gz, test); err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}
		gz.Close()
		bs := buf.Bytes()
		h := HandleBuffer("data.txt", bs, len(test), time.Now())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/file", nil)
		h.ServeHTTP(w, r)
		res := w.Result()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("test %d.1: unexpected error: %s", n+1, err)
			continue
		} else if res.ContentLength != int64(len(test)) {
			t.Errorf("test %d.1: expecting to read %d bytes, read %d", n+1, len(test), r.ContentLength)
			continue
		} else if string(b) != test {
			t.Errorf("test %d.1: expecting to read %q, got %q", n+1, test, b)
			continue
		}
		r.Header.Set("Accept-Encoding", "gzip")
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, err = io.ReadAll(w.Result().Body)
		if err != nil {
			t.Errorf("test %d.2: unexpected error: %s", n+1, err)
		} else if !bytes.Equal(b, bs) {
			t.Errorf("test %d.2: expecting to read %v, got %v", n+1, bs, b)
		}
	}
}

func TestBufferNoSize(t *testing.T) {
	var (
		buf bytes.Buffer
		gz  = gzip.NewWriter(&buf)
	)
	for n, test := range [...]string{
		"HELLO",
		"Hello, World!",
		"Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!",
	} {
		buf.Reset()
		gz.Reset(&buf)
		if _, err := io.WriteString(gz, test); err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}
		gz.Close()
		bs := buf.Bytes()
		h := HandleBuffer("data.txt", bs, 0, time.Now())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/file", nil)
		h.ServeHTTP(w, r)
		res := w.Result()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("test %d.1: unexpected error: %s", n+1, err)
			continue
		} else if res.ContentLength != int64(len(test)) {
			t.Errorf("test %d.1: expecting to read %d bytes, read %d", n+1, len(test), r.ContentLength)
			continue
		} else if string(b) != test {
			t.Errorf("test %d.1: expecting to read %q, got %q", n+1, test, b)
			continue
		}
		r.Header.Set("Accept-Encoding", "gzip")
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, err = io.ReadAll(w.Result().Body)
		if err != nil {
			t.Errorf("test %d.2: unexpected error: %s", n+1, err)
		} else if !bytes.Equal(b, bs) {
			t.Errorf("test %d.2: expecting to read %v, got %v", n+1, bs, b)
		}
	}
}

func TestReader(t *testing.T) {
	var (
		buf bytes.Buffer
		gz  = gzip.NewWriter(&buf)
	)
	for n, test := range [...]string{
		"HELLO",
		"Hello, World!",
		"Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!",
	} {
		buf.Reset()
		gz.Reset(&buf)
		if _, err := io.WriteString(gz, test); err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}
		gz.Close()
		bs := buf.Bytes()
		h := HandleReader("data.txt", &buf, buf.Len(), len(test), time.Now())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/file", nil)
		h.ServeHTTP(w, r)
		res := w.Result()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("test %d.1: unexpected error: %s", n+1, err)
			continue
		} else if res.ContentLength != int64(len(test)) {
			t.Errorf("test %d.1: expecting to read %d bytes, read %d", n+1, len(test), r.ContentLength)
			continue
		} else if string(b) != test {
			t.Errorf("test %d.1: expecting to read %q, got %q", n+1, test, b)
			continue
		}
		r.Header.Set("Accept-Encoding", "gzip")
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, err = io.ReadAll(w.Result().Body)
		if err != nil {
			t.Errorf("test %d.2: unexpected error: %s", n+1, err)
		} else if !bytes.Equal(b, bs) {
			t.Errorf("test %d.2: expecting to read %v, got %v", n+1, bs, b)
		}
	}
}

func TestReaderNoSize(t *testing.T) {
	var (
		buf bytes.Buffer
		gz  = gzip.NewWriter(&buf)
	)
	for n, test := range [...]string{
		"HELLO",
		"Hello, World!",
		"Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!  Long!",
	} {
		buf.Reset()
		gz.Reset(&buf)
		if _, err := io.WriteString(gz, test); err != nil {
			t.Fatalf("test %d: unexpected error: %s", n+1, err)
		}
		gz.Close()
		bs := buf.Bytes()
		h := HandleReader("data.txt", &buf, 0, len(test), time.Now())
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/file", nil)
		h.ServeHTTP(w, r)
		res := w.Result()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("test %d.1: unexpected error: %s", n+1, err)
			continue
		} else if res.ContentLength != int64(len(test)) {
			t.Errorf("test %d.1: expecting to read %d bytes, read %d", n+1, len(test), r.ContentLength)
			continue
		} else if string(b) != test {
			t.Errorf("test %d.1: expecting to read %q, got %q", n+1, test, b)
			continue
		}
		r.Header.Set("Accept-Encoding", "gzip")
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		b, err = io.ReadAll(w.Result().Body)
		if err != nil {
			t.Errorf("test %d.2: unexpected error: %s", n+1, err)
		} else if !bytes.Equal(b, bs) {
			t.Errorf("test %d.2: expecting to read %v, got %v", n+1, bs, b)
		}
	}
}
