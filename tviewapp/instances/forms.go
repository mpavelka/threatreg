package instances

import (
	"threatreg/internal/models"

	"github.com/rivo/tview"
)

func createInstanceEditForm(instance models.Instance, contentContainer ContentContainer) *tview.Form {
	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Edit Instance")

	form.AddInputField("Name", instance.Name, 30, nil, nil)
	form.AddInputField("Product", instance.Product.Name, 30, nil, nil) // Read-only for now

	form.AddButton("Save", func() {
		// TODO: Implement save logic using service.UpdateInstance
		// For now, just go back
		contentContainer.PopContent()
	})

	form.AddButton("Cancel", func() {
		contentContainer.PopContent()
	})

	return form
}

func createInstanceFilterForm() *tview.Form {
	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Filter").SetTitleAlign(tview.AlignLeft)
	form.AddInputField("Name", "", 0, nil, nil)
	form.AddInputField("Product", "", 0, nil, nil)
	return form
}
