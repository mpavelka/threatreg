package cmd

import (
	"fmt"
	"threatreg/internal/config"  
	"threatreg/internal/restapi"

	"github.com/spf13/cobra"
)

var restapiCmd *cobra.Command

// initRestapiCmd initializes the REST API command
func initRestapiCmd() {
	restapiCmd = &cobra.Command{
		Use:   "restapi",
		Short: "Start the REST API server",
		Long: `Start the REST API server for programmatic access to Threatreg.

The server provides RESTful endpoints for managing products, instances, 
threats, controls, domains, and tags. It supports JSON request/response 
format and includes CORS support for browser-based clients.

Configuration:
  The server can be configured using environment variables:
  - APP_API_HOST: Server host (default: localhost)
  - APP_API_PORT: Server port (default: 8080)

Examples:
  # Start with default settings (localhost:8080)
  threatreg restapi

  # Start with custom port
  APP_API_PORT=3000 threatreg restapi

  # Start with custom host and port  
  APP_API_HOST=0.0.0.0 APP_API_PORT=8080 threatreg restapi`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration before running
			if err := config.Load(); err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("🔧 Starting REST API server...\n")
			fmt.Printf("📊 Configuration:\n")
			fmt.Printf("   Environment: %s\n", config.GetEnvironment())
			fmt.Printf("   API Host: %s\n", config.GetAPIHost())
			fmt.Printf("   API Port: %s\n", config.GetAPIPort())
			fmt.Printf("   Database: %s\n", config.BuildDatabaseURL())
			fmt.Println()

			// Start the REST API server
			if err := restapi.RunServer(); err != nil {
				return fmt.Errorf("failed to start REST API server: %w", err)
			}

			return nil
		},
	}

	// Add flags for host and port override
	restapiCmd.Flags().StringP("host", "H", "", "API server host (overrides APP_API_HOST)")
	restapiCmd.Flags().StringP("port", "p", "", "API server port (overrides APP_API_PORT)")

	// Bind flags to configuration (this allows CLI flags to override config/env vars)
	restapiCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Load configuration first
		config.Load()

		// Override with CLI flags if provided
		if host, _ := cmd.Flags().GetString("host"); host != "" {
			// We can't modify the AppConfig directly, so we'll update the viper instance
			// This is a bit of a workaround, but it works with the existing config pattern
			config.AppConfig.APIHost = host
		}
		if port, _ := cmd.Flags().GetString("port"); port != "" {
			config.AppConfig.APIPort = port
		}
	}
}