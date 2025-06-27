package views

import (
	"fmt"
	"threatreg/internal/service"
	"threatreg/tviewapp/common"
	"threatreg/tviewapp/domains/modals"

	"github.com/rivo/tview"
)

type ContentContainer interface {
	PushContent(content tview.Primitive)
	PopContent() bool
	PushContentWithFactory(factory func() tview.Primitive)
}

func NewDomainsView(contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Domains")

	newDomainButton := tview.NewButton("New Domain").
		SetSelectedFunc(func() {
			contentContainer.PushContent(modals.CreateEditDomainModal(
				"", "",
				func(name, description string) {
					if _, err := service.CreateDomain(name, description); err == nil {
						contentContainer.PopContent()
					}
				},
				func() { contentContainer.PopContent() },
			))
		})

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)
	actionBar.AddItem(newDomainButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false) // Spacer

	domains, err := service.ListDomains()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading domains: %v", err))
	}

	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Domains").SetBorder(true)

	table.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Name").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Description").SetSelectable(false))
	table.SetCell(0, 3, tview.NewTableCell("[::b]Instances").SetSelectable(false))
	table.SetCell(0, 4, tview.NewTableCell("[::b]Actions").SetSelectable(false))

	for i, d := range domains {
		instances, err := service.GetInstancesByDomainId(d.ID)
		instanceCount := 0
		if err == nil {
			instanceCount = len(instances)
		}
		instanceText := fmt.Sprintf("%d instances", instanceCount)

		table.SetCell(i+1, 0, tview.NewTableCell(d.ID.String()))
		table.SetCell(i+1, 1, tview.NewTableCell(d.Name))
		table.SetCell(i+1, 2, tview.NewTableCell(d.Description))
		table.SetCell(i+1, 3, tview.NewTableCell(instanceText))
		removeButton := "[red]Remove[-]"
		table.SetCell(i+1, 4, tview.NewTableCell(removeButton).SetSelectable(true))
	}

	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(domains) {
			domain := domains[row-1]
			if column == 4 {
				// Remove button clicked
				contentContainer.PushContent(common.CreateConfirmationModal(
					"Remove Domain",
					fmt.Sprintf("Are you sure you want to remove domain '%s'?", domain.Name),
					func() {
						err := service.DeleteDomain(domain.ID)
						if err != nil {
							return
						}
						contentContainer.PopContent()
					},
					func() {
						contentContainer.PopContent()
					},
				))
			} else {
				// Navigate to domain detail
				contentContainer.PushContentWithFactory(func() tview.Primitive {
					return NewDomainDetailView(domain, contentContainer)
				})
			}
		}
	})

	table.SetSelectable(true, true)

	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(table, 0, 1, true)

	return flex
}
