package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionStr = "0.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println(versionStr)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
