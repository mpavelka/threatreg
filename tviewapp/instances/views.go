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

func NewInstanceDetailView(instanceID uuid.UUID, contentContainer ContentContainer) tview.Primitive {
	instance, err := service.GetInstance(instanceID)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading instance: %v", err))
	}

	form := tview.NewForm().SetHorizontal(false)
	form.SetBorder(true).SetTitle("Edit Instance")
	form.AddInputField("Name", instance.Name, 30, nil, nil)
	form.AddInputField("Product", instance.Product.Name, 30, nil, nil)
	form.AddButton("Save", func() {
		// TODO: Save logic
	})
	form.AddButton("Back", func() {
		contentContainer.PopContent()
	})

	threatsTable := tview.NewTable().SetBorders(true)
	threatsTable.SetTitle("Related Threats").SetBorder(true)
	threatsTable.SetCell(0, 0, tview.NewTableCell("[::b]ID").SetSelectable(false))
	threatsTable.SetCell(0, 1, tview.NewTableCell("[::b]Name").SetSelectable(false))
	// Dummy data for now
	for i := 1; i <= 3; i++ {
		threatsTable.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("T%d", i)))
		threatsTable.SetCell(i, 1, tview.NewTableCell(fmt.Sprintf("Threat %d", i)))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(form, 0, 1, true).
		AddItem(threatsTable, 0, 2, false)
	return flex
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
			return NewInstanceDetailView(instancesList[row-1].ID, contentContainer)
		})
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
