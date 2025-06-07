package tviewapp

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

// NewProductsView creates a tview.Primitive that lists all products
func NewProductsView() tview.Primitive {
	products, err := service.ListProducts()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading products: %v", err))
	}

	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Products").SetBorder(true)

	// Header (use color for bold effect)
	table.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Name").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Description").SetSelectable(false))

	// Data
	for i, p := range products {
		desc := p.Description
		table.SetCell(i+1, 0, tview.NewTableCell(p.ID.String()))
		table.SetCell(i+1, 1, tview.NewTableCell(p.Name))
		table.SetCell(i+1, 2, tview.NewTableCell(desc))
	}

	return table
}
