package instances

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

var (
	filterForm     *tview.Form
	instancesTable *tview.Table
	instancesList  []models.Instance
)

// NewInstancesView creates a tview.Primitive that lists all instances
func NewInstancesView(contentContainer ContentContainer) tview.Primitive {
	initFilter()
	reloadInstances()
	initInstancesTable(contentContainer)
	updateInstancesTable()

	// Create layout
	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(30, 0).
		AddItem(filterForm, 0, 0, 1, 1, 0, 0, true).
		AddItem(instancesTable, 0, 1, 1, 2, 0, 0, true)
	return grid
}

func reloadInstances() {
	instances, err := service.FilterInstances(
		filterForm.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		filterForm.GetFormItemByLabel("Product").(*tview.InputField).GetText(),
	)
	if err != nil {
		instancesTable.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Error loading instances: %v", err)))
		return
	}
	instancesList = instances
}

func initFilter() {
	filterForm = tview.NewForm().SetHorizontal(false)
	filterForm.SetBorder(true).SetTitle("Filter").SetTitleAlign(tview.AlignLeft)
	filterForm.AddInputField("Name", "", 0, nil, nil)
	filterForm.AddInputField("Product", "", 0, nil, nil)
	filterForm.AddButton("Filter", func() {
		// Apply filters
		reloadInstances()
		updateInstancesTable()
	})
	filterForm.AddButton("Reset", func() {
		// Reset filters
		filterForm.GetFormItemByLabel("Name").(*tview.InputField).SetText("")
		filterForm.GetFormItemByLabel("Product").(*tview.InputField).SetText("")
		reloadInstances()
		updateInstancesTable()
	})
}

func initInstancesTable(contentContainer ContentContainer) {
	instancesTable = tview.NewTable().SetBorders(true)
	instancesTable.SetTitle("Instances").SetTitleAlign(tview.AlignLeft).SetBorder(true)

	// Header (use color for bold effect)
	instancesTable.SetFixed(1, 0) // Keep header fixed
	instancesTable.SetSelectable(true, false)
	instancesTable.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return // header
		}
		contentContainer.PushContentWithFactory(func() tview.Primitive {
			return NewInstanceThreatManager(instancesList[row-1].ID, contentContainer)
		})
	})
}

// NewInstanceThreatManager creates a threat management view for an instance
func NewInstanceThreatManager(instanceID uuid.UUID, contentContainer ContentContainer) tview.Primitive {
	instance, err := service.GetInstance(instanceID)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading instance: %v", err))
	}

	// Left column - Instance section
	instanceText := tview.NewTextView()
	instanceText.SetBorder(true).SetTitle("Instance Information")
	instanceText.SetText(fmt.Sprintf("Name: %s\nID: %s", instance.Name, instance.ID.String()))

	instanceButton := tview.NewButton("Edit Instance").SetSelectedFunc(func() {
		// TODO: Navigate to instance edit view
	})

	instanceSection := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(instanceText, 0, 1, false).
		AddItem(instanceButton, 1, 0, true)

	// Left column - Product section
	productText := tview.NewTextView()
	productText.SetBorder(true).SetTitle("Product Information")
	productText.SetText(fmt.Sprintf("Name: %s\nID: %s", instance.Product.Name, instance.Product.ID.String()))

	productButton := tview.NewButton("Edit Product").SetSelectedFunc(func() {
		// TODO: Navigate to product edit view
	})

	productSection := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(productText, 0, 1, false).
		AddItem(productButton, 1, 0, false)

	// Left column container
	leftColumn := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(instanceSection, 0, 1, true).
		AddItem(productSection, 0, 1, false)

	// Right column - Actions section
	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)

	actionBar.AddItem(tview.NewButton("Add Instance Threat").SetSelectedFunc(func() {
		contentContainer.PushContent(createInstanceSelectThreatModal(instance.ID, func() {
			contentContainer.PopContent()
		}))
	}), 0, 1, false)
	actionBar.AddItem(tview.NewButton("Add Product Threat").SetSelectedFunc(func() {
		contentContainer.PushContent(createProductSelectThreatModal(instance.Product.ID, func() {
			contentContainer.PopContent()
		}))
	}), 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false) // Spacer

	// Right column - Threat Assignments table
	threatTable := tview.NewTable().SetBorders(true)
	threatTable.SetTitle("Threat Assignments").SetBorder(true)
	threatTable.SetCell(0, 0, tview.NewTableCell("[::b]Type").SetSelectable(false))
	threatTable.SetCell(0, 1, tview.NewTableCell("[::b]Threat").SetSelectable(false))
	threatTable.SetCell(0, 2, tview.NewTableCell("[::b]Description").SetSelectable(false))

	// Load actual threat assignments for this instance
	instanceAssignments, err := service.ListThreatAssignmentsByInstanceID(instance.ID)
	if err != nil {
		// Show error in table
		threatTable.SetCell(1, 0, tview.NewTableCell("Error"))
		threatTable.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("Failed to load: %v", err)))
		threatTable.SetCell(1, 2, tview.NewTableCell(""))
	} else {
		// Load product-level assignments as well (instance inherits product threats)
		productAssignments, err := service.ListThreatAssignmentsByProductID(instance.Product.ID)
		if err != nil {
			productAssignments = []models.ThreatAssignment{}
		}

		row := 1

		// Add instance-specific threats
		for _, assignment := range instanceAssignments {
			threatTable.SetCell(row, 0, tview.NewTableCell("Instance"))
			threatTable.SetCell(row, 1, tview.NewTableCell(assignment.Threat.Title))
			threatTable.SetCell(row, 2, tview.NewTableCell(assignment.Threat.Description))
			row++
		}

		// Add product-level threats (inherited)
		for _, assignment := range productAssignments {
			threatTable.SetCell(row, 0, tview.NewTableCell("Product"))
			threatTable.SetCell(row, 1, tview.NewTableCell(assignment.Threat.Title))
			threatTable.SetCell(row, 2, tview.NewTableCell(assignment.Threat.Description))
			row++
		}

		// Show message if no threats are assigned
		if len(instanceAssignments) == 0 && len(productAssignments) == 0 {
			threatTable.SetCell(1, 0, tview.NewTableCell(""))
			threatTable.SetCell(1, 1, tview.NewTableCell("No threats assigned"))
			threatTable.SetCell(1, 2, tview.NewTableCell(""))
		}
	}

	// Right column container
	rightColumn := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(actionBar, 3, 0, false).
		AddItem(threatTable, 0, 1, false)

	// Main layout - two columns
	mainLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftColumn, 30, 0, true).
		AddItem(rightColumn, 0, 1, false)

	// Add navigation back button
	wrapper := tview.NewFlex().SetDirection(tview.FlexRow)

	backButton := tview.NewButton("← Back").SetSelectedFunc(func() {
		contentContainer.PopContent()
	})
	backButton.SetBorder(true)

	wrapper.AddItem(backButton, 3, 0, false).
		AddItem(mainLayout, 0, 1, true)

	return mainLayout
}

func updateInstancesTable() {
	instancesTable.Clear()

	// Header
	instancesTable.SetCell(0, 0, tview.NewTableCell("[::b]Name"))
	instancesTable.SetCell(0, 1, tview.NewTableCell("[::b]Product"))

	// Data
	for i, instance := range instancesList {
		instancesTable.SetCell(i+1, 0, tview.NewTableCell(instance.Name))
		instancesTable.SetCell(i+1, 1, tview.NewTableCell(instance.Product.Name))
	}
}
