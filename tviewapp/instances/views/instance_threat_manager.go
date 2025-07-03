package views

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"
	instancesModals "threatreg/tviewapp/instances/modals"
	productsModals "threatreg/tviewapp/products/modals"
	threatsViews "threatreg/tviewapp/threats/views"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// NewInstanceThreatManager creates a threat management view for an instance
func NewInstanceThreatManager(instanceID uuid.UUID, contentContainer ContentContainer) tview.Primitive {
	instance, err := service.GetInstance(instanceID)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading instance: %v", err))
	}

	// Left column - Instance section
	instanceText := tview.NewTextView()
	instanceText.SetBorder(true).SetTitle("Instance Information")
	instanceText.SetText(fmt.Sprintf("Name: %s\n", instance.Name))

	instanceButton := tview.NewButton("Edit Instance").SetSelectedFunc(func() {
		contentContainer.PushContent(instancesModals.CreateEditInstanceModal(
			instance.Name,
			func(name string) {
				_, err := service.UpdateInstance(instance.ID, &name, nil)
				if err != nil {
					// TODO: Show error message
					return
				}
				// Refresh the view
				contentContainer.PopContent()
				contentContainer.PopContent()
				contentContainer.PushContentWithFactory(func() tview.Primitive {
					return NewInstanceThreatManager(instanceID, contentContainer)
				})
			},
			func() {
				contentContainer.PopContent()
			},
		))
	})

	instanceSection := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(instanceText, 5, 0, false).
		AddItem(instanceButton, 1, 0, true)

	// Left column - Product section
	productText := tview.NewTextView()
	productText.SetBorder(true).SetTitle("Product Information")
	productText.SetText(fmt.Sprintf("Name: %s\nDescription: %s", instance.Product.Name, instance.Product.Description))

	productButton := tview.NewButton("Edit Product").SetSelectedFunc(func() {
		contentContainer.PushContent(productsModals.CreateEditProductModal(
			instance.Product.Name,
			instance.Product.Description,
			func(name, description string) {
				_, err := service.UpdateProduct(instance.Product.ID, &name, &description)
				if err != nil {
					// TODO: Show error message
					return
				}
				// Refresh the view
				contentContainer.PopContent()
				contentContainer.PopContent()
				contentContainer.PushContentWithFactory(func() tview.Primitive {
					return NewInstanceThreatManager(instanceID, contentContainer)
				})
			},
			func() {
				contentContainer.PopContent()
			},
		))
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
		contentContainer.PushContent(instancesModals.CreateInstanceSelectThreatModal(instance.ID, func() {
			contentContainer.PopContent()
		}))
	}), 0, 1, false)
	actionBar.AddItem(tview.NewButton("Add Product Threat").SetSelectedFunc(func() {
		contentContainer.PushContent(instancesModals.CreateProductSelectThreatModal(instance.Product.ID, func() {
			contentContainer.PopContent()
		}))
	}), 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false) // Spacer

	// Right column - Threat Assignments table
	threatTable := tview.NewTable().SetBorders(true)
	threatTable.SetTitle("Threat Assignments").SetBorder(true)
	threatTable.SetCell(0, 0, tview.NewTableCell("[::b]Type").SetSelectable(false))
	threatTable.SetCell(0, 1, tview.NewTableCell("[::b]Resolution").SetSelectable(false))
	threatTable.SetCell(0, 2, tview.NewTableCell("[::b]Threat").SetSelectable(false))
	threatTable.SetCell(0, 3, tview.NewTableCell("[::b]Description").SetSelectable(false))

	// Load actual threat assignments for this instance
	instanceAssignments, err := service.ListThreatAssignmentsByInstanceID(instance.ID)
	if err != nil {
		// Show error in table
		threatTable.SetCell(1, 0, tview.NewTableCell("Error"))
		threatTable.SetCell(1, 1, tview.NewTableCell(fmt.Sprintf("Failed to load: %v", err)))
		threatTable.SetCell(1, 2, tview.NewTableCell(""))
		threatTable.SetCell(1, 3, tview.NewTableCell(""))
	} else {
		// Load product-level assignments as well (instance inherits product threats)
		productAssignments, err := service.ListThreatAssignmentsByProductID(instance.Product.ID)
		if err != nil {
			productAssignments = []models.ThreatAssignment{}
		}

		row := 1

		// Store all assignments for navigation
		allAssignments := append(instanceAssignments, productAssignments...)

		// Add instance-specific threats
		for _, assignment := range instanceAssignments {
			threatTable.SetCell(row, 0, tview.NewTableCell("Instance"))
			threatTable.SetCell(row, 1, tview.NewTableCell(getResolutionStatus(assignment)))
			threatTable.SetCell(row, 2, tview.NewTableCell(assignment.Threat.Title))
			threatTable.SetCell(row, 3, tview.NewTableCell(assignment.Threat.Description))
			row++
		}

		// Add product-level threats (inherited)
		for _, assignment := range productAssignments {
			threatTable.SetCell(row, 0, tview.NewTableCell("Product"))
			threatTable.SetCell(row, 1, tview.NewTableCell(getResolutionStatus(assignment)))
			threatTable.SetCell(row, 2, tview.NewTableCell(assignment.Threat.Title))
			threatTable.SetCell(row, 3, tview.NewTableCell(assignment.Threat.Description))
			row++
		}

		// Make table selectable and add click handler
		threatTable.SetSelectable(true, false)
		threatTable.SetSelectedFunc(func(row, column int) {
			if row > 0 && row-1 < len(allAssignments) {
				assignment := allAssignments[row-1]
				contentContainer.PushContentWithFactory(func() tview.Primitive {
					return threatsViews.NewInstanceLevelThreatResolverManager(assignment, instanceID, contentContainer)
				})
			}
		})

		// Show message if no threats are assigned
		if len(instanceAssignments) == 0 && len(productAssignments) == 0 {
			threatTable.SetCell(1, 0, tview.NewTableCell(""))
			threatTable.SetCell(1, 1, tview.NewTableCell(""))
			threatTable.SetCell(1, 2, tview.NewTableCell("No threats assigned"))
			threatTable.SetCell(1, 3, tview.NewTableCell(""))
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

func getResolutionStatus(assignment models.ThreatAssignment) string {
	// TODO: If this is a product-level assignment, we need to check both instance- and product-level resolutions
	// For now, we only check the first available resolution, which is obviously wrong (it might be instance or product)

	resolution, err := service.GetThreatResolutionByThreatAssignmentID(assignment.ID)
	if err != nil || resolution == nil {
		return "-"
	}

	switch resolution.Status {
	case models.ThreatAssignmentResolutionStatusResolved:
		return "[green]Resolved[-]"
	case models.ThreatAssignmentResolutionStatusAwaiting:
		return "[yellow]Awaiting[-]"
	case models.ThreatAssignmentResolutionStatusAccepted:
		return "[blue]Accepted[-]"
	default:
		return "Unknown"
	}
}
