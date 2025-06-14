package cmd

import (
	"fmt"
	"strings"
)

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
