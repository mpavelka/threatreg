package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "threatreg",
	Short: "A CLI application for threat registry management",
}

// Declare all subcommands here so they're available to all files in the cmd package
var (
	statusCmd *cobra.Command
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	initStatusCmd()
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(productCmd)
	rootCmd.AddCommand(instanceCmd)
	rootCmd.AddCommand(domainCmd)
	rootCmd.AddCommand(threatCmd)
	rootCmd.AddCommand(controlCmd)
}
