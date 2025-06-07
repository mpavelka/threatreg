package tviewapp

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NavbarFlex is a Flex that supports Tab/Shift+Tab navigation between its child buttons
// (for use as the navbar)
type NavbarFlex struct {
	*tview.Flex
	buttons []tview.Primitive
}

func NewNavbarFlex(buttons ...*tview.Button) *NavbarFlex {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	prims := make([]tview.Primitive, len(buttons))
	for i, btn := range buttons {
		flex.AddItem(btn, 0, 1, false)
		prims[i] = btn
	}
	return &NavbarFlex{Flex: flex, buttons: prims}
}

func (n *NavbarFlex) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	handler := n.Flex.InputHandler()
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Handle left/right arrow keys to move focus between buttons
		if event.Key() == tcell.KeyRight || event.Key() == tcell.KeyLeft {
			for i, btn := range n.buttons {
				if btn.HasFocus() {
					var next int
					if event.Key() == tcell.KeyRight {
						next = (i + 1) % len(n.buttons)
					} else {
						next = (i - 1 + len(n.buttons)) % len(n.buttons)
					}
					setFocus(n.buttons[next])
					return
				}
			}
			// If none focused, focus first
			if len(n.buttons) > 0 {
				setFocus(n.buttons[0])
			}
			return
		}
		if handler != nil {
			handler(event, setFocus)
		}
	}
}
