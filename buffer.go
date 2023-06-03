// Package httpembed aids with handling compressed 'embed' buffers, turning them into HTTP Handlers
package httpembed

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"time"

	"vimagination.zapto.org/httpencoding"
)

var isGzip = httpencoding.HandlerFunc(func(enc httpencoding.Encoding) bool { return enc == "gzip" })

type buffers struct {
	name                     string
	compressed, uncompressed []byte
	modTime                  time.Time
}

func (b *buffers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var br *bytes.Reader
	if httpencoding.HandleEncoding(r, isGzip) {
		br = bytes.NewReader(b.compressed)
		w.Header().Add("Content-Encoding", "gzip")
	} else {
		br = bytes.NewReader(b.uncompressed)
	}
	http.ServeContent(w, r, b.name, b.modTime, br)
}

// HandleBuffer takes filename, a gzip compressed data buffer, its decompressed
// size, and a last modified date, and turns it into a handler that will detect
// whether the client can handle the compressed data and send the data
// accordingly.
//
// If the decompressed size is 0, the decomplress buffer will be dynamically
// allocated.
func HandleBuffer(name string, compressed []byte, size int, lastMod time.Time) http.Handler {
	g, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		panic(err)
	}
	var uncompressed []byte
	if size == 0 {
		uncompressed, err = io.ReadAll(g)
		if err != nil {
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
