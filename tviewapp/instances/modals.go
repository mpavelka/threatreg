package instances

import (
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

func createInstanceSelectThreatModal(instanceID uuid.UUID, onClose func()) tview.Primitive {
	selectExistingForm := createInstanceSelectExistingThreatForm(instanceID, onClose)
	createNewForm := createNewThreatForInstanceForm(instanceID, onClose)
	return createThreatSelectionModal("Select Threat", selectExistingForm, createNewForm)
}

func createProductSelectThreatModal(productID uuid.UUID, onClose func()) tview.Primitive {
	selectExistingForm := createProductSelectExistingThreatForm(productID, onClose)
	createNewForm := createNewThreatForProductForm(productID, onClose)
	return createThreatSelectionModal("Select Threat for Product", selectExistingForm, createNewForm)
}

func CreateEditInstanceModal(
	name string,
	onSave func(
		name string,
	),
	onClose func(),
) tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Edit Instance")

	form.AddInputField("Name", name, 50, nil, func(text string) {
		name = text
	})
	form.AddButton("Save", func() {
		onSave(name)
	})

	form.AddButton("Close", func() {
		onClose()
	})

	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)
	centerFlex.AddItem(form, 80, 0, true)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)

	modalContainer.AddItem(centerFlex, 15, 0, true)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	return modalContainer
}
