package views

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"
	"threatreg/tviewapp/common"

	"github.com/rivo/tview"
)

type ContentContainer interface {
	PushContent(content tview.Primitive)
	PopContent() bool
	PushContentWithFactory(factory func() tview.Primitive)
}

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
	filterForm.SetBorder(true).SetTitle("Filter")

	filterForm.AddInputField("Name", "", 20, nil, func(text string) {
		reloadInstances()
		updateInstancesTable()
	})

	filterForm.AddInputField("Product", "", 20, nil, func(text string) {
		reloadInstances()
		updateInstancesTable()
	})
}

func initInstancesTable(contentContainer ContentContainer) {
	instancesTable = tview.NewTable().SetBorders(true)
	instancesTable.SetTitle("Instances").SetBorder(true)
	instancesTable.SetFixed(1, 0)

	instancesTable.SetCell(0, 0, tview.NewTableCell("[::b]Name").SetSelectable(false))
	instancesTable.SetCell(0, 1, tview.NewTableCell("[::b]Product").SetSelectable(false))
	instancesTable.SetCell(0, 2, tview.NewTableCell("[::b]Actions").SetSelectable(false))

	instancesTable.SetSelectable(true, true)
	instancesTable.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(instancesList) {
			if column == 2 {
				// Actions column clicked
				instance := instancesList[row-1]
				contentContainer.PushContent(common.CreateConfirmationModal(
					"Remove Instance",
					fmt.Sprintf("Are you sure you want to remove instance '%s'?", instance.Name),
					func() {
						err := service.DeleteInstance(instance.ID)
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
				// Navigate to instance detail for other columns
				contentContainer.PushContentWithFactory(func() tview.Primitive {
					return NewComponentThreatManager(instancesList[row-1].ID, contentContainer)
				})
			}
		}
	})
}

func updateInstancesTable() {
	// Clear existing rows (except header)
	for i := instancesTable.GetRowCount() - 1; i > 0; i-- {
		instancesTable.RemoveRow(i)
	}

	// Add current instances
	for i, instance := range instancesList {
		instancesTable.SetCell(i+1, 0, tview.NewTableCell(instance.Name))
		instancesTable.SetCell(i+1, 1, tview.NewTableCell(instance.Product.Name))
		removeButton := "[red]Remove[-]"
		instancesTable.SetCell(i+1, 2, tview.NewTableCell(removeButton).SetSelectable(true))
	}
}
