# Section Writer

SectionWriter is a Go library that allows for creating windows into a file for writing, much the same as the SectionReader in the standard library does. These SectionWriters are safe for concurrent writing as long as the windows do not overlap, it is your responsibility to make sure this is the case.

The use case for me is downloading sections of a file concurrently and writing them to sections of a local file.

SectionWriter, and NewSectionWriter, accept anything that conforms to [the io.WriterAt interface](https://golang.org/pkg/io/#WriterAt):

```go
type WriterAt interface {
        WriteAt(p []byte, off int64) (n int, err error)
}
```

## Usage

```go
package main

import (
	"fmt"
	"os"
	"sync"

	sw "github.com/wjessop/sectionwriter"
)

var data = []byte("0123456789abcdefghij")

func main() {
	f, err := os.OpenFile("somefile", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a window onto the file starting at byte 0 that is 10 bytes in length.
	s1 := sw.NewSectionWriter(f, 0, 10)

	// Create a window onto the file starting at byte 10 that is 10 bytes in length
	s2 := sw.NewSectionWriter(f, 10, 10)

	var wg sync.WaitGroup
	wg.Add(2)

	// Concurrent writes thread safe
	go func() {
		// Write to the first section
		defer wg.Done()
		n, err := s1.Write(data[0:10])
		fmt.Printf("Wrote %d bytes\n", n)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		// Write to the second section
		defer wg.Done()
		n, err := s2.Write(data[10:len(data)])
		fmt.Printf("Wrote %d bytes\n", n)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}

```

You can run this example yourself:

```
$ go run examples/example.go
Wrote 10 bytes
Wrote 10 bytes
$ cat somefile
0123456789abcdefghij
```

## License

MIT, see LICENSE file
