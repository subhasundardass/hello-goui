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
	ErrWindowNotFound = errors.New("window not found")
	ErrWindowClosed   = errors.New("window already closed")
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

// NativeWindow is implemented by the GOUI window adapter.
type NativeWindow interface {
	Show() error
	Close() error
	SetModal(modal bool) error
}

// WindowHost connects the manager to GOUI.
type WindowHost interface {
	Create(
		ctx context.Context,
		options WindowOptions,
		content any,
		parent NativeWindow,
	) (NativeWindow, error)
}

func (id ScreenID) String() string {
	return string(id)
}

func (id WindowID) String() string {
	return string(id)
}

func newWindowID() WindowID {
	return WindowID(fmt.Sprintf("window-%d", nextWindowSequence()))
}
