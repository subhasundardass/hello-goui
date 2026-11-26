package app

import (
	"sync"

	"github.com/golang-gui/goui/ui"
)

type ScreenWindows struct {
	mu     sync.RWMutex
	app    ui.App
	opened map[string]bool
}

func NewScreenWindows(a ui.App) *ScreenWindows {
	return &ScreenWindows{
		app:    a,
		opened: make(map[string]bool),
	}
}

func (m *ScreenWindows) Open(id string) {
	if m == nil {
		return
	}

	found := false
	for _, screen := range GeneratedScreens() {
		if screen.ID == id {
			found = true
			break
		}
	}

	if !found {
		return
	}

	m.mu.Lock()
	alreadyOpen := m.opened[id]
	m.opened[id] = true
	m.mu.Unlock()

	if alreadyOpen {
		m.app.Sync(func() {
			if window := m.app.FindWindow("screen-" + id); window != nil {
				_ = window.Show()
			}
		})
		return
	}

	m.app.RequestUpdate()
}

func (m *ScreenWindows) IsOpen(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.opened[id]
}

func (m *ScreenWindows) Destroyed(id string) {
	m.mu.Lock()
	delete(m.opened, id)
	m.mu.Unlock()

	m.app.RequestUpdate()
}

func (m *ScreenWindows) WindowViews() []ui.WindowView {
	var result []ui.WindowView

	for _, screen := range GeneratedScreens() {
		if !m.IsOpen(screen.ID) {
			continue
		}

		id := screen.ID

		window := ui.Window("screen-"+id).
			Title(id).
			Size(900, 600).
			Center(true).
			Chrome(ui.WindowChromeNative).
			Content(screen.Build(true)).
			OnDestroy(func() {
				m.Destroyed(id)
			})

		result = append(result, window)
	}

	return result
}
