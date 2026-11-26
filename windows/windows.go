package windows

import (
	"errors"
	"fmt"

	"github.com/golang-gui/goui/platform"
	"github.com/golang-gui/goui/ui"
)

var (
	ErrInvalidID      = errors.New("window ID cannot be empty")
	ErrWindowNotFound = errors.New("window not found")
	ErrWindowBusy     = errors.New("window already has an active modal child")
	ErrAlreadyOpen    = errors.New("window is already open as a modal")
	ErrResultType     = errors.New("modal result has an unexpected type")
	ErrNotModal       = errors.New("window is not an active modal")
	ErrUnsupported    = errors.New("native window management is unsupported")
)

type Config struct {
	ID      string
	Title   string
	Width   float32
	Height  float32
	Content ui.View
}

type modalFrame struct {
	parentID string
	accept   func(any) bool
	deliver  func(any)
}

type Manager struct {
	app ui.App

	// Key: child window ID.
	modals map[string]modalFrame

	// Key: parent window ID, value: active child window ID.
	activeChildren map[string]string
}

func New(app ui.App) *Manager {
	return &Manager{
		app:            app,
		modals:         make(map[string]modalFrame),
		activeChildren: make(map[string]string),
	}
}

// OpenModal opens an already-mounted native window over its parent.
//
// The parent is hidden while the child is active.
// The result callback is called when CloseModal is invoked.
func OpenModal[T any](
	manager *Manager,
	parentID string,
	childID string,
	onResult func(T),
) error {
	if manager == nil || manager.app == nil {
		return errors.New("window manager is nil")
	}

	if parentID == "" || childID == "" {
		return ErrInvalidID
	}

	if parentID == childID {
		return errors.New("parent and child cannot be the same window")
	}

	var result error

	manager.app.Sync(func() {
		if _, exists := manager.modals[childID]; exists {
			result = ErrAlreadyOpen
			return
		}

		if _, exists := manager.activeChildren[parentID]; exists {
			result = ErrWindowBusy
			return
		}

		parent := manager.app.FindWindow(parentID)
		if parent == nil {
			result = fmt.Errorf("%w: %s", ErrWindowNotFound, parentID)
			return
		}

		child := manager.app.FindWindow(childID)
		if child == nil {
			result = fmt.Errorf("%w: %s", ErrWindowNotFound, childID)
			return
		}

		parentWindow, ok := parent.(platform.WindowManagement)
		if !ok {
			result = ErrUnsupported
			return
		}

		childWindow, ok := child.(platform.WindowManagement)
		if !ok {
			result = ErrUnsupported
			return
		}

		if err := childWindow.SetParent(parentWindow); err != nil {
			result = fmt.Errorf("set native parent: %w", err)
			return
		}

		if err := parentWindow.Hide(); err != nil {
			result = fmt.Errorf("hide parent window: %w", err)
			return
		}

		if err := childWindow.Show(); err != nil {
			_ = parentWindow.Show()
			result = fmt.Errorf("show child window: %w", err)
			return
		}

		frame := modalFrame{
			parentID: parentID,
			accept: func(value any) bool {
				_, ok := value.(T)
				return ok
			},
		}

		if onResult != nil {
			frame.deliver = func(value any) {
				if typed, ok := value.(T); ok {
					onResult(typed)
				}
			}
		}

		manager.modals[childID] = frame
		manager.activeChildren[parentID] = childID
	})

	return result
}

// CloseModal closes a modal child and returns a typed result to its parent.
func CloseModal[T any](
	manager *Manager,
	childID string,
	value T,
) error {
	if manager == nil || manager.app == nil {
		return errors.New("window manager is nil")
	}

	if childID == "" {
		return ErrInvalidID
	}

	var result error

	manager.app.Sync(func() {
		frame, exists := manager.modals[childID]
		if !exists {
			result = ErrNotModal
			return
		}

		if frame.accept != nil && !frame.accept(value) {
			result = ErrResultType
			return
		}

		child := manager.app.FindWindow(childID)
		parent := manager.app.FindWindow(frame.parentID)

		if child == nil || parent == nil {
			result = ErrWindowNotFound
			return
		}

		childWindow, ok := child.(platform.WindowManagement)
		if !ok {
			result = ErrUnsupported
			return
		}

		parentWindow, ok := parent.(platform.WindowManagement)
		if !ok {
			result = ErrUnsupported
			return
		}

		if err := childWindow.Hide(); err != nil {
			result = fmt.Errorf("hide child window: %w", err)
			return
		}

		if err := parentWindow.Show(); err != nil {
			_ = childWindow.Show()
			result = fmt.Errorf("restore parent window: %w", err)
			return
		}

		delete(manager.modals, childID)
		delete(manager.activeChildren, frame.parentID)

		if frame.deliver != nil {
			frame.deliver(value)
		}
	})

	return result
}

// CancelModal closes a modal without returning a result.
func CancelModal(manager *Manager, childID string) error {
	return CloseModal[struct{}](manager, childID, struct{}{})
}
