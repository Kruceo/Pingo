package cmd

import (
	"fmt"
	"os"

	"pingo/pkg/app"
	"pingo/pkg/config"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [config-file]",
	Short: "Start the smoke ping service",
	Long:  "Start the smoke ping service with the provided configuration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := args[0]

		config, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		app, err := app.NewApplication(config)
		if err != nil {
			fmt.Printf("Error creating application: %v\n", err)
			os.Exit(1)
		}

		if err := app.Run(); err != nil {
			fmt.Printf("Error running application: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
