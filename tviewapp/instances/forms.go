package instances

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// createSelectExistingThreatForm creates a form for selecting and assigning an existing threat
func createSelectExistingThreatForm(buttonText, noThreatsMessage string, onAssign func(threatID uuid.UUID), onClose func()) tview.Primitive {
	form := tview.NewForm()

	// Get list of all threats
	threats, err := service.ListThreats()
	if err != nil {
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading threats: %v", err))
		return errorView
	}

	if len(threats) == 0 {
		noThreatsView := tview.NewTextView().SetText(noThreatsMessage)
		return noThreatsView
	}

	// Create threat options for dropdown
	threatOptions := make([]string, len(threats))
	threatMap := make(map[string]uuid.UUID)
	for i, threat := range threats {
		displayName := fmt.Sprintf("%s - %s", threat.Title, threat.Description)
		if len(displayName) > 80 {
			displayName = displayName[:77] + "..."
		}
		threatOptions[i] = displayName
		threatMap[displayName] = threat.ID
	}

	var selectedThreatID uuid.UUID
	form.AddDropDown("Threat", threatOptions, 0, func(option string, optionIndex int) {
		selectedThreatID = threatMap[option]
	})

	// Set initial selected threat
	if len(threats) > 0 {
		selectedThreatID = threats[0].ID
	}

	form.AddButton(buttonText, func() {
		onAssign(selectedThreatID)
		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}

func createInstanceSelectExistingThreatForm(instanceID uuid.UUID, onClose func()) tview.Primitive {
	return createSelectExistingThreatForm(
		"Assign to Instance",
		"No threats available to assign to this instance.",
		func(threatID uuid.UUID) {
			_, err := service.AssignThreatToInstance(instanceID, threatID)
			if err != nil {
				// TODO: Show error message in the future
			}
		},
		onClose,
	)
}
func createProductSelectExistingThreatForm(productID uuid.UUID, onClose func()) tview.Primitive {
	return createSelectExistingThreatForm(
		"Assign to Product",
		"No threats available to assign to this product.",
		func(threatID uuid.UUID) {
			_, err := service.AssignThreatToProduct(productID, threatID)
			if err != nil {
				// TODO: Show error message in the future
			}
		},
		onClose,
	)
}

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
func createNewThreatForInstanceForm(instanceID uuid.UUID, onClose func()) tview.Primitive {
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
func createNewThreatForProductForm(productID uuid.UUID, onClose func()) tview.Primitive {
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
