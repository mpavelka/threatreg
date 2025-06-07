package cmd

import (
	"threatreg/tviewapp"

	"github.com/spf13/cobra"
)

var tviewAppCmd = &cobra.Command{
	Use:   "tviewapp",
	Short: "Terminal application",
	Long:  "A command to manage threats in the threat registry application.",
	Run: func(cmd *cobra.Command, args []string) {
		tviewapp.Run()
	},
}

func init() {
	rootCmd.AddCommand(tviewAppCmd)
}
