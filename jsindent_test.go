package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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

	fileNames := getFileNames()

	for i := 1; i <= len(fileNames)/2; i++ {
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

		got := doIndent(string(input))

		if got != string(expected) {
			t.Errorf("%s does not match %s:\nGot\n%s\nWanted:\n%s", inputFile, expectedFile, got, expected)
		}
	}
}
