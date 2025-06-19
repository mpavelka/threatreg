package products

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

func NewProductsView(contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Products")

	newProductButton := tview.NewButton("New Product").
		SetSelectedFunc(func() {
			contentContainer.PushContent(createEditProductModal(
				"", "",
				func(name, description string) {
					if _, err := service.CreateProduct(name, description); err == nil {
						contentContainer.PopContent()
					}
				},
				func() { contentContainer.PopContent() },
			))
		})

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)
	actionBar.AddItem(newProductButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false)

	products, err := service.ListProducts()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading products: %v", err))
	}

	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Products").SetBorder(true)
	table.SetFixed(1, 0)

	table.SetCell(0, 0, tview.NewTableCell("[::b]Name").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Description").SetSelectable(false))

	for i, p := range products {
		table.SetCell(i+1, 0, tview.NewTableCell(p.Name))
		table.SetCell(i+1, 1, tview.NewTableCell(p.Description))
	}

	table.SetSelectable(true, false)

	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(table, 0, 1, true)

	return flex
}
