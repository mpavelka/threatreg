package tviewapp

import (
	"threatreg/tviewapp/domains"
	"threatreg/tviewapp/instances"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// Navbar is a Flex that supports Tab/Shift+Tab navigation between its child buttons
// (for use as the navbar)
type Navbar struct {
	*tview.Flex
}

func NewNavbar(contentContainer *ContentContainer) *Navbar {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Add Back button first
	backButton := tview.NewButton("← Back").
		SetSelectedFunc(func() {
			contentContainer.PopContent()
		})
	flex.AddItem(backButton, 8, 0, false)

	// Add a spacer
	flex.AddItem(tview.NewBox(), 1, 0, false)

	buttons := []*tview.Button{
		tview.NewButton("Domains").
			SetSelectedFunc(func() {
				contentContainer.SetContent(domains.NewDomainsView(contentContainer, func(instanceID uuid.UUID) tview.Primitive {
					return instances.NewInstanceDetailView(instanceID, contentContainer)
				}))
			}),
		tview.NewButton("Products").
			SetSelectedFunc(func() {
				contentContainer.SetContent(NewProductsView())
			}),
		tview.NewButton("Instances").
			SetSelectedFunc(func() {
				contentContainer.SetContent(instances.NewInstancesView(contentContainer))
			}),
		tview.NewButton("Threats").
			SetSelectedFunc(func() {
				contentContainer.SetContent(NewThreatsView())
			}),
		tview.NewButton("Controls").
			SetSelectedFunc(func() {
				contentContainer.SetContent(NewControlsView())
			}),
	}
	for _, btn := range buttons {
		flex.AddItem(btn, 0, 1, false)
	}
	return &Navbar{Flex: flex}
}
