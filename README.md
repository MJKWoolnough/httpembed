# httpembed

[![CI](https://github.com/MJKWoolnough/httpembed/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/httpembed/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/httpembed.svg)](https://pkg.go.dev/vimagination.zapto.org/httpembed)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/httpembed)](https://goreportcard.com/report/vimagination.zapto.org/httpembed)

--
    import "vimagination.zapto.org/httpembed"

Package httpembed aids with handling compressed 'embed' buffers and FSs, turning them into HTTP Handlers.

## Highlights

 - `HandleBuffer` function to automatically decompress gzip'd data into an HTTP Handler that will automatically server either compressed or decompressed data based on `Accept-Encoding` header.
 - `DecompressFS` function that decompresses gzip files in a `fs.FS` into a new `fs.FS`. Can be combined with `vimagination.zapto.org/httpgzip` to automatically server compressed or decompressed data based on `Accept-Encoding` header.

## Usage

```go
package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"vimagination.zapto.org/httpembed"
)

//go:embed hw.gz
var data []byte

func main() {
	handler := httpembed.HandleBuffer("hw", data, 14, time.Now())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-encoding", "identity")

	handler.ServeHTTP(w, r)

	fmt.Println(w.Body)

	// Output:
	// Hello, World!
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/httpembed
