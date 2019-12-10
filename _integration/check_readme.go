package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		println("Usage: go run check_readme.go SOURCE_FILE EXPECTED_FILE")
		return
	}
	sourceFile := os.Args[1]
	expectedFile := os.Args[2]

	source := strings.TrimSpace(readFile(sourceFile))
	expected := strings.TrimSpace(readFile(expectedFile))

	if !strings.Contains(source, expected) {
		log.Fatalf("Expected: %s is cotains %s, but not.\n", sourceFile, expectedFile)
	}
}

func readFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
