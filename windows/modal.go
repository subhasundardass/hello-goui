package windows

import (
	"context"
	"fmt"
)

type Modal[T any] struct {
	manager  *Manager
	instance *Instance
}

func (m *Modal[T]) ID() WindowID { return m.instance.ID }

func (m *Modal[T]) Close() error {
	return m.manager.Close(m.instance.ID)
}

func (m *Modal[T]) Await(ctx context.Context) (T, error) {
	var zero T

	value, err := m.instance.Results.AwaitContext(ctx)
	if err != nil {
		if ctx != nil && ctx.Err() != nil {
			// Do not leave a modal window behind after its caller has stopped
			// waiting. A concurrent close is harmless.
			_ = m.manager.Close(m.instance.ID)
			return zero, ctx.Err()
		}
		return zero, err
	}
	if value == nil {
		return zero, nil
	}
	result, ok := value.(T)
	if !ok {
		return zero, ErrInvalidResult
	}
	return result, nil
}

func OpenModal[P, R any](
	ctx context.Context,
	manager *Manager,
	screenID ScreenID,
	params P,
	parent WindowID,
) (*Modal[R], error) {
	instance, err := manager.Open(ctx, screenID, params, WindowOptions{
		Mode:   WindowModal,
		Parent: parent,
	})
	if err != nil {
		return nil, err
	}
	return &Modal[R]{manager: manager, instance: instance}, nil
}

func OpenWindow[P any](
	ctx context.Context,
	manager *Manager,
	screenID ScreenID,
	params P,
	options WindowOptions,
) (*Instance, error) {
	if options.Mode == WindowModal {
		return nil, fmt.Errorf("use OpenModal for modal windows")
	}
	return manager.Open(ctx, screenID, params, options)
}
