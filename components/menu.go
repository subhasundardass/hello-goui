package components

type Node struct {
	Label    string
	Action   func()
	Children []Node
	Sep      bool
}

type Menu struct {
	Title string
	Items []Node
}

func Item(label string, action func()) Node {
	return Node{
		Label:  label,
		Action: action,
	}
}

func Submenu(label string, children ...Node) Node {
	return Node{
		Label:    label,
		Children: children,
	}
}

func Separator() Node {
	return Node{Sep: true}
}

func Menus(openScreen func(string)) []Menu {
	noop := func() {}

	open := func(id string) func() {
		return func() {
			if openScreen != nil {
				openScreen(id)
			}
		}
	}

	return []Menu{
		{"File", []Node{
			Item("New", noop),
			Item("Open", noop),
			Item("Save", noop),
			Separator(),
			Item("Exit", noop),
		}},
		{"Edit", []Node{
			Item("Undo", noop),
			Item("Redo", noop),
			Separator(),
			Item("Cut", noop),
			Item("Copy", noop),
			Item("Paste", noop),
		}},
		{"View", []Node{
			Item("Dashboard", noop),
			Item("Refresh", noop),
			Separator(),
			Item("Full Screen", noop),
		}},
		{"Masters", []Node{
			Item("Customers", open("customers")),
			Item("Suppliers", open("suppliers")),
			Item("Products", open("products")),
			Separator(),
			Item("Users", noop),
		}},
		{"Transactions", []Node{
			Item("Sales", noop),
			Item("Purchase", noop),
			Item("Receipts", noop),
			Item("Payments", noop),
		}},
		{"Reports", []Node{
			Submenu("Sales Report",
				Item("Daily", noop),
				Item("Monthly", noop),
				Item("By Customer", noop),
			),
			Submenu("Purchase Report",
				Item("Daily", noop),
				Item("Monthly", noop),
				Item("By Supplier", noop),
			),
			Item("Stock Report", noop),
		}},
		{"Help", []Node{
			Item("Documentation", noop),
			Item("About", noop),
		}},
	}
}
