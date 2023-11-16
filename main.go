package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var dataFile string

func init() {
	flag.StringVar(&dataFile, "file", "housework.db", "data file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
		Usage: %s [flags] [add chore, ...|complete #]
		add			add comma-separated chores
		complete	complete designated chore

		Flags:
		`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {

}