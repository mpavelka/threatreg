package common

import (
	"github.com/rivo/tview"
)

func CreateConfirmationModal(title, message string, onYes, onNo func()) tview.Primitive {
	modal := tview.NewModal()
	modal.SetText(message)
	modal.AddButtons([]string{"Yes", "No"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			onYes()
		} else {
			onNo()
		}
	})

	return modal
}
