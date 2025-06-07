package tviewapp

import (
	"threatreg/internal/database"

	"github.com/rivo/tview"
)

// NewRootLayout creates the main application layout with content and navbar
func NewRootLayout(app *tview.Application) *RootLayout {
	mainContent := NewProductsView()
	contentContainer := NewContentContainer(mainContent)

	productsBtn := tview.NewButton("Products").SetSelectedFunc(func() {
		contentContainer.SetContent(NewProductsView())
	})
	applicationsBtn := tview.NewButton("Applications").SetSelectedFunc(func() {
		contentContainer.SetContent(NewApplicationsView(contentContainer))
	})
	threatsBtn := tview.NewButton("Threats").SetSelectedFunc(func() {
		contentContainer.SetContent(NewThreatsView())
	})
	controlsBtn := tview.NewButton("Controls").SetSelectedFunc(func() {
		contentContainer.SetContent(NewControlsView())
	})

	navbar := NewNavbarFlex(
		productsBtn,
		applicationsBtn,
		threatsBtn,
		controlsBtn,
	)

	return NewRootLayoutWithContent(app, contentContainer, navbar)
}

// RootLayout is a Flex that allows Tab/Shift+Tab to switch focus between content and navbar
// (contentContainer must be focusable)
type RootLayout struct {
	*tview.Flex
	contentContainer tview.Primitive
	navbar           tview.Primitive
	app              *tview.Application
}

func NewRootLayoutWithContent(app *tview.Application, contentContainer, navbar tview.Primitive) *RootLayout {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(contentContainer, 0, 3, false)
	flex.AddItem(navbar, 1, 0, true)
	return &RootLayout{
		Flex:             flex,
		contentContainer: contentContainer,
		navbar:           navbar,
		app:              app,
	}
}

func Run() {
	// Connect to the database
	database.Connect()

	app := tview.NewApplication().EnableMouse(true)
	rootLayout := NewRootLayout(app)

	// Set the root view and initial focus to the first button in the navbar
	navbar := rootLayout.GetItem(1) // navbar is the second item in the Flex
	if navFlex, ok := navbar.(*tview.Flex); ok {
		homeBtn := navFlex.GetItem(0)
		app.SetRoot(rootLayout.Flex, true).SetFocus(homeBtn)
	} else {
		app.SetRoot(rootLayout.Flex, true)
	}

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}
