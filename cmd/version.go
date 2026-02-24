package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.3"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show binary version",
	Long:  "Show the version of the pingo binary",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
