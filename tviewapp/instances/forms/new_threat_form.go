package forms

import (
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// createNewThreatForm creates a form for creating and assigning a new threat
func createNewThreatForm(onCreateAndAssign func(title, description string), onClose func()) tview.Primitive {
	form := tview.NewForm()

	titleField := ""
	descriptionField := ""

	form.AddInputField("Title", titleField, 50, nil, func(text string) {
		titleField = text
	})

	form.AddInputField("Description", descriptionField, 50, nil, func(text string) {
		descriptionField = text
	})

	form.AddButton("Create & Assign", func() {
		if titleField == "" {
			// TODO: Show validation error in the future
			return
		}

		onCreateAndAssign(titleField, descriptionField)
		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}

// createNewThreatForInstanceForm creates a form for creating and assigning a new threat to an instance
func CreateNewThreatForInstanceForm(instanceID uuid.UUID, onClose func()) tview.Primitive {
	return createNewThreatForm(func(title, description string) {
		// Create the threat
		threat, err := service.CreateThreat(title, description)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		// Assign threat to instance
		_, err = service.AssignThreatToInstance(instanceID, threat.ID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}
	}, onClose)
}

// createNewThreatForProductForm creates a form for creating and assigning a new threat to a product
func CreateNewThreatForProductForm(productID uuid.UUID, onClose func()) tview.Primitive {
	return createNewThreatForm(func(title, description string) {
		// Create the threat
		threat, err := service.CreateThreat(title, description)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		// Assign threat to product
		_, err = service.AssignThreatToProduct(productID, threat.ID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}
	}, onClose)
}
