package instances

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func createInstanceEditForm(instance models.Instance, contentContainer ContentContainer) *tview.Form {
	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Edit Instance")

	form.AddInputField("Name", instance.Name, 30, nil, nil)
	form.AddInputField("Product", instance.Product.Name, 30, nil, nil) // Read-only for now

	form.AddButton("Save", func() {
		// TODO: Implement save logic using service.UpdateInstance
		// For now, just go back
		contentContainer.PopContent()
	})

	form.AddButton("Cancel", func() {
		contentContainer.PopContent()
	})

	return form
}

func createInstanceFilterForm() *tview.Form {
	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Filter").SetTitleAlign(tview.AlignLeft)
	form.AddInputField("Name", "", 0, nil, nil)
	form.AddInputField("Product", "", 0, nil, nil)
	return form
}

func createSelectExistingThreatForm(instanceID uuid.UUID, onClose func()) tview.Primitive {
	form := tview.NewForm()

	// Get list of all threats
	threats, err := service.ListThreats()
	if err != nil {
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading threats: %v", err))
		return errorView
	}

	if len(threats) == 0 {
		noThreatsView := tview.NewTextView().SetText("No threats available to assign to this instance.")
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

	form.AddButton("Assign to Instance", func() {
		// Assign threat to instance
		_, err := service.AssignThreatToInstance(instanceID, selectedThreatID)
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

func createNewThreatForm(instanceID uuid.UUID, onClose func()) tview.Primitive {
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

		// Create the threat
		threat, err := service.CreateThreat(titleField, descriptionField)
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

		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}
