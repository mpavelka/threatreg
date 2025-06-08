package tviewapp

import (
	"github.com/rivo/tview"
)

// Navbar is a Flex that supports Tab/Shift+Tab navigation between its child buttons
// (for use as the navbar)
type Navbar struct {
	*tview.Flex
}

func NewNavbar(contentContainer *ContentContainer) *Navbar {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	buttons := []*tview.Button{
		tview.NewButton("Products").
			SetSelectedFunc(func() {
				contentContainer.SetContent(NewProductsView())
			}),
		tview.NewButton("Applications").
			SetSelectedFunc(func() {
				contentContainer.SetContent(NewApplicationsView(contentContainer))
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
