package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a .env file from a template file by prompting for missing values.",
		Long: `The generate command reads an environment template file like .env.sample or local.env and helps you create a .env file by prompting for any missing values. 
Examples Usage:
  genvy generate`,
		Run: generateEnv,
	}

	matchableTemplates []string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringSliceVarP(&matchableTemplates, "templates", "t", []string{}, "Additioinal environment file templates for which environment variables are required")
}

func generateEnv(cmd *cobra.Command, args []string) {
	fmt.Println(matchableTemplates)
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

	envFilesToCheck := []string{
		".env.sample",
		".env.local",
		"local.env",
		".env.example",
	}

	if len(matchableTemplates) > 0 {
		envFilesToCheck = append(envFilesToCheck, matchableTemplates...)
	}

	var envPaths []string

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isEnvFile(path, envFilesToCheck) {
			envPaths = append(envPaths, path)
		}
		return nil
	})

	return envPaths
}

func isEnvFile(path string, envFiles []string) bool {
	for _, envFile := range envFiles {
		if strings.HasSuffix(path, envFile) {
			return true
		}
	}
	return false
}

func generateJsonConfig(data map[string]string) {
	shouldGenerateConfig := Ask("Do you want genvy to generate a json file containing the environment variables for reference?")
	if !shouldGenerateConfig {
		return
	}

	jsonString, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./.genvy.config.json", jsonString, 0644)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("./.gitignore"); errors.Is(err, os.ErrNotExist) {
		return
	}

	shouldAddToGitIgnore := Ask("Do you want to add generated config file to .gitignore?")

	if shouldAddToGitIgnore {
		AddConfigToGitIgnore()
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
