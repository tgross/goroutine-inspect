package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/peterh/liner"
)

var (
	assignPattern = regexp.MustCompile(`^\s*[_a-zA-Z][_a-zA-Z0-9]*(\s*,\s*[_a-zA-Z][_a-zA-Z0-9]*)*\s*=\s*.*$`)
	cdPattern     = regexp.MustCompile(`^\s*cd\s*.*$`)

	commands = map[string]string{
		"?":     "Show this help",
		"cd":    "Change current working directory",
		"clear": "Clear the workspace",
		"exit":  "Exit the interactive shell",
		"help":  "Show this help",
		"ls":    "Show files in current directory",
		"pwd":   "Show current working directory",
		"quit":  "Quit the interactive shell",
		"whos":  "Show all varaibles in workspace",
	}
	cmds []string
	line *liner.State

	workspace = map[string]*GoroutineDump{}
)

func init() {
	cmds = make([]string, 0, len(commands))
	for k := range commands {
		cmds = append(cmds, k)
	}
	sort.Strings(cmds)
}

func main() {
	var err error
	line, err = createLiner()
	if err != nil {
		fmt.Println("could not read history file: %w", err)
	}

	defer func() {
		err := saveLiner(line)
		if err != nil {
			fmt.Println("could not save history file: %w", err)
		}
		line.Close()
	}()

	for {
		if cmd, err := line.Prompt(">> "); err == nil {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				continue
			}
			line.AppendHistory(cmd)

			switch cmd {
			case "?", "help":
				printHelp()
			case "clear":
				workspace = map[string]*GoroutineDump{}
				fmt.Println("Workspace cleared.")
			case "exit", "quit":
				return
			case "ls":
				wd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					continue
				}
				printDir(wd)
			case "pwd":
				wd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(wd)
			case "whos":
				if len(workspace) == 0 {
					fmt.Println("No variables defined.")
					continue
				}
				for k := range workspace {
					fmt.Printf("%s\t", k)
				}
				fmt.Println()
			default:
				if cdPattern.MatchString(cmd) {
					// Change directory.
					idx := strings.Index(cmd, "cd")
					dir := strings.TrimSpace(cmd[idx+2:])
					if dir == "" {
						fmt.Println("Expect command \"cd <dir>\"")
						continue
					}
					if strings.HasPrefix(dir, "~") {
						home, _ := os.UserHomeDir()
						dir = filepath.Join(home, dir[1:])
					}
					dir = os.ExpandEnv(dir)

					if err := os.Chdir(dir); err != nil {
						fmt.Println(err)
					}
					continue
				}

				// Assignment.
				if assignPattern.MatchString(cmd) {
					if err := assign(cmd); err != nil {
						fmt.Printf("Error, %s.\n", err.Error())
					}
					continue
				}

				if err := expr(cmd); err != nil {
					fmt.Printf("Error, %s.\n", err.Error())
				}
			}
		} else if err == liner.ErrPromptAborted || err == io.EOF {
			fmt.Println()
			break
		} else {
			log.Print("Error reading line: ", err)
		}
	}
}

func printDir(wd string) {
	f, err := os.Open(wd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(fis, func(i, j int) bool {
		return fis[i].Name() < fis[j].Name()
	})

	for _, fi := range fis {
		if fi.IsDir() {
			fmt.Printf("%s%s%s\n", fgBlue, fi.Name(), reset)
		} else {
			fmt.Println(fi.Name())
		}
	}
}

func printHelp() {
	fmt.Println("Commands:")
	for _, k := range cmds {
		fmt.Printf("  %12s: %s\n", k, commands[k])
	}
	fmt.Println()
	fmt.Println("Statements:")
	fmt.Println("\t<var>")
	fmt.Println("\t<var> = load(\"<file-name>\")")
	fmt.Println("\t<var> = <another-var>")
	fmt.Println("\t<var> = <another-var>.copy()")
	fmt.Println("\t<var> = <another-var>.copy(\"<condition>\")")
	fmt.Println("\t<var>.delete(\"<condition>\")")
	fmt.Println("\tleft = <var>.diff(<another-var>)")
	fmt.Println("\tleft, common = <var>.diff(<another-var>)")
	fmt.Println("\tleft, common, right = <var>.diff(<another-var>)")
	fmt.Println("\t<var>.keep(\"<condition>\")")
	fmt.Println("\t<var>.save(\"<output-file-name>\")")
	fmt.Println("\t<var>.search(\"<condition>\")")
	fmt.Println("\t<var>.search(\"<condition>\", offset)")
	fmt.Println("\t<var>.search(\"<condition>\", offset, limit)")
	fmt.Println("\t<var>.show()")
	fmt.Println("\t<var>.show(offset)")
	fmt.Println("\t<var>.show(offset, limit)")
	fmt.Println()
}

const (
	fgRed   = "\x1b[38;05;1m"
	fgGreen = "\x1b[38;05;2m"
	fgBlue  = "\x1b[38;05;4m"
	reset   = "\x1b[0m"
)
