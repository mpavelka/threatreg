package cmd

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains",
}

var domainGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a domain",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		domain, err := service.GetDomain(uuid)
		if err != nil {
			fmt.Printf("Error retrieving domain: %v\n", err)
			return
		}
		fmt.Printf("Domain details: uuid=%s, name=%s, description=%s\n", domain.ID, domain.Name, domain.Description)
		fmt.Println("No instances associated with this domain (Intentionally broken)")

	},
}

var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		if name == "" {
			fmt.Println("Error: name is required")
			return
		}
		domain, err := service.CreateDomain(name, description)
		if err != nil {
			fmt.Printf("Error creating domain: %v\n", err)
			return
		}
		fmt.Printf("Domain created: uuid=%s, name=%s, description=%s\n", domain.ID, domain.Name, domain.Description)
	},
}

var domainUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing domain",
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

		domain, err := service.UpdateDomain(uuid, name, description)
		if err != nil {
			fmt.Printf("Error updating domain: %v\n", err)
			return
		}
		fmt.Println("✅ Domain updated:")
		fmt.Printf("- uuid=%s, name=%s, description=%s\n", domain.ID, domain.Name, domain.Description)
	},
}

var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println("Error: invalid ID (must be a uuid)")
			return
		}
		err = service.DeleteDomain(uuid)
		if err != nil {
			fmt.Printf("Error deleting domain: %v\n", err)
			return
		}
		fmt.Printf("Domain deleted: uuid=%s\n", uuid)
	},
}

var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains",
	Run: func(cmd *cobra.Command, args []string) {
		domains, err := service.ListDomains()
		if err != nil {
			fmt.Printf("Error listing domains: %v\n", err)
			return
		}
		if len(domains) == 0 {
			fmt.Println("No domains found")
			return
		}
		fmt.Println("Domains:")
		for _, domain := range domains {
			fmt.Printf("- uuid=%s, name=%s, description=%s\n", domain.ID, domain.Name, domain.Description)
		}
	},
}

func init() {
	domainCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		database.Connect()
	}

	domainGetCmd.Flags().String("id", "", "UUID of the domain (required)")
	domainCreateCmd.Flags().String("name", "", "Name of the domain (required)")
	domainCreateCmd.Flags().String("description", "", "Description of the domain")
	domainUpdateCmd.Flags().String("id", "", "UUID of the domain (required)")
	domainUpdateCmd.Flags().String("name", "", "New name of the domain")
	domainUpdateCmd.Flags().String("description", "", "New description of the domain")
	domainDeleteCmd.Flags().String("id", "", "UUID of the domain (required)")

	domainCmd.AddCommand(domainGetCmd)
	domainCmd.AddCommand(domainCreateCmd)
	domainCmd.AddCommand(domainUpdateCmd)
	domainCmd.AddCommand(domainDeleteCmd)
	domainCmd.AddCommand(domainListCmd)
}
