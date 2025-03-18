# httpembed
--
    import "vimagination.zapto.org/httpembed"

Package httpembed aids with handling compressed 'embed' buffers and FSs, turning
them into HTTP Handlers.

## Usage

#### func  DecompressFS

```go
func DecompressFS(files fs.FS) (fs.FS, error)
```
DecompressFS takes a FS with compressed (.gz) files and returns a new FS with
those files decompressed and stored under the same name with the .gz suffix
removed.

The output of this is intended to be use with httpgzip.FileServer.

#### func  HandleBuffer

```go
func HandleBuffer(name string, compressed []byte, size int, lastMod time.Time) http.Handler
```
HandleBuffer takes filename, a gzip compressed data buffer, its uncompressed
size, and a last modified date, and turns it into a handler that will detect
whether the client can handle the compressed data and send the data accordingly.

If the uncompressed size is 0, the decompress buffer will be dynamically
allocated.

#### func  HandleReader

```go
func HandleReader(name string, r io.Reader, compressedSize, uncompressedSize int, lastMod time.Time) http.Handler
```
HandleReader takes filename, a gzip compressed data buffer, its compressed and
uncompressed size, and a last modified date, and turns it into a handler that
will detect whether the client can handle the compressed data and send the data
accordingly.

If the either the compressed size or uncompressed size is 0, the buffers will be
dynamically allocated.
