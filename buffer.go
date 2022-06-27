package httpembed

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"time"

	"vimagination.zapto.org/httpencoding"
	"vimagination.zapto.org/memio"
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

func HandleBuffer(compressed []byte, size int, lastMod time.Time) http.Handler {
	uncompressed := make(memio.Buffer, 0, size)
	g, err := gzip.NewReader(&uncompressed)
	if err != nil {
		panic(err)
	}
	if _, err := g.Read(compressed); err != nil {
		panic(err)
	}
	return &buffers{
		compressed:   compressed,
		uncompressed: uncompressed,
		modTime:      lastMod,
	}
}
