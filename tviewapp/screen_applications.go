package tviewapp

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/rivo/tview"
)

var (
	filterForm        *tview.Form
	applicationsTable *tview.Table
	applicationsList  []models.Application
)

// NewApplicationsView creates a tview.Primitive that lists all applications
func NewApplicationsView(contentContainer *ContentContainer) tview.Primitive {

	initFilter()
	reloadApplications()
	initApplicationsTable(contentContainer)
	updateApplicationsTable()

	// Create layout
	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(30, 0).
		AddItem(filterForm, 0, 0, 1, 1, 0, 0, true).
		AddItem(applicationsTable, 0, 1, 1, 2, 0, 0, true)
	return grid
}

func reloadApplications() {
	applications, err := service.FilterApplications(
		filterForm.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		filterForm.GetFormItemByLabel("Product").(*tview.InputField).GetText(),
	)
	if err != nil {
		applicationsTable.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Error loading applications: %v", err)))
		return
	}
	applicationsList = applications
}

func initFilter() {
	filterForm = tview.NewForm().SetHorizontal(false)
	filterForm.SetBorder(true).SetTitle("Filter").SetTitleAlign(tview.AlignLeft)
	filterForm.AddInputField("Name", "", 0, nil, nil)
	filterForm.AddInputField("Product", "", 0, nil, nil)
	filterForm.AddButton("Filter", func() {
		// Apply filters
		reloadApplications()
		updateApplicationsTable()
	})
	filterForm.AddButton("Reset", func() {
		// Reset filters
		filterForm.GetFormItemByLabel("Name").(*tview.InputField).SetText("")
		filterForm.GetFormItemByLabel("Product").(*tview.InputField).SetText("")
		reloadApplications()
		updateApplicationsTable()
	})
}

func initApplicationsTable(contentContainer *ContentContainer) {
	applicationsTable = tview.NewTable().SetBorders(true)
	applicationsTable.SetTitle("Applications").SetTitleAlign(tview.AlignLeft).SetBorder(true)

	// Header (use color for bold effect)
	applicationsTable.SetFixed(1, 0) // Keep header fixed
	applicationsTable.SetSelectable(true, false)
	applicationsTable.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return // header
		}
		contentContainer.SetContent(NewApplicationDetailScreen(applicationsList[row-1].ID))
	})
}

func updateApplicationsTable() {
	applicationsTable.Clear()

	// Header
	applicationsTable.SetCell(0, 0, tview.NewTableCell("[::b]Name"))
	applicationsTable.SetCell(0, 1, tview.NewTableCell("[::b]Product"))

	// Data
	for i, a := range applicationsList {
		applicationsTable.SetCell(i+1, 0, tview.NewTableCell(a.Name))
		applicationsTable.SetCell(i+1, 1, tview.NewTableCell(a.Product.Name))
	}
}
