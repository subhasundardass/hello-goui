package app

import (
	mainlayout "example.com/hello-goui/layouts/main"
	"example.com/hello-goui/screens/home"
	"example.com/hello-goui/styles"

	"github.com/golang-gui/goui/theme/modern"
	"github.com/golang-gui/goui/ui"
)

func Run() error {
	theme := modern.Options{
		Dark:     false,
		FontSize: 11,
	}

	var screenWindows *ScreenWindows

	return ui.Run("com.example.hello-goui", func(a ui.App) ui.RootView {
		if screenWindows == nil {
			screenWindows = NewScreenWindows(a)
		}

		// Main application window.
		main := ui.Window("main").
			Title("Welcome to GoUI").
			Size(1000, 650).
			Content(
				mainlayout.View(
					a,
					home.Screen(),
					screenWindows.Open,
				),
			)

		// Main window plus currently opened native screen windows.
		windows := []ui.WindowView{main}
		windows = append(windows, screenWindows.WindowViews()...)

		return ui.Root().
			StyleSheet(styles.App(theme)).
			Windows(windows...)
	})
}
