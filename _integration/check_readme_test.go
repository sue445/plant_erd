package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestReadme(t *testing.T) {
	tests := []struct {
		subCommand string
	}{
		{
			subCommand: "sqlite3",
		},
		{
			subCommand: "mysql",
		},
		{
			subCommand: "postgresql",
		},
		{
			subCommand: "oracle",
		},
	}

	readme := readFile("../README.md")
	for _, tt := range tests {
		t.Run(tt.subCommand, func(t *testing.T) {
			out, err := exec.Command("../bin/plant_erd", tt.subCommand, "--help").Output()

			if assert.NoError(t, err) {
				subCommandHelp := strings.TrimSpace(string(out))
				assert.Contains(t, readme, subCommandHelp)
			}
		})
	}
}

func readFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
