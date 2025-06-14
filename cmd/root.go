package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genvy",
	Short: "Simplify environment variable setup with genvy",
	Long: `genvy is a simple cli tool to quickly bootstrap environment variables for your project.
genvy uses your local.env or .env.sample file definitions to setup environment variables for your project
so that all necessary environment variables are configured before you run the application.
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
