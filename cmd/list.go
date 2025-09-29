package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config entries",
	Args:  cobra.ExactArgs(0),
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	entries := MainConfig.ReadAllEntries()
	if len(entries) == 0 {
		log.Printf("No config entries found.")
		return
	}

	for index, entry := range entries {
		log.Printf("Index: %d", index)
		log.Printf("  Branch: %s", entry.Branch)
		log.Printf("  File Paths: %v", entry.FilePaths)
		//log.Infof("  Local Last Update: %d", entry.LocalLastUpdate)
		//log.Infof("  Last SHA256: %s", entry.LastSha256)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
