package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	startLinePattern = regexp.MustCompile(`^goroutine\s+(\d+)\s+\[(.*)\]:$`)
)

func load(fn string) (*GoroutineDump, error) {
	fn = strings.Trim(fn, "\"")

	if strings.HasPrefix(fn, "~") {
		home, _ := os.UserHomeDir()
		fn = filepath.Join(home, fn[1:])
	}
	fn = os.ExpandEnv(fn)

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dump := NewGoroutineDump()
	var goroutine *Goroutine

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if startLinePattern.MatchString(line) {
			goroutine, err = NewGoroutine(line)
			if err != nil {
				return nil, err
			}
			dump.Add(goroutine)
		} else if line == "" {
			// End of a goroutine section.
			if goroutine != nil {
				goroutine.Freeze()
			}
			goroutine = nil
		} else if goroutine != nil {
			goroutine.AddLine(line)
		}
	}

	if goroutine != nil {
		goroutine.Freeze()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return dump, nil
}
