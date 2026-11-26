package customers

import "github.com/golang-gui/goui/ui"

func Screen() ui.View {
	return ui.VBox(
		ui.Label("Customers"),
		ui.Label("Manage your customers here."),
	).Spacing(15).Padding(30)
}
