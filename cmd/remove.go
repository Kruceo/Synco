package cmd

import (
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [index]",
	Short: "Remove a config entry by its index",
	Args:  cobra.ExactArgs(1),
	Run:   remove,
}

func remove(cmd *cobra.Command, args []string) {
	index, err := strconv.ParseInt(args[0], 10, 8)
	if err != nil {
		log.Error("Invalid index argument", err)
		return
	}

	err = MainConfig.RemoveEntry(int(index))
	if err != nil {
		log.Error("Failed to remove entry", err)
		return
	}

	log.Infof("Successfully removed entry at index %d", index)
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
