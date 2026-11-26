package customers

import "github.com/golang-gui/goui/ui"

func View(visible bool) ui.View {
	return ui.VBox(
		ui.Label("Customer List"),
		// Customer form and table will go here.
	).ID("customers-screen").Visible(visible)
}
