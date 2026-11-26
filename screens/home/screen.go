package home

import (
	"fmt"

	"github.com/golang-gui/goui/ui"
)

func Screen() ui.View {
	count := ui.MakeState(0)

	return ui.VBox(
		ui.Label("Dashboard"),
		ui.Label("Welcome to My Application"),
		ui.Label(""),
		ui.Label(fmt.Sprintf("Clicked %d times", count.Get())),
		ui.Button("Click me").OnClick(func() {
			count.Update(func(old int) int {
				return old + 1
			})
		}),
	).Spacing(15).Padding(30)
}
