package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

func getConfDir() string {
	userDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Join(userDir, "goroutine-inspect")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	return dir
}

func getHistoryFile() string {
	return filepath.Join(getConfDir(), "history")
}

func createLiner() (*liner.State, error) {
	line := liner.NewLiner()
	line.SetCompleter(func(line string) (c []string) {
		for n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(getHistoryFile()); err == nil {
		defer f.Close()
		_, err := line.ReadHistory(f)
		if err != nil {
			return nil, err
		}
	}

	return line, nil
}

func saveLiner(liner *liner.State) error {
	f, err := os.Create(getHistoryFile())
	if err != nil {
		log.Fatal("Error writing history file: ", err)
	}
	defer f.Close()

	_, err = liner.WriteHistory(f)
	return err
}
