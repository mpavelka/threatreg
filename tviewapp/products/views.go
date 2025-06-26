package products

import (
	"fmt"
	"threatreg/internal/service"
	"threatreg/tviewapp/common"

	"github.com/rivo/tview"
)

func NewProductsView(contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Products")

	newProductButton := tview.NewButton("New Product").
		SetSelectedFunc(func() {
			contentContainer.PushContent(CreateEditProductModal(
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
	table.SetCell(0, 2, tview.NewTableCell("[::b]Actions").SetSelectable(false))

	for i, p := range products {
		table.SetCell(i+1, 0, tview.NewTableCell(p.Name))
		table.SetCell(i+1, 1, tview.NewTableCell(p.Description))
		removeButton := "[red]Remove[-]"
		table.SetCell(i+1, 2, tview.NewTableCell(removeButton).SetSelectable(true))
	}

	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(products) && column == 2 {
			product := products[row-1]
			contentContainer.PushContent(common.CreateConfirmationModal(
				"Remove Product",
				fmt.Sprintf("Are you sure you want to remove product '%s'?", product.Name),
				func() {
					err := service.DeleteProduct(product.ID)
					if err != nil {
						return
					}
					contentContainer.PopContent()
				},
				func() {
					contentContainer.PopContent()
				},
			))
		}
	})

	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(table, 0, 1, true)

	return flex
}
