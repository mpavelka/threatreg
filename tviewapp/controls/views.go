package controls

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

type ContentContainer interface {
	PushContent(content tview.Primitive)
	PopContent() bool
	PushContentWithFactory(factory func() tview.Primitive)
}

func NewControlsView(contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Controls")

	newControlButton := tview.NewButton("New Control").
		SetSelectedFunc(func() {
			contentContainer.PushContent(createEditControlModal(
				"", "",
				func(title, description string) {
					if _, err := service.CreateControl(title, description); err == nil {
						contentContainer.PopContent()
					}
				},
				func() { contentContainer.PopContent() },
			))
		})

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)
	actionBar.AddItem(newControlButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false)

	controls, err := service.ListControls()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading controls: %v", err))
	}

	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Controls").SetBorder(true)
	table.SetFixed(1, 0)

	table.SetCell(0, 0, tview.NewTableCell("[::b]Title").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Description").SetSelectable(false))

	for i, c := range controls {
		table.SetCell(i+1, 0, tview.NewTableCell(c.Title))
		table.SetCell(i+1, 1, tview.NewTableCell(c.Description))
	}

	table.SetSelectable(true, false)

	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(table, 0, 1, true)

	return flex
}