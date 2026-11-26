package header

import (
	"image/color"
	"strings"

	"example.com/hello-goui/components"

	"github.com/golang-gui/goui/layout"
	"github.com/golang-gui/goui/style"
	"github.com/golang-gui/goui/ui"
)

const (
	CompanyName = "ABC Company"
	UserName    = "Administrator"
)

const (
	menuButtonStyle = "header-menu-button"
	headerBarStyle  = "app-header"
	logoStyle       = "header-logo"
	companyStyle    = "header-company"

	barPadding   float32 = 0
	menuPadX     float32 = 5
	menuPadY     float32 = 0
	headerHeight float32 = 20
)

// -----------------------------------------------------------------------------
// Stylesheet
// -----------------------------------------------------------------------------

func Rules(s ui.Settings) []style.Rule {
	fg := color.NRGBA{
		R: 35, G: 38, B: 42, A: 255,
	}

	barBG := color.NRGBA{
		R: 245, G: 246, B: 248, A: 255,
	}

	hover := color.NRGBA{
		R: 100, G: 110, B: 125, A: 35,
	}

	pressed := color.NRGBA{
		R: 100, G: 110, B: 125, A: 65,
	}

	brand := color.NRGBA{
		R: 35, G: 105, B: 175, A: 255,
	}

	if s.ColorScheme() == ui.ColorSchemeDark {
		fg = color.NRGBA{
			R: 230, G: 232, B: 236, A: 255,
		}

		barBG = color.NRGBA{
			R: 35, G: 37, B: 41, A: 255,
		}

		hover = color.NRGBA{
			R: 150, G: 160, B: 175, A: 35,
		}

		pressed = color.NRGBA{
			R: 150, G: 160, B: 175, A: 65,
		}

		brand = color.NRGBA{
			R: 105, G: 175, B: 235, A: 255,
		}
	}

	var rules []style.Rule

	// Top menu buttons.
	rules = append(rules,
		style.Name(menuButtonStyle).
			ForegroundColor(fg).
			FontFamily(s.FontFamily()).
			FontSize(s.FontSize()).
			BackgroundColor(color.Transparent).
			BorderWidth(0).
			Radius(0),

		style.Name(menuButtonStyle).
			State(style.Hovered).
			BackgroundColor(hover).
			BorderWidth(0),

		style.Name(menuButtonStyle).
			State(style.Pressed).
			BackgroundColor(pressed).
			BorderWidth(0),
	)

	// Header background.
	rules = append(rules,
		style.Name(headerBarStyle).
			BackgroundColor(barBG).
			BorderWidth(0).
			Radius(0),
	)

	// Application logo.
	rules = append(rules,
		style.Name(logoStyle).
			ForegroundColor(brand).
			FontFamily(s.FontFamily()).
			FontSize(s.FontSize()).
			BackgroundColor(color.Transparent).
			BorderWidth(0).
			Radius(0),
	)

	// Company name.
	rules = append(rules,
		style.Name(companyStyle).
			ForegroundColor(fg).
			FontFamily(s.FontFamily()).
			FontSize(s.FontSize()).
			BackgroundColor(color.Transparent).
			BorderWidth(0).
			Radius(0),
	)

	return rules
}

// -----------------------------------------------------------------------------
// Menu item conversion
// -----------------------------------------------------------------------------

// buildMenuItems converts the existing components.Node tree into
// native GoUI menu entries.
//
// GoUI v0.8.0's declarative MenuItemView currently does not expose
// a public submenu constructor. Therefore, nested nodes are represented
// as disabled group headings followed by indented child entries.
func buildMenuItems(nodes []components.Node, depth int) []*ui.MenuItemView {
	var items []*ui.MenuItemView

	for _, node := range nodes {
		if node.Sep {
			items = append(items, ui.MenuSeparator())
			continue
		}

		prefix := strings.Repeat("   ", depth)
		label := prefix + node.Label

		if len(node.Children) > 0 {
			// Group heading.
			items = append(
				items,
				ui.MenuItem(label).Enabled(false),
			)

			// Separate heading from its children.
			items = append(items, ui.MenuSeparator())

			// Add child entries.
			items = append(
				items,
				buildMenuItems(node.Children, depth+1)...,
			)

			continue
		}

		// Normal menu item.
		items = append(
			items,
			ui.MenuItem(label, node.Action),
		)
	}

	return items
}

// -----------------------------------------------------------------------------
// Top-level menu button
// -----------------------------------------------------------------------------

func topMenu(title string, nodes []components.Node) ui.View {
	return ui.MenuButton(title).
		Style(menuButtonStyle).
		Menu(buildMenuItems(nodes, 0)...)
}

// -----------------------------------------------------------------------------
// User menu
// -----------------------------------------------------------------------------

func userMenu() ui.View {
	return ui.MenuButton(UserName).
		Style(menuButtonStyle).
		MinHeight(headerHeight).
		Menu(
			ui.MenuItem("My Profile", func() {
				// Open profile screen.
			}),
			ui.MenuItem("Change Password", func() {
				// Open change password screen.
			}),
			ui.MenuSeparator(),
			ui.MenuItem("Logout", func() {
				// Perform logout.
			}),
		)
}

// -----------------------------------------------------------------------------
// Header
// -----------------------------------------------------------------------------

func View(content ui.View) ui.View {
	all := components.Menus()

	// Application identity.
	brand := ui.HBox(

		ui.Label("DORII").
			Style(logoStyle),

		// ui.Label(CompanyName).
		// 	Style(companyStyle),
	).
		Spacing(12).
		CrossAlign(layout.CrossCenter)

	// Top-level menus.
	menus := make([]ui.View, 0, len(all))

	for _, m := range all {
		menus = append(
			menus,
			topMenu(m.Title, m.Items),
		)
	}

	// Left side: application identity and menus.
	left := ui.HBox(
		append(
			[]ui.View{brand},
			menus...,
		)...,
	).
		Spacing(8).
		CrossAlign(layout.CrossCenter)

		// Main header bar.
	bar := ui.HBox(
		ui.HBox(
			ui.HBox().MinWidth(10),
			left,
		),
		userMenu(),
	).
		MainAlign(layout.MainSpaceBetween).
		CrossAlign(layout.CrossCenter).
		Padding(barPadding).
		MinHeight(headerHeight).
		Style(headerBarStyle)

	// Main application layout.
	page := ui.VBox(
		bar,
		ui.Modify(content).MainWeight(1),
	).
		Spacing(0).
		CrossAlign(layout.CrossStretch)

	return page.MainWeight(1)
}
