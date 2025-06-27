package views

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// NewProductLevelThreatResolverManager creates a threat resolver management view for a specific product
func NewProductLevelThreatResolverManager(assignment models.ThreatAssignment, resolverProductId uuid.UUID, contentContainer ContentContainer) tview.Primitive {
	resolver, err := service.GetProduct(resolverProductId)
	if err != nil {
		return tview.NewTextView().SetText(fmt.Sprintf("Error loading resolver product: %v", err))
	}

	// Build main layout
	main := tview.NewFlex().SetDirection(tview.FlexRow)

	// Top section with two columns
	top := buildProductTopSection(assignment, resolver, contentContainer)
	controls := buildControlsTable(assignment)

	main.AddItem(top, 11, 0, false)
	main.AddItem(controls, 0, 1, true)

	return main
}

// buildProductTopSection creates the top two-column section with threat info and product resolver details
func buildProductTopSection(assignment models.ThreatAssignment, resolver *models.Product, contentContainer ContentContainer) *tview.Flex {
	top := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Left: Threat Assignment Information
	threatInfo := buildThreatInfoSection(assignment)

	// Right: Resolver info, resolution, and actions
	resolverColumn := buildProductResolverColumn(assignment, resolver, contentContainer)

	top.AddItem(threatInfo, 0, 1, false)
	top.AddItem(resolverColumn, 0, 1, false)

	return top
}

// buildProductResolverColumn creates the right column with product resolver info, resolution, and actions
func buildProductResolverColumn(assignment models.ThreatAssignment, resolver *models.Product, contentContainer ContentContainer) *tview.Flex {
	column := tview.NewFlex().SetDirection(tview.FlexRow)

	// Resolver information
	resolverInfo := tview.NewTextView()
	resolverInfo.SetBorder(true).SetTitle("Resolver Information")
	resolverInfo.SetText(fmt.Sprintf("Product: %s\nDescription: %s",
		resolver.Name, resolver.Description))

	// Resolution status
	resolutionInfo := buildResolutionSection(assignment)

	// Action buttons
	actions := buildActionBar(assignment, contentContainer)

	column.AddItem(resolverInfo, 3, 0, false)
	column.AddItem(resolutionInfo, 5, 0, false)
	column.AddItem(actions, 3, 0, false)

	return column
}
