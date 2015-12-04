package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Unable to read file:", err)
	}

	fmt.Print(doIndent(string(file)))
}
