package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func getFileNames() []string {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("can't get working directory: %v", err))
	}

	fileNames, err := filepath.Glob(filepath.Join(wd, "testdata/*.js"))
	if err != nil {
		panic(fmt.Sprintf("can't get working directory: %v", err))
	}

	return fileNames
}

func TestIndenter(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("can't get working directory: %v", err))
	}

	//fileNames := getFileNames()

	//for i := 1; i <= len(fileNames)/2; i++ {
	for i := 6; i <= 6; i++ {
		expectedFile := fmt.Sprintf(filepath.Join(wd, "testdata/%02d_expected.js"), i)
		inputFile := fmt.Sprintf(filepath.Join(wd, "testdata/%02d_expected.js"), i)
		expected, err := ioutil.ReadFile(expectedFile)
		if err != nil {
			t.Errorf("Cannot read file %s: %s", expectedFile, err)
		}
		input, err := ioutil.ReadFile(inputFile)
		if err != nil {
			t.Errorf("Cannot read file %s: %s", inputFile, err)
		}

		got := []byte(doIndent(string(input)))

		if bytes.Compare(got, expected) != 0 {
			t.Errorf("%02d.js does not match %02d.js:\nGot\n%s\nWanted:\n%s", i, i, got, expected)
			t.Errorf("%02d.js does not match %02d.js:\nGot\n%s\nWanted:\n%s", i, i, spew.Sdump(got), spew.Sdump(expected))
		}
	}
}
