package modals

import (
	"github.com/rivo/tview"
)

func CreateEditInstanceModal(
	name string,
	onSave func(name string),
	onClose func(),
) tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Edit Instance")

	form.AddInputField("Name", name, 50, nil, func(text string) {
		name = text
	})

	form.AddButton("Save", func() {
		onSave(name)
	})

	form.AddButton("Close", func() {
		onClose()
	})

	// Create a centered modal-like container
	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Top spacer

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Left spacer
	centerFlex.AddItem(form, 80, 0, true)           // Form with fixed width
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Right spacer

	modalContainer.AddItem(centerFlex, 12, 0, true)     // Form area
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Bottom spacer

	return modalContainer
}
