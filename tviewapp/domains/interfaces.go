package domains

import (
	"github.com/google/uuid"
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

// InstanceDetailScreenBuilder is a function type for creating instance detail screens
type InstanceDetailScreenBuilder func(instanceID uuid.UUID) tview.Primitive
