package cmd

import (
	"fmt"
	"strconv"

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
		fmt.Printf("Invalid index argument: %v\n", err)
		return
	}

	err = MainConfig.RemoveEntry(int(index))
	if err != nil {
		fmt.Println("Failed to remove entry", err)
		return
	}

	fmt.Printf("Successfully removed entry at index %d\n", index)
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
