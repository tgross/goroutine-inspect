package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mattn/go-colorable"
	"github.com/shoenig/test/must"
)

const dummyGoroutineTmpl = `goroutine %d [running]:
testing.tRunner.func1.2({0x542140, 0xc0002a4000})
        /usr/local/go/src/testing/testing.go:1631 +0x24a
testing.tRunner.func1()`

func TestShowOffset(t *testing.T) {

	var buf bytes.Buffer
	dump := NewGoroutineDump(colorable.NewNonColorable(&buf))
	for i := 0; i < 20; i++ {
		gr, err := NewGoroutine(fmt.Sprintf(dummyGoroutineTmpl, i))
		if err != nil {
			t.Fatal(err)
		}
		dump.goroutines = append(dump.goroutines, gr)

	}

	getIDs := func(t *testing.T, buf bytes.Buffer) []int {
		got := []int{}
		out, err := loadFrom(&buf)
		must.NoError(t, err)
		for _, goroutine := range out.goroutines {
			got = append(got, goroutine.id)
		}
		return got
	}

	dump.Show(0, 10)
	must.Eq(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, getIDs(t, buf))

	buf.Reset()
	dump.Show(0, 25)
	must.Eq(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
	}, getIDs(t, buf))

	buf.Reset()
	dump.Show(10, 5)
	must.Eq(t, []int{10, 11, 12, 13, 14}, getIDs(t, buf))

	buf.Reset()
	dump.Show(10, 20)
	must.Eq(t, []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, getIDs(t, buf))

}
