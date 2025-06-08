package tviewapp

import (
	"github.com/rivo/tview"
)

// ContentContainer is a Flex that always contains exactly one item (the content view)
type ContentContainer struct {
	*tview.Flex
}

func NewContentContainer(content tview.Primitive) *ContentContainer {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(content, 0, 3, false)
	return &ContentContainer{Flex: flex}
}

// SetContent replaces the current content with a new view
func (c *ContentContainer) SetContent(content tview.Primitive) {
	c.Clear()
	c.AddItem(content, 0, 3, false)
}
