package windows

import (
	"context"
	"errors"

	"github.com/golang-gui/goui/ui"
)

// WindowCreator creates a native window using the application's
// GOUI-specific window creation implementation.
//
// Keep GOUI-specific creation logic outside the Manager.
type WindowCreator func(
	ctx context.Context,
	options WindowOptions,
	content ui.View,
	parent NativeWindow,
) (NativeWindow, error)

// FactoryHost implements WindowHost using an injected window creator.
type FactoryHost struct {
	create WindowCreator
}

// NewFactoryHost creates a host backed by the supplied window creator.
func NewFactoryHost(create WindowCreator) *FactoryHost {
	return &FactoryHost{
		create: create,
	}
}

// Create creates a native window.
func (h *FactoryHost) Create(
	ctx context.Context,
	options WindowOptions,
	content ui.View,
	parent NativeWindow,
) (NativeWindow, error) {
	if h == nil || h.create == nil {
		return nil, errors.New("window creator is not configured")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if content == nil {
		return nil, ErrInvalidScreen
	}

	native, err := h.create(ctx, options, content, parent)
	if err != nil {
		return nil, err
	}

	if native == nil {
		return nil, ErrInvalidWindow
	}

	return native, nil
}
