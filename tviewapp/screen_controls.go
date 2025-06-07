package tviewapp

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

// NewControlsView creates a tview.Primitive that lists all controls
func NewControlsView() tview.Primitive {
	controls, err := service.ListControls()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading controls: %v", err))
	}
	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Controls").SetBorder(true)
	table.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Title").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Description").SetSelectable(false))
	for i, c := range controls {
		table.SetCell(i+1, 0, tview.NewTableCell(c.ID.String()))
		table.SetCell(i+1, 1, tview.NewTableCell(c.Title))
		table.SetCell(i+1, 2, tview.NewTableCell(c.Description))
	}
	return table
}
