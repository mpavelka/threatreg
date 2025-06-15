package domains

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func createSelectExistingInstanceForm(domainID uuid.UUID, onClose func()) tview.Primitive {
	form := tview.NewForm()

	// Get list of all instances
	instances, err := service.ListInstances()
	if err != nil {
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading instances: %v", err))
		return errorView
	}

	if len(instances) == 0 {
		noInstancesView := tview.NewTextView().SetText("No instances available to add to this domain.")
		return noInstancesView
	}

	// Create instance options for dropdown with fuzzy search
	instanceOptions := make([]string, len(instances))
	instanceMap := make(map[string]uuid.UUID)
	for i, instance := range instances {
		displayName := fmt.Sprintf("%s (%s)", instance.Name, instance.Product.Name)
		instanceOptions[i] = displayName
		instanceMap[displayName] = instance.ID
	}

	var selectedInstanceID uuid.UUID
	form.AddDropDown("Instance", instanceOptions, 0, func(option string, optionIndex int) {
		selectedInstanceID = instanceMap[option]
	})

	// Set initial selected instance
	if len(instances) > 0 {
		selectedInstanceID = instances[0].ID
	}

	form.AddButton("Add to Domain", func() {
		// Associate instance with domain
		err := service.AddInstanceToDomain(domainID, selectedInstanceID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}

func createNewInstanceForm(domainID uuid.UUID, onClose func()) tview.Primitive {
	form := tview.NewForm()

	nameField := ""
	var selectedProductID uuid.UUID

	// Get list of products for dropdown
	products, err := service.ListProducts()
	if err != nil {
		// If we can't load products, show error and return
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading products: %v", err))
		return errorView
	}

	// Create product options for dropdown
	productOptions := make([]string, len(products))
	productMap := make(map[string]uuid.UUID)
	for i, product := range products {
		productOptions[i] = product.Name
		productMap[product.Name] = product.ID
	}

	form.AddInputField("Instance Name", "", 50, nil, func(text string) {
		nameField = text
	})

	form.AddDropDown("Product", productOptions, 0, func(option string, optionIndex int) {
		selectedProductID = productMap[option]
	})

	// Set initial selected product if we have products
	if len(products) > 0 {
		selectedProductID = products[0].ID
	}

	form.AddButton("Create & Add", func() {
		if nameField == "" {
			// TODO: Show validation error in the future
			return
		}

		// Create the instance
		instance, err := service.CreateInstance(nameField, selectedProductID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		// Associate instance with domain
		err = service.AddInstanceToDomain(domainID, instance.ID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}
