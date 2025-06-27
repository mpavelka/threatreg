package modals

import (
	instancesForms "threatreg/tviewapp/instances/forms"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// createThreatSelectionModal creates a modal with tabbed interface for threat selection
func createThreatSelectionModal(title string, selectExistingForm, createNewForm tview.Primitive) tview.Primitive {
	// Create tabs
	tabs := tview.NewPages()
	tabs.SetBorder(true).SetTitle(title)

	// Add tabs
	tabs.AddPage("Select Existing", selectExistingForm, true, true)
	tabs.AddPage("Create New", createNewForm, true, false)

	// Create tab navigation
	tabNavigation := tview.NewFlex()
	selectExistingButton := tview.NewButton("Select Existing").
		SetSelectedFunc(func() {
			tabs.SwitchToPage("Select Existing")
		})
	createNewButton := tview.NewButton("Create New").
		SetSelectedFunc(func() {
			tabs.SwitchToPage("Create New")
		})

	tabNavigation.AddItem(selectExistingButton, 0, 1, true)
	tabNavigation.AddItem(createNewButton, 0, 1, false)

	// Combine tabs with navigation
	tabContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	tabContainer.AddItem(tabNavigation, 3, 0, false)
	tabContainer.AddItem(tabs, 0, 1, true)

	// Create a centered modal-like container
	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Top spacer

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Left spacer
	centerFlex.AddItem(tabContainer, 80, 0, true)   // Tab container with fixed width
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Right spacer

	modalContainer.AddItem(centerFlex, 18, 0, true)     // Form area (increased height for tabs)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Bottom spacer

	return modalContainer
}

func CreateInstanceSelectThreatModal(instanceID uuid.UUID, onClose func()) tview.Primitive {
	selectExistingForm := instancesForms.CreateInstanceSelectExistingThreatForm(instanceID, onClose)
	createNewForm := instancesForms.CreateNewThreatForInstanceForm(instanceID, onClose)
	return createThreatSelectionModal("Select Threat", selectExistingForm, createNewForm)
}

func CreateProductSelectThreatModal(productID uuid.UUID, onClose func()) tview.Primitive {
	selectExistingForm := instancesForms.CreateProductSelectExistingThreatForm(productID, onClose)
	createNewForm := instancesForms.CreateNewThreatForProductForm(productID, onClose)
	return createThreatSelectionModal("Select Threat for Product", selectExistingForm, createNewForm)
}
