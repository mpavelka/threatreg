package main

import (
	"log"
	"threatreg/cmd"
	"threatreg/internal/config"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Execute the root command
	cmd.Execute()
}
