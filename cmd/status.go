package cmd

import (
	"fmt"
	"threatreg/internal/config"
	"threatreg/internal/database"
	"os"

	"github.com/spf13/cobra"
)

// initStatusCmd initializes the status and serve commands
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

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the application server",
		Long:  "Starts the HTTP server (implement your server logic here)",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")
			
			fmt.Printf("🚀 Starting server on port %s...\n", port)
			fmt.Println("Environment:", config.GetEnvironment())
			
			// Connect to database
			if err := database.Connect(); err != nil {
				fmt.Printf("❌ Failed to connect to database: %v\n", err)
				os.Exit(1)
			}
			defer database.Close()
			
			fmt.Println("✅ Database connected")
			
			// TODO: Implement your HTTP server here
			// Example with gorilla/mux or gin-gonic
			fmt.Printf("🌐 Server would start on http://localhost:%s\n", port)
			fmt.Println("📝 Implement your HTTP server logic in cmd/status.go")
			
			// For now, just keep the process running
			select {}
		},
	}

	// Add flags to serve command
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
}
