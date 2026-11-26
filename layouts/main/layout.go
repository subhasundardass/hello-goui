package mainlayout

import (
	"example.com/hello-goui/components/footer"
	"example.com/hello-goui/components/header"
	"example.com/hello-goui/components/sidebar"

	"github.com/golang-gui/goui/layout"
	"github.com/golang-gui/goui/ui"
)

func View(app ui.App, content ui.View) ui.View {
	body := ui.VBox(
		ui.HBox(
			sidebar.View(),
			content,
		).
			Spacing(0).
			MainWeight(1),

		footer.View(app),
	).
		Spacing(0).
		CrossAlign(layout.CrossStretch)

	return header.View(body)
}
