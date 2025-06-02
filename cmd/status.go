package cmd

import (
	"fmt"
	"threatreg/internal/config"
	"threatreg/internal/database"

	"github.com/spf13/cobra"
)

// initStatusCmd initializes the status command
func initStatusCmd() {
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show application status",
		Long:  "Displays the current status of the application and database connection",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔧 DEBUG: Status command is running!") // Debug line
			fmt.Println("📊 Application Status")
			fmt.Println("====================")

			// Show configuration
			fmt.Printf("Environment: %s\n", config.GetEnvironment())

			// Show database configuration with obfuscated password
			if config.GetDatabaseURL() != "" {
				fmt.Printf("Database URL: %s\n", config.GetDatabaseURL())
			} else {
				fmt.Printf("Database Protocol: %s\n", config.GetDatabaseProtocol())
				fmt.Printf("Database Host: %s\n", config.GetDatabaseHost())
				fmt.Printf("Database Port: %s\n", config.GetDatabasePort())
				fmt.Printf("Database Name: %s\n", config.GetDatabaseName())
				fmt.Printf("Database Username: %s\n", config.GetDatabaseUsername())

				password := config.GetDatabasePassword()
				if password != "" {
					fmt.Printf("Database Password: ***\n")
				} else {
					fmt.Printf("Database Password: [NOT_SET]\n")
				}
			}

			// Test database connection
			fmt.Print("Database Connection: ")
			if err := database.Connect(); err != nil {
				fmt.Printf("❌ Failed (%v)\n", err)
				return
			}
			defer database.Close()

			fmt.Println("✅ Connected")

			fmt.Println("\n🚀 Application is running successfully!")
		},
	}
}
