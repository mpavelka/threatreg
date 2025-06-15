package domains

import (
	"github.com/rivo/tview"
)

// ContentContainer interface defines what the domains package needs from tviewapp
type ContentContainer interface {
	SetContent(content tview.Primitive)
	SetContentWithFactory(factory func() tview.Primitive)
	PushContent(content tview.Primitive)
	PushContentWithFactory(factory func() tview.Primitive)
	PopContent() bool
}
