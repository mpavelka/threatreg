package domains

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func NewDomainsView(contentContainer ContentContainer, instanceDetailScreenBuilder InstanceDetailScreenBuilder) tview.Primitive {
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
			contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailScreenBuilder))
		}
	})

	table.SetSelectable(true, false)

	return table
}

func NewDomainDetailView(domain models.Domain, contentContainer ContentContainer, instanceDetailScreenBuilder InstanceDetailScreenBuilder) tview.Primitive {
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
			modal := createEditDomainModal(domain, contentContainer, instanceDetailScreenBuilder, func(updatedDomain models.Domain) {
				contentContainer.SetContent(NewDomainDetailView(updatedDomain, contentContainer, instanceDetailScreenBuilder))
			})
			contentContainer.SetContent(modal)
		})

	addInstanceButton := tview.NewButton("Add Instance").
		SetSelectedFunc(func() {
			modal := createSelectInstanceModal(domain.ID, contentContainer, instanceDetailScreenBuilder, func() {
				contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailScreenBuilder))
			})
			contentContainer.SetContent(modal)
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
			contentContainer.SetContent(instanceDetailScreenBuilder(instance.ID))
		}
	})

	backButton := tview.NewButton("Back to Domains").
		SetSelectedFunc(func() {
			contentContainer.SetContent(NewDomainsView(contentContainer, instanceDetailScreenBuilder))
		})

	flex.AddItem(info, 5, 1, false)
	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(instancesTable, 0, 2, true)
	flex.AddItem(backButton, 1, 0, false)

	return flex
}

func createEditDomainModal(domain models.Domain, contentContainer ContentContainer, instanceDetailScreenBuilder InstanceDetailScreenBuilder, onSave func(models.Domain)) tview.Primitive {
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
			contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailScreenBuilder))
			return
		}

		onSave(*updatedDomain)
	})

	form.AddButton("Close", func() {
		contentContainer.SetContent(NewDomainDetailView(domain, contentContainer, instanceDetailScreenBuilder))
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

func createSelectInstanceModal(domainID uuid.UUID, contentContainer ContentContainer, instanceDetailScreenBuilder InstanceDetailScreenBuilder, onClose func()) tview.Primitive {
	// Create tabs
	tabs := tview.NewPages()
	tabs.SetBorder(true).SetTitle("Select Instance")

	// Create Select Existing tab
	selectExistingForm := createSelectExistingInstanceForm(domainID, onClose)
	tabs.AddPage("Select Existing", selectExistingForm, true, true)

	// Create New tab
	createNewForm := createNewInstanceForm(domainID, onClose)
	tabs.AddPage("Create New", createNewForm, true, false)

	// Create tab navigation
	tabNavigation := tview.NewFlex()
	selectExistingButton := tview.NewButton("Select Existing").
		SetSelectedFunc(func() {
			tabs.SwitchToPage("Select Existing")
		})
	createNewButton := tview.NewButton("Create New").
		SetSelectedFunc(func() {
			tabs.SwitchToPage("Create New")
		})

	tabNavigation.AddItem(selectExistingButton, 0, 1, true)
	tabNavigation.AddItem(createNewButton, 0, 1, false)

	// Combine tabs with navigation
	tabContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	tabContainer.AddItem(tabNavigation, 3, 0, false)
	tabContainer.AddItem(tabs, 0, 1, true)

	// Create a centered modal-like container
	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Top spacer

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Left spacer
	centerFlex.AddItem(tabContainer, 80, 0, true)   // Tab container with fixed width
	centerFlex.AddItem(tview.NewBox(), 0, 1, false) // Right spacer

	modalContainer.AddItem(centerFlex, 18, 0, true)     // Form area (increased height for tabs)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false) // Bottom spacer

	return modalContainer
}

func createSelectExistingInstanceForm(domainID uuid.UUID, onClose func()) tview.Primitive {
	form := tview.NewForm()

	// Get list of all instances
	instances, err := service.ListInstances()
	if err != nil {
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading instances: %v", err))
		return errorView
	}

	if len(instances) == 0 {
		noInstancesView := tview.NewTextView().SetText("No instances available to add to this domain.")
		return noInstancesView
	}

	// Create instance options for dropdown with fuzzy search
	instanceOptions := make([]string, len(instances))
	instanceMap := make(map[string]uuid.UUID)
	for i, instance := range instances {
		displayName := fmt.Sprintf("%s (%s)", instance.Name, instance.Product.Name)
		instanceOptions[i] = displayName
		instanceMap[displayName] = instance.ID
	}

	var selectedInstanceID uuid.UUID
	form.AddDropDown("Instance", instanceOptions, 0, func(option string, optionIndex int) {
		selectedInstanceID = instanceMap[option]
	})

	// Set initial selected instance
	if len(instances) > 0 {
		selectedInstanceID = instances[0].ID
	}

	form.AddButton("Add to Domain", func() {
		// Associate instance with domain
		err := service.AddInstanceToDomain(domainID, selectedInstanceID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}

func createNewInstanceForm(domainID uuid.UUID, onClose func()) tview.Primitive {
	form := tview.NewForm()

	nameField := ""
	var selectedProductID uuid.UUID

	// Get list of products for dropdown
	products, err := service.ListProducts()
	if err != nil {
		// If we can't load products, show error and return
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading products: %v", err))
		return errorView
	}

	// Create product options for dropdown
	productOptions := make([]string, len(products))
	productMap := make(map[string]uuid.UUID)
	for i, product := range products {
		productOptions[i] = product.Name
		productMap[product.Name] = product.ID
	}

	form.AddInputField("Instance Name", "", 50, nil, func(text string) {
		nameField = text
	})

	form.AddDropDown("Product", productOptions, 0, func(option string, optionIndex int) {
		selectedProductID = productMap[option]
	})

	// Set initial selected product if we have products
	if len(products) > 0 {
		selectedProductID = products[0].ID
	}

	form.AddButton("Create & Add", func() {
		if nameField == "" {
			// TODO: Show validation error in the future
			return
		}

		// Create the instance
		instance, err := service.CreateInstance(nameField, selectedProductID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		// Associate instance with domain
		err = service.AddInstanceToDomain(domainID, instance.ID)
		if err != nil {
			// TODO: Show error message in the future
			return
		}

		onClose()
	})

	form.AddButton("Cancel", func() {
		onClose()
	})

	return form
}
