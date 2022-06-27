// Package httpembed aids with handling compressed 'embed' buffers, turning them into HTTP Handlers
package httpembed

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"time"

	"vimagination.zapto.org/httpencoding"
)

var isGzip = httpencoding.HandlerFunc(func(enc httpencoding.Encoding) bool { return enc == "gzip" })

type buffers struct {
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
	http.ServeContent(w, r, "index.html", b.modTime, br)
}

// HandleBuffer takes a gzip compressed data buffer, its decompressed size, and
// a last modified data, and turns it into a handler that will detect whether
// the client can handle the compressed data and send the data accordingly.
func HandleBuffer(compressed []byte, size int, lastMod time.Time) http.Handler {
	g, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		panic(err)
	}
	uncompressed := make([]byte, size)
	if _, err := g.Read(uncompressed); err != nil {
		panic(err)
	}
	return &buffers{
		compressed:   compressed,
		uncompressed: uncompressed,
		modTime:      lastMod,
	}
}
