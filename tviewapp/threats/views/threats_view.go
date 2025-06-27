package views

import (
	"fmt"
	"threatreg/internal/service"
	"threatreg/tviewapp/common"
	"threatreg/tviewapp/threats/modals"

	"github.com/rivo/tview"
)

type ContentContainer interface {
	PushContent(content tview.Primitive)
	PopContent() bool
	PushContentWithFactory(factory func() tview.Primitive)
}

func NewThreatsView(contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Threats")

	newThreatButton := tview.NewButton("New Threat").
		SetSelectedFunc(func() {
			contentContainer.PushContent(modals.CreateEditThreatModal(
				"", "",
				func(title, description string) {
					if _, err := service.CreateThreat(title, description); err == nil {
						contentContainer.PopContent()
					}
				},
				func() { contentContainer.PopContent() },
			))
		})

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)
	actionBar.AddItem(newThreatButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false)

	threats, err := service.ListThreats()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading threats: %v", err))
	}

	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Threats").SetBorder(true)
	table.SetFixed(1, 0)

	table.SetCell(0, 0, tview.NewTableCell("[::b]Title").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Description").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Actions").SetSelectable(false))

	for i, t := range threats {
		table.SetCell(i+1, 0, tview.NewTableCell(t.Title))
		table.SetCell(i+1, 1, tview.NewTableCell(t.Description))
		removeButton := "[red]Remove[-]"
		table.SetCell(i+1, 2, tview.NewTableCell(removeButton).SetSelectable(true))
	}

	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(threats) && column == 2 {
			threat := threats[row-1]
			contentContainer.PushContent(common.CreateConfirmationModal(
				"Remove Threat",
				fmt.Sprintf("Are you sure you want to remove threat '%s'?", threat.Title),
				func() {
					err := service.DeleteThreat(threat.ID)
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
