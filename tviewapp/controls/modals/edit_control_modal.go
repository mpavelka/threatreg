package modals

import (
	"github.com/rivo/tview"
)

func CreateEditControlModal(
	title string,
	description string,
	onSave func(
		title string,
		description string,
	),
	onClose func(),
) tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Edit Control")

	form.AddInputField("Title", title, 50, nil, func(text string) {
		title = text
	})
	form.AddInputField("Description", description, 50, nil, func(text string) {
		description = text
	})

	form.AddButton("Save", func() {
		onSave(title, description)
	})

	form.AddButton("Close", func() {
		onClose()
	})

	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)
	centerFlex.AddItem(form, 80, 0, true)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)

	modalContainer.AddItem(centerFlex, 15, 0, true)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	return modalContainer
}
