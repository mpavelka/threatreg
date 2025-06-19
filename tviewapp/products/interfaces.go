package products

import (
	"github.com/rivo/tview"
)

// ContentContainer interface defines what the instances package needs from tviewapp

type ContentContainer interface {
	PushContent(content tview.Primitive)
	PopContent() bool
	PushContentWithFactory(factory func() tview.Primitive)
}
