package domains

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

func NewDomainsView(contentContainer ContentContainer, instanceDetailFunc InstanceDetailScreenFunc) tview.Primitive {
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
	}

	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(domains) {
			domain := domains[row-1]
			contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailFunc))
		}
	})

	table.SetSelectable(true, false)

	return table
}

func NewDomainDetailView(domain models.Domain, contentContainer ContentContainer, instanceDetailFunc InstanceDetailScreenFunc) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle(fmt.Sprintf("Domain: %s", domain.Name)).SetBorder(true)

	info := tview.NewTextView()
	info.SetTitle("Domain Information").SetBorder(true)
	info.SetText(fmt.Sprintf("ID: %s\nName: %s\nDescription: %s",
		domain.ID.String(), domain.Name, domain.Description))

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)

	editButton := tview.NewButton("Edit").
		SetSelectedFunc(func() {
			modal := createEditDomainModal(domain, contentContainer, instanceDetailFunc, func(updatedDomain models.Domain) {
				contentContainer.SetContent(NewDomainDetailView(updatedDomain, contentContainer, instanceDetailFunc))
			})
			contentContainer.SetContent(modal)
		})

	addInstanceButton := tview.NewButton("Add Instance").
		SetSelectedFunc(func() {
			// TODO: Implement add instance functionality
		})

	actionBar.AddItem(editButton, 0, 1, false)
	actionBar.AddItem(addInstanceButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false) // Spacer

	instancesTable := tview.NewTable().SetBorders(true)
	instancesTable.SetTitle("Instances in Domain").SetBorder(true)

	instancesTable.SetCell(0, 0, tview.NewTableCell("[::b]Name").SetSelectable(false))
	instancesTable.SetCell(0, 1, tview.NewTableCell("[::b]Product").SetSelectable(false))

	instances, err := service.GetInstancesByDomainId(domain.ID)
	if err != nil {
		// If we can't load instances, show an error message in the table
		instancesTable.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Error loading instances: %v", err)))
		instancesTable.SetCell(1, 1, tview.NewTableCell(""))
		instances = []models.Instance{} // Set to empty slice to avoid nil access
	} else {
		for i, instance := range instances {
			productName := ""
			if instance.Product.Name != "" {
				productName = instance.Product.Name
			}
			instancesTable.SetCell(i+1, 0, tview.NewTableCell(instance.Name))
			instancesTable.SetCell(i+1, 1, tview.NewTableCell(productName))
		}
	}

	instancesTable.SetSelectable(true, false)
	instancesTable.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(instances) {
			instance := instances[row-1]
			contentContainer.SetContent(instanceDetailFunc(instance.ID))
		}
	})

	backButton := tview.NewButton("Back to Domains").
		SetSelectedFunc(func() {
			contentContainer.SetContent(NewDomainsView(contentContainer, instanceDetailFunc))
		})

	flex.AddItem(info, 5, 1, false)
	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(instancesTable, 0, 2, true)
	flex.AddItem(backButton, 1, 0, false)

	return flex
}

func createEditDomainModal(domain models.Domain, contentContainer ContentContainer, instanceDetailFunc InstanceDetailScreenFunc, onSave func(models.Domain)) tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Edit Domain")

	nameField := domain.Name
	descField := domain.Description

	form.AddInputField("Name", domain.Name, 50, nil, func(text string) {
		nameField = text
	})
	form.AddInputField("Description", domain.Description, 50, nil, func(text string) {
		descField = text
	})

	form.AddButton("Save", func() {
		updatedDomain, err := service.UpdateDomain(domain.ID, &nameField, &descField)
		if err != nil {
			// TODO: Show error message in the future
			contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailFunc))
			return
		}

		onSave(*updatedDomain)
	})

	form.AddButton("Close", func() {
		contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailFunc))
	})

	// Create a centered modal-like container
	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Top spacer

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Left spacer
	centerFlex.AddItem(form, 80, 0, true)           // Form with fixed width
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Right spacer

	modalContainer.AddItem(centerFlex, 15, 0, true)     // Form area
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Bottom spacer

	return modalContainer
}
