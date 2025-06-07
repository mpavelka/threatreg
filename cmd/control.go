package cmd

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Manage controls",
}

var controlGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a control",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		control, err := service.GetControl(uuid)
		if err != nil {
			fmt.Printf("Error retrieving control: %v\n", err)
			return
		}
		fmt.Printf("Control details: uuid=%s, title=%s, description=%s\n", control.ID, control.Title, control.Description)
	},
}

var controlCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new control",
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		if title == "" {
			fmt.Println("Error: title is required")
			return
		}
		control, err := service.CreateControl(title, description)
		if err != nil {
			fmt.Printf("Error creating control: %v\n", err)
			return
		}
		fmt.Printf("Control created: title=%s, description=%s\n", control.Title, control.Description)
	},
}

var controlUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing control",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
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

		control, err := service.UpdateControl(uuid, title, description)
		if err != nil {
			fmt.Printf("Error updating control: %v\n", err)
			return
		}
		fmt.Println("✅ Control updated:")
		fmt.Printf("- uuid=%s, title=%s, description=%s\n", control.ID, control.Title, control.Description)
	},
}

var controlDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a control",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		err = service.DeleteControl(uuid)
		if err != nil {
			fmt.Printf("Error deleting control: %v\n", err)
			return
		}
		fmt.Printf("Control deleted: uuid=%s\n", uuid)
	},
}

var controlListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all controls",
	Run: func(cmd *cobra.Command, args []string) {
		controls, err := service.ListControls()
		if err != nil {
			fmt.Printf("Error listing controls: %v\n", err)
			return
		}
		fmt.Println("Controls:")
		for _, control := range controls {
			fmt.Printf("- uuid=%s, title=%s, description=%s\n", control.ID, control.Title, control.Description)
		}
	},
}

func init() {
	controlCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		database.Connect()
	}

	controlGetCmd.Flags().String("id", "", "UUID of the control (required)")
	controlCreateCmd.Flags().String("title", "", "Title of the control (required)")
	controlCreateCmd.Flags().String("description", "", "Description of the control")
	controlUpdateCmd.Flags().String("id", "", "UUID of the control (required)")
	controlUpdateCmd.Flags().String("title", "", "New title of the control")
	controlUpdateCmd.Flags().String("description", "", "New description of the control")
	controlDeleteCmd.Flags().String("id", "", "UUID of the control (required)")

	controlCmd.AddCommand(controlGetCmd)
	controlCmd.AddCommand(controlCreateCmd)
	controlCmd.AddCommand(controlUpdateCmd)
	controlCmd.AddCommand(controlDeleteCmd)
	controlCmd.AddCommand(controlListCmd)
}
