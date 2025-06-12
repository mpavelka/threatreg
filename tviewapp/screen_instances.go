package tviewapp

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

var (
	filterForm     *tview.Form
	instancesTable *tview.Table
	instancesList  []models.Instance
)

// NewInstancesView creates a tview.Primitive that lists all instances
func NewInstancesView(contentContainer *ContentContainer) tview.Primitive {

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

func initInstancesTable(contentContainer *ContentContainer) {
	instancesTable = tview.NewTable().SetBorders(true)
	instancesTable.SetTitle("Instances").SetTitleAlign(tview.AlignLeft).SetBorder(true)

	// Header (use color for bold effect)
	instancesTable.SetFixed(1, 0) // Keep header fixed
	instancesTable.SetSelectable(true, false)
	instancesTable.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return // header
		}
		contentContainer.SetContent(NewInstanceDetailScreen(instancesList[row-1].ID))
	})
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
