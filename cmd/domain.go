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
		if len(domain.Instances) > 0 {
			fmt.Printf("Instances (%d):\n", len(domain.Instances))
			for _, instance := range domain.Instances {
				productName := ""
				if instance.Product.Name != "" {
					productName = fmt.Sprintf(" (Product: %s)", instance.Product.Name)
				}
				fmt.Printf("  - %s: %s%s\n", instance.ID, instance.Name, productName)
			}
		} else {
			fmt.Println("No instances associated with this domain")
		}
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
			instanceCount := len(domain.Instances)
			fmt.Printf("- uuid=%s, name=%s, description=%s, instances=%d\n", domain.ID, domain.Name, domain.Description, instanceCount)
		}
	},
}

var domainAddInstanceCmd = &cobra.Command{
	Use:   "add-instance",
	Short: "Add an instance to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domainID, _ := cmd.Flags().GetString("domain-id")
		instanceID, _ := cmd.Flags().GetString("instance-id")

		domainUUID, err := uuid.Parse(domainID)
		if err != nil {
			fmt.Println("Error: invalid domain ID (must be a uuid)")
			return
		}

		instanceUUID, err := uuid.Parse(instanceID)
		if err != nil {
			fmt.Println("Error: invalid instance ID (must be a uuid)")
			return
		}

		err = service.AddInstanceToDomain(domainUUID, instanceUUID)
		if err != nil {
			fmt.Printf("Error adding instance to domain: %v\n", err)
			return
		}
		fmt.Printf("✅ Instance %s added to domain %s\n", instanceUUID, domainUUID)
	},
}

var domainRemoveInstanceCmd = &cobra.Command{
	Use:   "remove-instance",
	Short: "Remove an instance from a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domainID, _ := cmd.Flags().GetString("domain-id")
		instanceID, _ := cmd.Flags().GetString("instance-id")

		domainUUID, err := uuid.Parse(domainID)
		if err != nil {
			fmt.Println("Error: invalid domain ID (must be a uuid)")
			return
		}

		instanceUUID, err := uuid.Parse(instanceID)
		if err != nil {
			fmt.Println("Error: invalid instance ID (must be a uuid)")
			return
		}

		err = service.RemoveInstanceFromDomain(domainUUID, instanceUUID)
		if err != nil {
			fmt.Printf("Error removing instance from domain: %v\n", err)
			return
		}
		fmt.Printf("✅ Instance %s removed from domain %s\n", instanceUUID, domainUUID)
	},
}

var domainListInstancesCmd = &cobra.Command{
	Use:   "list-instances",
	Short: "List all instances in a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domainID, _ := cmd.Flags().GetString("domain-id")

		domainUUID, err := uuid.Parse(domainID)
		if err != nil {
			fmt.Println("Error: invalid domain ID (must be a uuid)")
			return
		}

		instances, err := service.GetInstancesByDomain(domainUUID)
		if err != nil {
			fmt.Printf("Error listing instances in domain: %v\n", err)
			return
		}

		if len(instances) == 0 {
			fmt.Printf("No instances found in domain %s\n", domainUUID)
			return
		}

		fmt.Printf("Instances in domain %s:\n", domainUUID)
		for _, instance := range instances {
			productName := ""
			if instance.Product.Name != "" {
				productName = fmt.Sprintf(" (Product: %s)", instance.Product.Name)
			}
			fmt.Printf("- uuid=%s, name=%s%s\n", instance.ID, instance.Name, productName)
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
	domainAddInstanceCmd.Flags().String("domain-id", "", "UUID of the domain (required)")
	domainAddInstanceCmd.Flags().String("instance-id", "", "UUID of the instance (required)")
	domainRemoveInstanceCmd.Flags().String("domain-id", "", "UUID of the domain (required)")
	domainRemoveInstanceCmd.Flags().String("instance-id", "", "UUID of the instance (required)")
	domainListInstancesCmd.Flags().String("domain-id", "", "UUID of the domain (required)")

	domainCmd.AddCommand(domainGetCmd)
	domainCmd.AddCommand(domainCreateCmd)
	domainCmd.AddCommand(domainUpdateCmd)
	domainCmd.AddCommand(domainDeleteCmd)
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainAddInstanceCmd)
	domainCmd.AddCommand(domainRemoveInstanceCmd)
	domainCmd.AddCommand(domainListInstancesCmd)
}
