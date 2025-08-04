package views

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"
	"threatreg/tviewapp/threats/modals"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// NewComponentLevelThreatResolverManager creates a threat resolver management view for a specific instance
func NewComponentLevelThreatResolverManager(assignment models.ThreatAssignment, resolverComponentId uuid.UUID, contentContainer ContentContainer) tview.Primitive {
	resolver, err := service.GetComponent(resolverComponentId)
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
func buildTopSection(assignment models.ThreatAssignment, resolver *models.Component, contentContainer ContentContainer) *tview.Flex {
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
	if assignment.ComponentID != uuid.Nil {
		instanceName = assignment.Component.Name
		info.SetText(fmt.Sprintf("Component: %s\nThreat: %s\n\n%s",
			instanceName,
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
	resolverComponent *models.Component,
	contentContainer ContentContainer,
) *tview.Flex {
	column := tview.NewFlex().SetDirection(tview.FlexRow)

	// Resolver information
	resolverInfo := tview.NewTextView()
	resolverInfo.SetBorder(true).SetTitle("Resolver Information")
	resolverInfo.SetText(fmt.Sprintf("Component: %s\nProduct: %s",
		resolverComponent.Name, "Currently no way to determine product name")) // Product name not available in this context

	// Resolution status
	resolutionInfo := buildResolutionSection(assignment, resolverComponent)

	// Action buttons
	actions := buildActionBar(assignment, resolverComponent, contentContainer)

	column.AddItem(resolverInfo, 3, 0, false)
	column.AddItem(resolutionInfo, 5, 0, false)
	column.AddItem(actions, 3, 0, false)

	return column
}

// buildResolutionSection creates the resolution information display
func buildResolutionSection(assignment models.ThreatAssignment, resolverComponent *models.Component) *tview.TextView {
	section := tview.NewTextView()
	section.SetBorder(true).SetTitle("Resolution")

	resolution, err := service.GetComponentLevelThreatResolution(assignment.ID, resolverComponent.ID)
	if err == nil && resolution != nil {
		// Check if this resolution has been delegated by looking for delegation info
		targetResolution, err := service.GetDelegatedToResolutionByDelegatedByID(resolution.ID)
		if err == nil && targetResolution != nil {
			section.SetText(fmt.Sprintf("Delegated to: %s (%s)",
				targetResolution.Component.Name,
				targetResolution.ThreatAssignment.Threat.Title))
		} else {
			section.SetText(fmt.Sprintf("Status: %s\nDescription: %s",
				string(resolution.Status), resolution.Description))
		}
	} else {
		section.SetText("No resolution assigned yet. Click 'Edit Resolution' to create one.")
	}

	return section
}

// buildActionBar creates the action buttons bar
func buildActionBar(assignment models.ThreatAssignment, resolverComponent *models.Component, contentContainer ContentContainer) *tview.Flex {
	bar := tview.NewFlex().SetDirection(tview.FlexColumn)
	bar.SetTitle("Actions").SetBorder(true)

	editBtn := tview.NewButton("Edit Resolution").SetSelectedFunc(func() {
		resolutionWithDelegation, _ := service.GetComponentLevelThreatResolutionWithDelegation(assignment.ID, resolverComponent.ID)
		var resolutionPtr *models.ThreatAssignmentResolution
		if resolutionWithDelegation != nil {
			resolutionPtr = &resolutionWithDelegation.Resolution
		}
		contentContainer.PushContent(modals.CreateEditThreatAssignmentResolutionModal(
			assignment, resolutionPtr,
			&resolverComponent.ID,
			nil, // resolverProductId is nil for instance-level resolution
			func() { contentContainer.PopContent() },
			func() { contentContainer.PopContent() },
		))
	})

	addControlBtn := tview.NewButton("Add Control").SetSelectedFunc(func() {
		// TODO: Implement add control functionality
	})

	delegateBtn := tview.NewButton("Delegate").SetSelectedFunc(func() {
		resolution, err := service.GetComponentLevelThreatResolution(assignment.ID, resolverComponent.ID)
		if err != nil || resolution == nil {
			// TODO: Show error message that resolution must exist to delegate
			return
		}

		contentContainer.PushContent(modals.CreateThreatAssignmentDelegationModal(
			*resolution,
			func(targetResolution models.ThreatAssignmentResolution) {
				// Perform delegation using the service
				err := service.DelegateResolution(*resolution, targetResolution)
				if err == nil {
					contentContainer.PopContent() // Close modal on success
				}
				// TODO: Handle error case
			},
			func() { contentContainer.PopContent() }, // Close modal
		))
	})

	bar.AddItem(editBtn, 0, 1, false)
	bar.AddItem(addControlBtn, 0, 1, false)
	bar.AddItem(delegateBtn, 0, 1, false)
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
