package threats

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"
	"threatreg/tviewapp/common"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

type ContentContainer interface {
	PushContent(content tview.Primitive)
	PopContent() bool
	PushContentWithFactory(factory func() tview.Primitive)
}

func NewThreatsView(contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Threats")

	newThreatButton := tview.NewButton("New Threat").
		SetSelectedFunc(func() {
			contentContainer.PushContent(createEditThreatModal(
				"", "",
				func(title, description string) {
					if _, err := service.CreateThreat(title, description); err == nil {
						contentContainer.PopContent()
					}
				},
				func() { contentContainer.PopContent() },
			))
		})

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)
	actionBar.AddItem(newThreatButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false)

	threats, err := service.ListThreats()
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading threats: %v", err))
	}

	table := tview.NewTable().SetBorders(true)
	table.SetTitle("Threats").SetBorder(true)
	table.SetFixed(1, 0)

	table.SetCell(0, 0, tview.NewTableCell("[::b]Title").SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("[::b]Description").SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("[::b]Actions").SetSelectable(false))

	for i, t := range threats {
		table.SetCell(i+1, 0, tview.NewTableCell(t.Title))
		table.SetCell(i+1, 1, tview.NewTableCell(t.Description))
		removeButton := "[red]Remove[-]"
		table.SetCell(i+1, 2, tview.NewTableCell(removeButton).SetSelectable(true))
	}

	table.SetSelectable(true, true)
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(threats) && column == 2 {
			threat := threats[row-1]
			contentContainer.PushContent(common.CreateConfirmationModal(
				"Remove Threat",
				fmt.Sprintf("Are you sure you want to remove threat '%s'?", threat.Title),
				func() {
					err := service.DeleteThreat(threat.ID)
					if err != nil {
						return
					}
					contentContainer.PopContent()
				},
				func() {
					contentContainer.PopContent()
				},
			))
		}
	})

	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(table, 0, 1, true)

	return flex
}

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

	instanceName := "(Inherited from Product)"
	if assignment.InstanceID != uuid.Nil {
		instanceName = assignment.Instance.Name
	}

	info.SetText(fmt.Sprintf("Instance: %s\nThreat: %s\n\n%s",
		instanceName,
		assignment.Threat.Title,
		assignment.Threat.Description))

	return info
}

// buildResolverColumn creates the right column with resolver info, resolution, and actions
func buildResolverColumn(assignment models.ThreatAssignment, resolver *models.Instance, contentContainer ContentContainer) *tview.Flex {
	column := tview.NewFlex().SetDirection(tview.FlexRow)

	// Resolver information
	resolverInfo := tview.NewTextView()
	resolverInfo.SetBorder(true).SetTitle("Resolver Information")
	resolverInfo.SetText(fmt.Sprintf("Instance: %s\nProduct: %s",
		resolver.Name, resolver.Product.Name))

	// Resolution status
	resolutionInfo := buildResolutionSection(assignment)

	// Action buttons
	actions := buildActionBar(assignment, contentContainer)

	column.AddItem(resolverInfo, 3, 0, false)
	column.AddItem(resolutionInfo, 5, 0, false)
	column.AddItem(actions, 3, 0, false)

	return column
}

// buildResolutionSection creates the resolution information display
func buildResolutionSection(assignment models.ThreatAssignment) *tview.TextView {
	section := tview.NewTextView()
	section.SetBorder(true).SetTitle("Resolution")

	resolution, err := service.GetThreatResolutionByThreatAssignmentID(assignment.ID)
	if err == nil && resolution != nil {
		section.SetText(fmt.Sprintf("Status: %s\nDescription: %s",
			string(resolution.Status), resolution.Description))
	} else {
		section.SetText("No resolution assigned yet. Click 'Edit Resolution' to create one.")
	}

	return section
}

// buildActionBar creates the action buttons bar
func buildActionBar(assignment models.ThreatAssignment, contentContainer ContentContainer) *tview.Flex {
	bar := tview.NewFlex().SetDirection(tview.FlexColumn)
	bar.SetTitle("Actions").SetBorder(true)

	editBtn := tview.NewButton("Edit Resolution").SetSelectedFunc(func() {
		resolution, _ := service.GetThreatResolutionByThreatAssignmentID(assignment.ID)
		contentContainer.PushContent(createEditThreatAssignmentResolutionModal(
			assignment, resolution,
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
