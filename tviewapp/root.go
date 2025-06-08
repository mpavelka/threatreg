package tviewapp

import (
	"threatreg/internal/database"

	"github.com/rivo/tview"
)

func Run() {
	// Connect to the database
	database.Connect()

	app := tview.NewApplication()
	app.EnableMouse(true)
	app.SetRoot(NewRootLayout(), true)

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}
