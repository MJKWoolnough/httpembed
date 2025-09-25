package httpembed_test

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

func Example() {
	handler := httpembed.HandleBuffer("hw", data, 14, time.Now())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-encoding", "identity")

	handler.ServeHTTP(w, r)

	fmt.Println(w.Body)

	// Output:
	// Hello, World!
}
