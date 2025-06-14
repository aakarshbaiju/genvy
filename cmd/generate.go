package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a .env file from a template file by prompting for missing values.",
	Long: `The generate command reads an environment template file like .env.sample or local.env and helps you create a .env file by prompting for any missing values. 
Examples Usage:
  genvy generate`,
	Run: generateEnv,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func generateEnv(cmd *cobra.Command, args []string) {
	templates := findEnvTemplates()
	envVars := processTemplates(templates)
	generateJsonConfig(envVars)
	generateEnvFiles(envVars)
}

func processTemplates(templates []string) map[string]string {
	environmentVariables := make(map[string]string)

	for _, template := range templates {
		file, err := os.Open(template)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			pair := strings.Split(scanner.Text(), "=")
			if len(pair) > 0 {
				var value string

				if len(pair) > 1 && pair[1] != "" {
					value = pair[1]
				} else {
					fmt.Print(pair[0] + ":")
					fmt.Scan(&value)
				}
				environmentVariables[pair[0]] = value
			}
		}
	}

	return environmentVariables
}

func findEnvTemplates() []string {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envRegex, err := regexp.Compile(`(?i)(^|/)(\.env\.(sample|local)|local\.env)$`)
	if err != nil {
		log.Fatal(err)
	}

	var envPaths []string

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if envRegex.MatchString(path) {
			fmt.Println(path)
			envPaths = append(envPaths, path)
		}
		return nil
	})

	return envPaths
}

func generateJsonConfig(data map[string]string) {
	jsonString, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		panic("Couldn't generate genvy config file")
	}

	err = os.WriteFile("./.genvy.config.json", jsonString, 0644)

	if err != nil {
		panic("Failed to write genvy config")
	}
}

func generateEnvFiles(envVars map[string]string) {
	f, err := os.Create(".env")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for key, value := range envVars {
		f.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
}
