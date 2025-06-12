package cmd

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage instances",
}

var instanceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of an instance",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		appUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		instance, err := service.GetInstance(appUUID)
		if err != nil {
			fmt.Printf("Error retrieving instance: %v\n", err)
			return
		}
		fmt.Printf("Instance details: uuid=%s, name=%s, instance_of=%s, product_name=%s\n",
			instance.ID, instance.Name, instance.InstanceOf, instance.Product.Name)
	},
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new instance of a product",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		instanceOf, _ := cmd.Flags().GetString("instance-of")
		if name == "" {
			fmt.Println("Error: name is required")
			return
		}
		if instanceOf == "" {
			fmt.Println("Error: instance-of is required")
			return
		}
		productUUID, err := uuid.Parse(instanceOf)
		if err != nil {
			fmt.Println("Error: invalid instance-of ID (must be a uuid)")
			return
		}
		instance, err := service.CreateInstance(name, productUUID)
		if err != nil {
			fmt.Printf("Error creating instance: %v\n", err)
			return
		}
		fmt.Printf("Instance created: uuid=%s, name=%s, instance_of=%s\n",
			instance.ID, instance.Name, instance.InstanceOf)
	},
}

var instanceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing instance",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		appUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}

		var name *string
		nameStr, err := cmd.Flags().GetString("name")
		if err == nil && nameStr != "" {
			name = &nameStr
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

		instance, err := service.UpdateInstance(appUUID, name, instanceOf)
		if err != nil {
			fmt.Printf("Error updating instance: %v\n", err)
			return
		}
		fmt.Println("✅ Instance updated:")
		fmt.Printf("- uuid=%s, name=%s, instance_of=%s, product_name=%s\n",
			instance.ID, instance.Name, instance.InstanceOf, instance.Product.Name)
	},
}

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an instance",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		appUUID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		err = service.DeleteInstance(appUUID)
		if err != nil {
			fmt.Printf("Error deleting instance: %v\n", err)
			return
		}
		fmt.Printf("Instance deleted: uuid=%s\n", appUUID)
	},
}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances",
	Run: func(cmd *cobra.Command, args []string) {
		productID, _ := cmd.Flags().GetString("product-id")

		var instances []models.Instance
		var err error

		if productID != "" {
			productUUID, err := uuid.Parse(productID)
			if err != nil {
				fmt.Println("Error: invalid product-id (must be a uuid)")
				return
			}
			instances, err = service.ListInstancesByProductID(productUUID)
		} else {
			instances, err = service.ListInstances()
		}

		if err != nil {
			fmt.Printf("Error listing instances: %v\n", err)
			return
		}

		if productID != "" {
			fmt.Printf("Instances for product %s:\n", productID)
		} else {
			fmt.Println("Instances:")
		}

		for _, instance := range instances {
			fmt.Printf("- uuid=%s, name=%s, instance_of=%s, product_name=%s\n",
				instance.ID, instance.Name, instance.InstanceOf, instance.Product.Name)
		}
	},
}

func init() {
	instanceCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		database.Connect()
	}

	instanceGetCmd.Flags().String("id", "", "UUID of the instance (required)")
	instanceCreateCmd.Flags().String("name", "", "Name of the instance (required)")
	instanceCreateCmd.Flags().String("instance-of", "", "UUID of the product this instance is an instance of (required)")
	instanceUpdateCmd.Flags().String("id", "", "UUID of the instance (required)")
	instanceUpdateCmd.Flags().String("name", "", "New name of the instance")
	instanceUpdateCmd.Flags().String("instance-of", "", "New product UUID this instance is an instance of")
	instanceDeleteCmd.Flags().String("id", "", "UUID of the instance (required)")
	instanceListCmd.Flags().String("product-id", "", "Filter instances by product UUID (optional)")

	instanceCmd.AddCommand(instanceGetCmd)
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceUpdateCmd)
	instanceCmd.AddCommand(instanceDeleteCmd)
	instanceCmd.AddCommand(instanceListCmd)
}
