package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "threatreg",
	Short: "A Go CLI application for threat registry management",
	Long: `Threatreg - A complete Go application for threat registry management with database migrations and CLI interface.
	
Features:
- Database migrations using golang-migrate
- User management
- SQLite and PostgreSQL support
- Configuration management`,
	// No Run function - this allows subcommands to work properly
}

// Declare all subcommands here so they're available to all files in the cmd package
var (
	dbCmd     *cobra.Command
	userCmd   *cobra.Command
	statusCmd *cobra.Command
	serveCmd  *cobra.Command
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Initialize commands (will be defined in their respective files)
	initDBCmd() // Only has db setup now
	// initUserCmd()
	initStatusCmd()

	// Add subcommands to root
	rootCmd.AddCommand(dbCmd)
	// rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(serveCmd)
}
