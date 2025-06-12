package tviewapp

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

func NewDomainsView(contentContainer *ContentContainer) tview.Primitive {
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

	for i, d := range domains {
		instanceCount := len(d.Instances)
		instanceText := fmt.Sprintf("%d instances", instanceCount)

		table.SetCell(i+1, 0, tview.NewTableCell(d.ID.String()))
		table.SetCell(i+1, 1, tview.NewTableCell(d.Name))
		table.SetCell(i+1, 2, tview.NewTableCell(d.Description))
		table.SetCell(i+1, 3, tview.NewTableCell(instanceText))
	}

	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(domains) {
			domain := domains[row-1]
			contentContainer.SetContent(NewDomainDetailView(domain, contentContainer))
		}
	})

	table.SetSelectable(true, false)

	return table
}

func NewDomainDetailView(domain models.Domain, contentContainer *ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle(fmt.Sprintf("Domain: %s", domain.Name)).SetBorder(true)

	info := tview.NewTextView()
	info.SetTitle("Domain Information").SetBorder(true)
	info.SetText(fmt.Sprintf("ID: %s\nName: %s\nDescription: %s",
		domain.ID.String(), domain.Name, domain.Description))

	instancesTable := tview.NewTable().SetBorders(true)
	instancesTable.SetTitle("Instances in Domain").SetBorder(true)

	instancesTable.SetCell(0, 0, tview.NewTableCell("[::b]Name").SetSelectable(false))
	instancesTable.SetCell(0, 1, tview.NewTableCell("[::b]Product").SetSelectable(false))

	for i, instance := range domain.Instances {
		productName := ""
		if instance.Product.Name != "" {
			productName = instance.Product.Name
		}
		instancesTable.SetCell(i+1, 0, tview.NewTableCell(instance.Name))
		instancesTable.SetCell(i+1, 1, tview.NewTableCell(productName))
	}

	instancesTable.SetSelectable(true, false)
	instancesTable.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(domain.Instances) {
			instance := domain.Instances[row-1]
			contentContainer.SetContent(NewInstanceDetailScreen(instance.ID))
		}
	})

	backButton := tview.NewButton("Back to Domains").
		SetSelectedFunc(func() {
			contentContainer.SetContent(NewDomainsView(contentContainer))
		})

	flex.AddItem(info, 0, 1, false)
	flex.AddItem(instancesTable, 0, 2, true)
	flex.AddItem(backButton, 1, 0, false)

	return flex
}
