package views

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"
	"threatreg/tviewapp/threats/modals"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// NewInstanceLevelThreatResolverManager creates a threat resolver management view for a specific instance
func NewInstanceLevelThreatResolverManager(assignment models.ThreatAssignment, resolverInstanceId uuid.UUID, contentContainer ContentContainer) tview.Primitive {
	resolver, err := service.GetInstance(resolverInstanceId)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading resolver instance: %v", err))
	}

	// Build main layout
	main := tview.NewFlex().SetDirection(tview.FlexRow)

	// Top section with two columns
	top := buildTopSection(assignment, resolver, contentContainer)
	controls := buildControlsTable(assignment)

	main.AddItem(top, 11, 0, false)
	main.AddItem(controls, 0, 1, true)

	return main
}

// buildTopSection creates the top two-column section with threat info and resolver details
func buildTopSection(assignment models.ThreatAssignment, resolver *models.Instance, contentContainer ContentContainer) *tview.Flex {
	top := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Left: Threat Assignment Information
	threatInfo := buildThreatInfoSection(assignment)

	// Right: Resolver info, resolution, and actions
	resolverColumn := buildResolverColumn(assignment, resolver, contentContainer)

	top.AddItem(threatInfo, 0, 1, false)
	top.AddItem(resolverColumn, 0, 1, false)

	return top
}

// buildThreatInfoSection creates the threat assignment information section
func buildThreatInfoSection(assignment models.ThreatAssignment) *tview.TextView {
	info := tview.NewTextView()
	info.SetBorder(true).SetTitle("Threat Assignment Information")

	var instanceName string
	if assignment.InstanceID != uuid.Nil {
		instanceName = assignment.Instance.Name
		info.SetText(fmt.Sprintf("Instance: %s\nThreat: %s\n\n%s",
			instanceName,
			assignment.Threat.Title,
			assignment.Threat.Description))
	} else if assignment.ProductID != uuid.Nil {
		info.SetText(fmt.Sprintf("Product: %s\nThreat: %s\n\n%s",
			assignment.Product.Name,
			assignment.Threat.Title,
			assignment.Threat.Description))
	} else {
		info.SetText(fmt.Sprintf("Threat: %s\n\n%s",
			assignment.Threat.Title,
			assignment.Threat.Description))
	}

	return info
}

// buildResolverColumn creates the right column with resolver info, resolution, and actions
func buildResolverColumn(
	assignment models.ThreatAssignment,
	resolverInstance *models.Instance,
	contentContainer ContentContainer,
) *tview.Flex {
	column := tview.NewFlex().SetDirection(tview.FlexRow)

	// Resolver information
	resolverInfo := tview.NewTextView()
	resolverInfo.SetBorder(true).SetTitle("Resolver Information")
	resolverInfo.SetText(fmt.Sprintf("Instance: %s\nProduct: %s",
		resolverInstance.Name, resolverInstance.Product.Name))

	// Resolution status
	resolutionInfo := buildResolutionSection(assignment, resolverInstance)

	// Action buttons
	actions := buildActionBar(assignment, resolverInstance, contentContainer)

	column.AddItem(resolverInfo, 3, 0, false)
	column.AddItem(resolutionInfo, 5, 0, false)
	column.AddItem(actions, 3, 0, false)

	return column
}

// buildResolutionSection creates the resolution information display
func buildResolutionSection(assignment models.ThreatAssignment, resolverInstance *models.Instance) *tview.TextView {
	section := tview.NewTextView()
	section.SetBorder(true).SetTitle("Resolution")

	resolution, err := service.GetInstanceLevelThreatResolution(assignment.ID, resolverInstance.ID)
	if err == nil && resolution != nil {
		section.SetText(fmt.Sprintf("Status: %s\nDescription: %s",
			string(resolution.Status), resolution.Description))
	} else {
		section.SetText("No resolution assigned yet. Click 'Edit Resolution' to create one.")
	}

	return section
}

// buildActionBar creates the action buttons bar
func buildActionBar(assignment models.ThreatAssignment, resolverInstance *models.Instance, contentContainer ContentContainer) *tview.Flex {
	bar := tview.NewFlex().SetDirection(tview.FlexColumn)
	bar.SetTitle("Actions").SetBorder(true)

	editBtn := tview.NewButton("Edit Resolution").SetSelectedFunc(func() {
		resolution, _ := service.GetInstanceLevelThreatResolution(assignment.ID, resolverInstance.ID)
		contentContainer.PushContent(modals.CreateEditThreatAssignmentResolutionModal(
			assignment, resolution,
			&resolverInstance.ID,
			nil, // resolverProductId is nil for instance-level resolution
			func() { contentContainer.PopContent() },
			func() { contentContainer.PopContent() },
		))
	})

	addControlBtn := tview.NewButton("Add Control").SetSelectedFunc(func() {
		// TODO: Implement add control functionality
	})

	bar.AddItem(editBtn, 0, 1, false)
	bar.AddItem(addControlBtn, 0, 1, false)
	bar.AddItem(tview.NewBox(), 0, 1, false) // Spacer

	return bar
}

// buildControlsTable creates the controls table section
func buildControlsTable(assignment models.ThreatAssignment) *tview.Table {
	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Controls").SetBorder(true)
	table.SetSelectable(true, false)

	// Headers
	table.SetCell(0, 0, tview.NewTableCell("[::b]Title").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Description").SetSelectable(false))

	// Data rows
	if len(assignment.ControlAssignments) > 0 {
		for i, ctrl := range assignment.ControlAssignments {
			table.SetCell(i+1, 0, tview.NewTableCell(ctrl.Control.Title))
			table.SetCell(i+1, 1, tview.NewTableCell(ctrl.Control.Description))
		}
	} else {
		table.SetCell(1, 0, tview.NewTableCell("No controls assigned"))
		table.SetCell(1, 1, tview.NewTableCell(""))
	}

	return table
}
