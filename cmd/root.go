package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pinger",
	Short: "A CLI smoke ping service",
	Long:  "A CLI tool for performing smoke ping tests and storing metrics in ClickHouse",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}