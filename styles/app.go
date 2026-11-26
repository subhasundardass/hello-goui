package styles

import (
	"image/color"

	"github.com/golang-gui/goui/style"
	"github.com/golang-gui/goui/theme/modern"
)

func App(theme modern.Options) style.StyleSheet {
	return style.Sheet(
		append(
			modern.Rules(theme),

			// Footer background
			style.Name("app-footer").
				BackgroundColor(color.Gray{Y: 225}),

			// Footer text
			style.Name("footer-text").
				FontSize(8).
				TextColor(color.Gray{Y: 32}),

			// Top highlight
			style.Name("footer-top").
				BackgroundColor(color.White),

			// Bottom shadow
			style.Name("footer-bottom").
				BackgroundColor(color.Gray{Y: 128}),

			// Divider dark edge
			style.Name("footer-divider-dark").
				BackgroundColor(color.Gray{Y: 128}),

			// Divider light edge
			style.Name("footer-divider-light").
				BackgroundColor(color.White),

			// Header background
			style.Name("app-header").
				BackgroundColor(color.Gray{Y: 225}),

			style.Name("header-menu-panel").
				BackgroundColor(color.Gray{Y: 225}).
				BorderWidth(0).
				Radius(0),

			// Application title
			style.Name("header-logo").
				FontSize(10).
				FontWeight(style.FontWeightBold).
				TextColor(color.Black),

			// Menu labels
			style.Name("header-menu").
				FontSize(9).
				TextColor(color.Black),

			// Company name
			style.Name("header-company").
				FontSize(9).
				FontWeight(style.FontWeightBold).
				TextColor(color.Black),

			// User dropdown
			style.Name("header-user-menu").
				FontSize(9).
				TextColor(color.Black),

			// Embossed bottom edge
			style.Name("header-bottom-shadow").
				BackgroundColor(color.Gray{Y: 128}),

			style.Name("header-bottom-highlight").
				BackgroundColor(color.White),
		)...,
	)
}
