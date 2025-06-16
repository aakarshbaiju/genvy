package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const GENVY_CONFIG_FILE = ".genvy.config.json"

func Ask(prompt string) bool {
	var input string

	for {
		fmt.Printf("%s (yes|y|no|n)", prompt)
		fmt.Scanln(&input)

		if strings.ToLower(input) == "yes" || strings.ToLower(input) == "y" {
			return true
		}

		if strings.ToLower(input) == "no" || strings.ToLower(input) == "n" {
			return false
		}
	}
}

func CheckIfAlreadyGitignored() bool {
	f, err := os.Open(".gitignore")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == GENVY_CONFIG_FILE {
			return true
		}
	}

	return false
}

func AddConfigToGitIgnore() {
	if CheckIfAlreadyGitignored() {
		fmt.Printf("%s already in .gitignore, skipping operation\n", GENVY_CONFIG_FILE)
		return
	}

	f, err := os.OpenFile("./.gitignore", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err = f.WriteString(fmt.Sprintf("\n%s\n", GENVY_CONFIG_FILE)); err != nil {
		panic(err)
	}
	fmt.Println("Added .genvy.config.json to .gitignore")
}
