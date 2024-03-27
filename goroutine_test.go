package main

import (
	"fmt"
	"testing"
)

const dummyGoroutineTmpl = `goroutine %d [running]:
testing.tRunner.func1.2({0x542140, 0xc0002a4000})
        /usr/local/go/src/testing/testing.go:1631 +0x24a
testing.tRunner.func1()`

func TestShowOffset(t *testing.T) {
	dump := NewGoroutineDump()
	for i := 0; i < 20; i++ {
		gr, err := NewGoroutine(fmt.Sprintf(dummyGoroutineTmpl, i))
		if err != nil {
			t.Fatal(err)
		}
		dump.goroutines = append(dump.goroutines, gr)

	}

	// TODO: it'd be nice if we had a way of capturing output
	dump.Show(0, 10)
	dump.Show(0, 25)
	dump.Show(10, 5)
	dump.Show(10, 20)
}
