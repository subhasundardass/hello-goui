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
		// FontFamily: "JetBrains Mono Medium",
	}

	return ui.Run("com.example.hello-goui", func(a ui.App) ui.RootView {
		return ui.Root().
			StyleSheet(styles.App(theme)).
			Windows(
				ui.Window("main").
					Title("Welcome to GoUI").
					Size(1000, 650).
					Content(
						mainlayout.View(a, home.Screen()),
					),
			)
	})
}
