package tviewapp

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// NewApplicationDetailScreen returns a Flex with a form to edit the application and a dummy table for related threats
func NewApplicationDetailScreen(appID uuid.UUID) tview.Primitive {
	app, err := service.GetApplication(appID)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading application: %v", err))
	}

	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Edit Application")
	form.AddInputField("Name", app.Name, 30, nil, nil)
	form.AddInputField("Product", app.Product.Name, 30, nil, nil)
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
