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

func createLiner() *liner.State {
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
		line.ReadHistory(f)
		f.Close()
	}

	return line
}

func saveLiner(liner *liner.State) {
	f, err := os.Create(getHistoryFile())
	if err != nil {
		log.Fatal("Error writing history file: ", err)
	}
	defer f.Close()

	liner.WriteHistory(f)
}
