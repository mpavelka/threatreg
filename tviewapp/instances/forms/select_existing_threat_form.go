package forms

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

func CreateInstanceSelectExistingThreatForm(instanceID uuid.UUID, onClose func()) tview.Primitive {
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
func CreateProductSelectExistingThreatForm(productID uuid.UUID, onClose func()) tview.Primitive {
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
