package tviewapp

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// NewInstanceDetailScreen returns a Flex with a form to edit the instance and a dummy table for related threats
func NewInstanceDetailScreen(instanceID uuid.UUID) tview.Primitive {
	instance, err := service.GetInstance(instanceID)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading instance: %v", err))
	}

	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Edit Instance")
	form.AddInputField("Name", instance.Name, 30, nil, nil)
	form.AddInputField("Product", instance.Product.Name, 30, nil, nil)
	form.AddButton("Save", func() {
		// TODO: Save logic
	})
	form.AddButton("Cancel", func() {
		// TODO: Cancel logic
	})

	threatsTable := tview.NewTable().SetBorders(true)
	threatsTable.SetTitle("Related Threats").SetBorder(true)
	threatsTable.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	threatsTable.SetCell(0, 1, tview.NewTableCell("[::b]Name").SetSelectable(false))
	// Dummy data for now
	for i := 1; i <= 3; i++ {
		threatsTable.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("T%d", i)))
		threatsTable.SetCell(i, 1, tview.NewTableCell(fmt.Sprintf("Threat %d", i)))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(form, 0, 1, true).
		AddItem(threatsTable, 0, 2, false)
	return flex
}
