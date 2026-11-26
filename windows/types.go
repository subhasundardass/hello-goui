package windows

import (
	"context"
	"errors"
	"fmt"
)

type WindowID string
type ScreenID string

type WindowMode uint8

const (
	WindowNormal WindowMode = iota
	WindowModal
)

var (
	ErrScreenNotFound = errors.New("screen not registered")
	ErrWindowClosed   = errors.New("window is closed")
	ErrInvalidResult  = errors.New("invalid window result type")
	ErrManagerClosed  = errors.New("window manager closed")
)

type WindowOptions struct {
	Title  string
	Width  float32
	Height float32
	Mode   WindowMode
	Parent WindowID
}

type Screen interface {
	View() any
}

type ScreenFactory func(params any) (Screen, error)

// NativeWindow is implemented by the application-specific GOUI adapter.
// SetOnClose must notify the manager when the user closes the native window.
// The callback must not be called merely because Close() is invoked by the
// manager, or the adapter must make that callback safe to call more than once.
type NativeWindow interface {
	Show() error
	Close() error
	SetOnClose(func()) error
}

// WindowHost connects the manager to the application's GOUI window layer.
// The host owns the actual GOUI integration and must interpret WindowMode.
// GOUI's modal input boundary is window-scoped; this interface does not
// promise OS-level modality across separate native windows.
type WindowHost interface {
	Create(
		ctx context.Context,
		options WindowOptions,
		content any,
		parent NativeWindow,
	) (NativeWindow, error)
}

func (id ScreenID) String() string { return string(id) }
func (id WindowID) String() string { return string(id) }

func newWindowID() WindowID {
	return WindowID(fmt.Sprintf("window-%d", nextWindowSequence()))
}
