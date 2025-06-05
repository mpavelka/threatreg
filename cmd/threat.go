package cmd

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var threatCmd = &cobra.Command{
	Use:   "threat",
	Short: "Manage threats",
}

var threatGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a threat",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		threatUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		threat, err := service.GetThreat(threatUUID)
		if err != nil {
			fmt.Printf("Error retrieving threat: %v\n", err)
			return
		}
		fmt.Printf("Threat details: uuid=%s, title=%s, description=%s\n", 
			threat.ID, threat.Title, threat.Description)
	},
}

var threatCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new threat",
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		if title == "" {
			fmt.Println("Error: title is required")
			return
		}
		threat, err := service.CreateThreat(title, description)
		if err != nil {
			fmt.Printf("Error creating threat: %v\n", err)
			return
		}
		fmt.Printf("Threat created: uuid=%s, title=%s, description=%s\n", 
			threat.ID, threat.Title, threat.Description)
	},
}

var threatUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing threat",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		threatUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}

		var title *string
		titleStr, err := cmd.Flags().GetString("title")
		if err == nil && titleStr != "" {
			title = &titleStr
		}

		var description *string
		descriptionStr, err := cmd.Flags().GetString("description")
		if err == nil && descriptionStr != "" {
			description = &descriptionStr
		}

		threat, err := service.UpdateThreat(threatUUID, title, description)
		if err != nil {
			fmt.Printf("Error updating threat: %v\n", err)
			return
		}
		fmt.Println("✅ Threat updated:")
		fmt.Printf("- uuid=%s, title=%s, description=%s\n", 
			threat.ID, threat.Title, threat.Description)
	},
}

var threatDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a threat",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		threatUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		err = service.DeleteThreat(threatUUID)
		if err != nil {
			fmt.Printf("Error deleting threat: %v\n", err)
			return
		}
		fmt.Printf("Threat deleted: uuid=%s\n", threatUUID)
	},
}

var threatListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all threats",
	Run: func(cmd *cobra.Command, args []string) {
		threats, err := service.ListThreats()
		if err != nil {
			fmt.Printf("Error listing threats: %v\n", err)
			return
		}
		fmt.Println("Threats:")
		for _, threat := range threats {
			fmt.Printf("- uuid=%s, title=%s, description=%s\n", 
				threat.ID, threat.Title, threat.Description)
		}
	},
}

func init() {
	threatCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		database.Connect()
	}

	threatGetCmd.Flags().String("id", "", "UUID of the threat (required)")
	threatCreateCmd.Flags().String("title", "", "Title of the threat (required)")
	threatCreateCmd.Flags().String("description", "", "Description of the threat")
	threatUpdateCmd.Flags().String("id", "", "UUID of the threat (required)")
	threatUpdateCmd.Flags().String("title", "", "New title of the threat")
	threatUpdateCmd.Flags().String("description", "", "New description of the threat")
	threatDeleteCmd.Flags().String("id", "", "UUID of the threat (required)")

	threatCmd.AddCommand(threatGetCmd)
	threatCmd.AddCommand(threatCreateCmd)
	threatCmd.AddCommand(threatUpdateCmd)
	threatCmd.AddCommand(threatDeleteCmd)
	threatCmd.AddCommand(threatListCmd)
}