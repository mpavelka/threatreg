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
			fmt.Println("🔧 Setting up database tables...")

			// Connect to database
			if err := database.Connect(); err != nil {
				fmt.Printf("❌ Failed to connect to database: %v\n", err)
				os.Exit(1)
			}
			defer database.Close()

			// Create tables
			if err := database.CreateTables(); err != nil {
				fmt.Printf("❌ Failed to create tables: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ Database tables created successfully!")
			fmt.Println("📝 For production, use migrations with the 'migrate' CLI tool")
		},
	}

	// Only add the setup command for development
	dbCmd.AddCommand(dbSetupCmd)
}
