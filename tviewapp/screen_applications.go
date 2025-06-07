package tviewapp

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

// NewApplicationsView creates a tview.Primitive that lists all applications
func NewApplicationsView() tview.Primitive {
	applications, err := service.ListApplications()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading applications: %v", err))
	}
	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Applications").SetBorder(true)
	table.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Name").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Product").SetSelectable(false))
	for i, a := range applications {
		table.SetCell(i+1, 0, tview.NewTableCell(a.ID.String()))
		table.SetCell(i+1, 1, tview.NewTableCell(a.Name))
		table.SetCell(i+1, 2, tview.NewTableCell(a.Product.Name))
	}
	return table
}
