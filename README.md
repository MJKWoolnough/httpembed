# httpembed
--
    import "vimagination.zapto.org/httpembed"

Package httpembed aids with handling compressed 'embed' buffers, turning them
into HTTP Handlers

## Usage

#### func  HandleBuffer

```go
func HandleBuffer(name string, compressed []byte, size int, lastMod time.Time) http.Handler
```
HandleBuffer takes filename, a gzip compressed data buffer, its decompressed
size, and a last modified date, and turns it into a handler that will detect
whether the client can handle the compressed data and send the data accordingly.

If the decompressed size is 0, the decomplress buffer will be dynamically
allocated.
