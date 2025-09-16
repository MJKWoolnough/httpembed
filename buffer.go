// Package httpembed aids with handling compressed 'embed' buffers and FSs, turning them into HTTP Handlers.
package httpembed // import "vimagination.zapto.org/python"

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"time"

	"vimagination.zapto.org/httpencoding"
)

type requestGzip bool

func (r *requestGzip) Handle(enc httpencoding.Encoding) bool {
	if enc == "gzip" || httpencoding.IsWildcard(enc) && !httpencoding.IsDisallowedInWildcard(enc, "gzip") {
		*r = true

		return true
	}

	return enc == "" || httpencoding.IsWildcard(enc) && !httpencoding.IsDisallowedInWildcard(enc, "")
}

type buffers struct {
	name                     string
	compressed, uncompressed []byte
	modTime                  time.Time
}

func (b *buffers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		br          *bytes.Reader
		requestGzip requestGzip
	)

	if !httpencoding.HandleEncoding(r, &requestGzip) {
		httpencoding.InvalidEncoding(w)

		return
	}

	if requestGzip {
		w.Header().Add("Content-Encoding", "gzip")

		br = bytes.NewReader(b.compressed)
		w = &wrapResponseWriter{
			ResponseWriter: w,
			size:           int64(len(b.compressed)),
		}
	} else {
		br = bytes.NewReader(b.uncompressed)
	}

	http.ServeContent(w, r, b.name, b.modTime, br)
}

type wrapResponseWriter struct {
	http.ResponseWriter
	size int64
}

func (w *wrapResponseWriter) WriteHeader(code int) {
	if w.Header().Get("Content-Length") == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(w.size, 10))
	}

	w.ResponseWriter.WriteHeader(code)
}

// HandleBuffer takes filename, a gzip compressed data buffer, its uncompressed
// size, and a last modified date, and turns it into a handler that will detect
// whether the client can handle the compressed data and send the data
// accordingly.
//
// If the uncompressed size is 0, the decompress buffer will be dynamically
// allocated.
func HandleBuffer(name string, compressed []byte, size int, lastMod time.Time) http.Handler {
	g, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		panic(err)
	}

	var uncompressed []byte

	if size == 0 {
		if uncompressed, err = io.ReadAll(g); err != nil {
			panic(err)
		}
	} else {
		uncompressed = make([]byte, size)

		if n, err := io.ReadFull(g, uncompressed); n != size {
			panic(err)
		}
	}

	return &buffers{
		name:         name,
		compressed:   compressed,
		uncompressed: uncompressed,
		modTime:      lastMod,
	}
}

// HandleReader takes filename, a gzip compressed data buffer, its compressed
// and uncompressed size, and a last modified date, and turns it into a handler
// that will detect whether the client can handle the compressed data and send
// the data accordingly.
//
// If the either the compressed size or uncompressed size is 0, the buffers will
// be dynamically allocated.
func HandleReader(name string, r io.Reader, compressedSize, uncompressedSize int, lastMod time.Time) http.Handler {
	var (
		compressed []byte
		err        error
	)

	if compressedSize == 0 {
		compressed, err = io.ReadAll(r)
	} else {
		compressed = make([]byte, compressedSize)
		_, err = io.ReadFull(r, compressed)
	}

	if err != nil {
		panic(err)
	}

	return HandleBuffer(name, compressed, uncompressedSize, lastMod)
}
