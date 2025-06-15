package tviewapp

import (
	"github.com/rivo/tview"
)

// ContentContainer is a Flex that always contains exactly one item (the content view)
// and maintains a navigation stack
type ContentContainer struct {
	*tview.Flex
	navigationStack []struct {
		View    tview.Primitive
		Factory func() tview.Primitive // optional, can be nil
	}
}

func NewContentContainer(content tview.Primitive) *ContentContainer {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(content, 0, 3, false)
	return &ContentContainer{
		Flex: flex,
		navigationStack: []struct {
			View    tview.Primitive
			Factory func() tview.Primitive
		}{
			{View: content, Factory: nil},
		},
	}
}

func NewContentContainerWithFactory(factory func() tview.Primitive) *ContentContainer {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	content := factory()
	flex.AddItem(content, 0, 3, false)
	return &ContentContainer{
		Flex: flex,
		navigationStack: []struct {
			View    tview.Primitive
			Factory func() tview.Primitive
		}{
			{View: content, Factory: factory},
		},
	}
}

// SetContent replaces the current content with a new view (clears navigation stack)
func (c *ContentContainer) SetContent(content tview.Primitive) {
	c.Clear()
	c.AddItem(content, 0, 3, false)
	c.navigationStack = []struct {
		View    tview.Primitive
		Factory func() tview.Primitive
	}{{
		View:    content,
		Factory: nil,
	}}
}

func (c *ContentContainer) SetContentWithFactory(factory func() tview.Primitive) {
	c.Clear()
	content := factory()
	c.AddItem(content, 0, 3, false)
	c.navigationStack = []struct {
		View    tview.Primitive
		Factory func() tview.Primitive
	}{{
		View:    content,
		Factory: factory,
	}}
}

// PushContent adds a new view to the stack and displays it
func (c *ContentContainer) PushContent(content tview.Primitive) {
	c.Clear()
	c.AddItem(content, 0, 3, false)
	c.navigationStack = append(c.navigationStack, struct {
		View    tview.Primitive
		Factory func() tview.Primitive
	}{View: content, Factory: nil})
}

// PopContent removes the current view and shows the previous one
// Returns true if navigation was successful, false if already at root
func (c *ContentContainer) PopContent() bool {
	if len(c.navigationStack) <= 1 {
		return false // Already at root, can't go back
	}

	// Remove current view from stack
	c.navigationStack = c.navigationStack[:len(c.navigationStack)-1]

	// Show previous view
	previousView := c.navigationStack[len(c.navigationStack)-1]
	c.Clear()
	if previousView.Factory != nil {
		c.AddItem(previousView.Factory(), 0, 3, false)
	} else {
		c.AddItem(previousView.View, 0, 3, false)
	}

	return true
}
