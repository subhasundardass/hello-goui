package footer

import (
	"time"

	"github.com/golang-gui/goui/layout"
	"github.com/golang-gui/goui/ui"
)

const (
	SoftwareName    = "Hello GOUI"
	SoftwareVersion = "v1.0.0"
)

func text(value string) ui.View {
	return ui.Label(value).Style("footer-text")
}

// Two adjacent lines create an embossed divider.
func divider() ui.View {
	return ui.HBox(
		ui.VBox().MinWidth(1).Style("footer-divider-dark"),
		ui.VBox().MinWidth(1).Style("footer-divider-light"),
	).
		Spacing(0).
		CrossAlign(layout.CrossStretch)
}

func View(app ui.App) ui.View {
	// Reactive clock state.
	now := ui.MakeState(time.Now())

	// Update every second.
	app.TickFunc(time.Second, func() {
		now.Set(time.Now())
	})

	// Reactive date/time label.
	clock := ui.Build(func() ui.View {
		return text(now.Get().Format("02-Jan-2006 03:04:05 PM"))
	})

	right := ui.HBox(

		// text("CAPS: OFF"),
		// divider(),
		text("All systems operational"),
		divider(),
		text(SoftwareName+" "+SoftwareVersion),
		divider(),
		clock,
	).
		Spacing(8).
		CrossAlign(layout.CrossStretch)

	content := ui.HBox(
		text("Ready"),
		right,
	).
		MainAlign(layout.MainSpaceBetween).
		CrossAlign(layout.CrossStretch).
		Spacing(0).
		Padding(5).
		MinHeight(26).
		Style("app-footer")

	return ui.VBox(
		ui.HBox().MinHeight(1).Style("footer-top"),
		content,
		ui.HBox().MinHeight(1).Style("footer-bottom"),
	).
		Spacing(0).
		CrossAlign(layout.CrossStretch)
}
