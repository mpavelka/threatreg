package tviewapp

import (
	products "threatreg/tviewapp/products/views"

	"github.com/rivo/tview"
)

// RootLayout is a Flex that allows Tab/Shift+Tab to switch focus between content and navbar
// (contentContainer must be focusable)
type RootLayout struct {
	*tview.Flex
}

// NewRootLayout creates the main application layout with content and navbar
func NewRootLayout() *RootLayout {
	contentContainer := NewContentContainer(nil)
	mainContent := products.NewProductsView(contentContainer)
	contentContainer.SetContent(mainContent)

	navbar := NewNavbar(
		contentContainer,
	)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(navbar, 2, 0, true)
	flex.AddItem(contentContainer, 0, 3, false)
	return &RootLayout{
		Flex: flex,
	}
}
