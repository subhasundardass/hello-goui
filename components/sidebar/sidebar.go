package sidebar

import "github.com/golang-gui/goui/ui"

func View() ui.View {
	return ui.VBox(
		ui.Label("MY APP"),
		ui.Label(""),
		ui.Label("Home"),
		ui.Label("Customers"),
		ui.Label("Settings"),
	).Spacing(15).Padding(20)
}
