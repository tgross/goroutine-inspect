package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mattn/go-colorable"
	"github.com/shoenig/test/must"
)

const dummyGoroutineMetaTmpl = `goroutine %d [%s]:`

func TestShowOffset(t *testing.T) {

	var buf bytes.Buffer
	dump := NewGoroutineDump(colorable.NewNonColorable(&buf))
	for i := 0; i < 20; i++ {
		gr, err := NewGoroutine(fmt.Sprintf(dummyGoroutineMetaTmpl, i, "running"))
		must.NoError(t, err)
		dump.goroutines = append(dump.goroutines, gr)
	}

	getIDs := func(t *testing.T, buf bytes.Buffer) []int {
		t.Helper()
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

func TestSearchOffset(t *testing.T) {

	var buf bytes.Buffer
	dump := NewGoroutineDump(colorable.NewNonColorable(&buf))
	states := []string{"running", "select"}
	for i := 0; i < 20; i++ {
		gr, err := NewGoroutine(fmt.Sprintf(dummyGoroutineMetaTmpl, i, states[i%2]))
		must.NoError(t, err)
		dump.goroutines = append(dump.goroutines, gr)
	}

	getIDs := func(t *testing.T, buf bytes.Buffer) []int {
		t.Helper()
		got := []int{}
		out, err := loadFrom(&buf)
		must.NoError(t, err)
		for _, goroutine := range out.goroutines {
			got = append(got, goroutine.id)
		}
		return got
	}

	dump.Search("state == 'select'", 0, 30)
	must.Eq(t, []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}, getIDs(t, buf))

	buf.Reset()
	dump.Search("state == 'select'", 0, 5)
	must.Eq(t, []int{1, 3, 5, 7, 9}, getIDs(t, buf))

	buf.Reset()
	dump.Search("state == 'select'", 5, 4)
	must.Eq(t, []int{11, 13, 15, 17}, getIDs(t, buf))
}
