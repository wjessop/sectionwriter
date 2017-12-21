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
