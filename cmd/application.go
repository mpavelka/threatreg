package cmd

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Manage applications",
}

var applicationGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of an application",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		appUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		application, err := service.GetApplication(appUUID)
		if err != nil {
			fmt.Printf("Error retrieving application: %v\n", err)
			return
		}
		fmt.Printf("Application details: uuid=%s, instance_of=%s, product_name=%s\n", 
			application.ID, application.InstanceOf, application.Product.Name)
	},
}

var applicationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new application instance of a product",
	Run: func(cmd *cobra.Command, args []string) {
		instanceOf, _ := cmd.Flags().GetString("instance-of")
		if instanceOf == "" {
			fmt.Println("Error: instance-of is required")
			return
		}
		productUUID, err := uuid.Parse(instanceOf)
		if err != nil {
			fmt.Println("Error: invalid instance-of ID (must be a uuid)")
			return
		}
		application, err := service.CreateApplication(productUUID)
		if err != nil {
			fmt.Printf("Error creating application: %v\n", err)
			return
		}
		fmt.Printf("Application created: uuid=%s, instance_of=%s\n", application.ID, application.InstanceOf)
	},
}

var applicationUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing application",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		appUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}

		var instanceOf *uuid.UUID
		instanceOfStr, err := cmd.Flags().GetString("instance-of")
		if err == nil && instanceOfStr != "" {
			parsedUUID, err := uuid.Parse(instanceOfStr)
			if err != nil {
				fmt.Println("Error: invalid instance-of ID (must be a uuid)")
				return
			}
			instanceOf = &parsedUUID
		}

		application, err := service.UpdateApplication(appUUID, instanceOf)
		if err != nil {
			fmt.Printf("Error updating application: %v\n", err)
			return
		}
		fmt.Println("✅ Application updated:")
		fmt.Printf("- uuid=%s, instance_of=%s, product_name=%s\n", 
			application.ID, application.InstanceOf, application.Product.Name)
	},
}

var applicationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an application",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		appUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		err = service.DeleteApplication(appUUID)
		if err != nil {
			fmt.Printf("Error deleting application: %v\n", err)
			return
		}
		fmt.Printf("Application deleted: uuid=%s\n", appUUID)
	},
}

var applicationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Run: func(cmd *cobra.Command, args []string) {
		productID, _ := cmd.Flags().GetString("product-id")
		
		var applications []models.Application
		var err error
		
		if productID != "" {
			productUUID, err := uuid.Parse(productID)
			if err != nil {
				fmt.Println("Error: invalid product-id (must be a uuid)")
				return
			}
			applications, err = service.ListApplicationsByProductID(productUUID)
		} else {
			applications, err = service.ListApplications()
		}
		
		if err != nil {
			fmt.Printf("Error listing applications: %v\n", err)
			return
		}
		
		if productID != "" {
			fmt.Printf("Applications for product %s:\n", productID)
		} else {
			fmt.Println("Applications:")
		}
		
		for _, application := range applications {
			fmt.Printf("- uuid=%s, instance_of=%s, product_name=%s\n", 
				application.ID, application.InstanceOf, application.Product.Name)
		}
	},
}

func init() {
	applicationCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		database.Connect()
	}

	applicationGetCmd.Flags().String("id", "", "UUID of the application (required)")
	applicationCreateCmd.Flags().String("instance-of", "", "UUID of the product this application is an instance of (required)")
	applicationUpdateCmd.Flags().String("id", "", "UUID of the application (required)")
	applicationUpdateCmd.Flags().String("instance-of", "", "New product UUID this application is an instance of")
	applicationDeleteCmd.Flags().String("id", "", "UUID of the application (required)")
	applicationListCmd.Flags().String("product-id", "", "Filter applications by product UUID (optional)")

	applicationCmd.AddCommand(applicationGetCmd)
	applicationCmd.AddCommand(applicationCreateCmd)
	applicationCmd.AddCommand(applicationUpdateCmd)
	applicationCmd.AddCommand(applicationDeleteCmd)
	applicationCmd.AddCommand(applicationListCmd)
}