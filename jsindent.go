package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pmezard/go-difflib/difflib"
)

var (
	flagDebug = flag.Bool("debug", false, "Turns on debug mode")
	flagWrite = flag.Bool("w", false, "Writes the file instead of printing to stdout")
	flagDiff  = flag.Bool("diff", false, "Writes a diff to stdout.")
)

func main() {
	flag.Parse()

	if *flagWrite && *flagDiff {
		fmt.Println("diff and w are exclusive flags")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to read file:", err)
		os.Exit(1)
	}

	indented := doIndent(string(file))

	if *flagDiff {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(file)),
			B:        difflib.SplitLines(indented),
			FromFile: filepath.Base(filename),
			FromDate: time.Now().Format("2006-01-02 15:04:05"),
			ToFile:   filepath.Base(filename),
			ToDate:   time.Now().Format("2006-01-02 15:04:05"),
			Context:  3,
		}

		result, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			fmt.Println("Failed to generate diff:", err)
			os.Exit(1)
		}

		fmt.Println(strings.Replace(result, "\t", "  ", -1))
	} else if *flagWrite {
		if err := ioutil.WriteFile(filename, []byte(indented), 664); err != nil {
			fmt.Printf("Failed to write to file (%s): %v", filename, err)
		}
	} else {
		fmt.Print(indented)
	}
}
