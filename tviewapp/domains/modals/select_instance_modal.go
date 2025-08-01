package modals

import (
	"threatreg/tviewapp/domains/forms"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func CreateSelectInstanceModal(domainID uuid.UUID, onClose func()) tview.Primitive {
	// Create tabs
	tabs := tview.NewPages()
	tabs.SetBorder(true).SetTitle("Select Instance")

	// Create Select Existing tab
	selectExistingForm := forms.CreateSelectExistingInstanceForm(domainID, onClose)
	tabs.AddPage("Select Existing", selectExistingForm, true, true)

	// Create New tab
	createNewForm := forms.CreateNewInstanceForm(domainID, onClose)
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
