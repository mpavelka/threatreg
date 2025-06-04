package cmd

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Manage products",
}
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a product",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		product, err := service.GetProduct(uuid)
		if err != nil {
			fmt.Printf("Error retrieving product: %v\n", err)
			return
		}
		fmt.Printf("Product details: uuid=%s, name=%s, description=%s\n", product.ID, product.Name, product.Description)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new product",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		if name == "" {
			fmt.Println("Error: name is required")
			return
		}
		product, err := service.CreateProduct(name, description)
		if err != nil {
			fmt.Printf("Error creating product: %v\n", err)
			return
		}
		fmt.Printf("Product created: name=%s, description=%s\n", product.Name, product.Description)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing product",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		var name *string
		nameStr, err := cmd.Flags().GetString("name")
		if err == nil && nameStr != "" {
			name = &nameStr
		}

		var description *string
		descriptionStr, err := cmd.Flags().GetString("description")
		if err == nil && descriptionStr != "" {
			description = &descriptionStr
		}

		product, err := service.UpdateProduct(uuid, name, description)
		if err != nil {
			fmt.Printf("Error updating product: %v\n", err)
			return
		}
		fmt.Println("✅ Product updated:")
		fmt.Printf("- uuid=%s, name=%s, description=%s\n", product.ID, product.Name, product.Description)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a product",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		err = service.DeleteProduct(uuid)
		if err != nil {
			fmt.Printf("Error deleting product: %v\n", err)
			return
		}
		fmt.Printf("Product deleted: uuid=%s\n", uuid)
	},
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all products",
	Run: func(cmd *cobra.Command, args []string) {
		products, err := service.ListProducts()
		if err != nil {
			fmt.Printf("Error listing products: %v\n", err)
			return
		}
		fmt.Println("Products:")
		for _, product := range products {
			fmt.Printf("- uuid=%s, name=%s, description=%s\n", product.ID, product.Name, product.Description)
		}
	},
}

func init() {
	productCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		database.Connect()
	}

	getCmd.Flags().String("uuid", "", "UUID of the product (required)")
	createCmd.Flags().String("name", "", "Name of the product (required)")
	createCmd.Flags().String("description", "", "Description of the product")
	updateCmd.Flags().String("id", "", "UUID of the product (required)")
	updateCmd.Flags().String("name", "", "New name of the product")
	updateCmd.Flags().String("description", "", "New description of the product")
	deleteCmd.Flags().String("id", "", "UUID of the product (required)")

	productCmd.AddCommand(getCmd)
	productCmd.AddCommand(createCmd)
	productCmd.AddCommand(updateCmd)
	productCmd.AddCommand(deleteCmd)
	productCmd.AddCommand(listCmd)
}
