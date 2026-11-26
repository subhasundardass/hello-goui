package settings

import "github.com/golang-gui/goui/ui"

func Screen() ui.View {
	return ui.VBox(
		ui.Label("Settings"),
		ui.Label("Application settings will appear here."),
	).Spacing(15).Padding(30)
}
