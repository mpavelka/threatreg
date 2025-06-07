package tviewapp

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

// NewThreatsView creates a tview.Primitive that lists all threats
func NewThreatsView() tview.Primitive {
	threats, err := service.ListThreats()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading threats: %v", err))
	}
	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Threats").SetBorder(true)
	table.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Title").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Description").SetSelectable(false))
	for i, t := range threats {
		table.SetCell(i+1, 0, tview.NewTableCell(t.ID.String()))
		table.SetCell(i+1, 1, tview.NewTableCell(t.Title))
		table.SetCell(i+1, 2, tview.NewTableCell(t.Description))
	}
	return table
}
