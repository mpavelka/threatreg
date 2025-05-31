package cmd

import (
	"fmt"
	"os"
	"threatreg/internal/database"

	"github.com/spf13/cobra"
)

// initDBCmd initializes the database command and its subcommands
func initDBCmd() {
	dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Database management commands",
		Long:  "Commands for managing database connections and direct table operations (for development)",
	}

	var dbSetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Create database tables (development only)",
		Long:  "Creates the database tables directly without migrations (for development only)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔧 Initializing database...")

			// Connect to database
			if err := database.Connect(); err != nil {
				fmt.Printf("❌ Failed to connect to database: %v\n", err)
				os.Exit(1)
			}
			defer database.Close()

			// TODO: initialize database (this should be done via migrations)
			fmt.Println("✅ Database connection verified!")
			fmt.Println("📝 Use migrations with the 'migrate' CLI tool to create tables")
		},
	}

	// Only add the setup command for development
	dbCmd.AddCommand(dbSetupCmd)
}
