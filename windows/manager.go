package windows

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Manager struct {
	mu       sync.RWMutex
	registry *Registry
	host     WindowHost
	windows  map[WindowID]*Instance
	closed   bool
}

var (
	ErrWindowNotFound = errors.New("window not found")
	ErrInvalidScreen  = errors.New("invalid screen")
	ErrInvalidWindow  = errors.New("invalid native window")
)

func NewManager(registry *Registry, host WindowHost) *Manager {
	if registry == nil {
		registry = NewRegistry()
	}

	return &Manager{
		registry: registry,
		host:     host,
		windows:  make(map[WindowID]*Instance),
	}
}

func (m *Manager) Registry() *Registry {
	return m.registry
}

func (m *Manager) Open(
	ctx context.Context,
	screenID ScreenID,
	params any,
	options WindowOptions,
) (*Instance, error) {
	if m == nil {
		return nil, ErrManagerClosed
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if m.host == nil {
		return nil, fmt.Errorf("window host is not configured")
	}

	instance := newInstance(screenID, options)

	// Reserve the instance and register its parent-child relationship.
	m.mu.Lock()

	if m.closed {
		m.mu.Unlock()
		return nil, ErrManagerClosed
	}

	var (
		parent         NativeWindow
		parentInstance *Instance
	)

	if options.Parent != "" {
		parentInstance = m.windows[options.Parent]

		if parentInstance == nil || parentInstance.IsClosed() {
			m.mu.Unlock()
			return nil, fmt.Errorf("parent: %w", ErrWindowNotFound)
		}

		parent = parentInstance.nativeWindow()

		if parent == nil {
			m.mu.Unlock()
			return nil, fmt.Errorf("parent window is not ready")
		}

		parentInstance.AddChild(instance.ID)
	}

	m.windows[instance.ID] = instance
	m.mu.Unlock()

	cleanup := func() {
		m.removeInstance(instance)
		instance.markClosing()
		instance.markClosed()
		instance.Results.Complete(nil, ErrWindowClosed)
	}

	// Resolve screen factory.
	factory, err := m.registry.Get(screenID)
	if err != nil {
		cleanup()
		return nil, err
	}

	// Create screen.
	screen, err := factory(params)
	if err != nil {
		cleanup()
		return nil, err
	}

	if screen == nil {
		cleanup()
		return nil, ErrInvalidScreen
	}

	instance.setScreen(screen)

	if err := ctx.Err(); err != nil {
		cleanup()
		return nil, err
	}

	// Create native window.
	native, err := m.host.Create(ctx, options, screen.View(), parent)
	if err != nil {
		cleanup()
		return nil, err
	}

	if native == nil {
		cleanup()
		return nil, ErrInvalidWindow
	}

	instance.setNative(native)

	// Check manager and parent lifecycle.
	m.mu.RLock()

	managerClosed := m.closed
	registered := m.windows[instance.ID] == instance
	parentStillOpen := true

	if options.Parent != "" {
		currentParent := m.windows[options.Parent]

		parentStillOpen = currentParent == parentInstance &&
			currentParent != nil &&
			!currentParent.IsClosed()
	}

	m.mu.RUnlock()

	if managerClosed || !registered || !parentStillOpen || ctx.Err() != nil {
		_ = native.Close()
		cleanup()

		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if managerClosed {
			return nil, ErrManagerClosed
		}

		return nil, fmt.Errorf("parent window closed while child was opening")
	}

	// Register native close callback.
	if err := native.SetOnClose(func() {
		_ = m.Close(instance.ID)
	}); err != nil {
		_ = native.Close()
		cleanup()
		return nil, err
	}

	// Transition from creating to open.
	if !instance.markOpen() {
		_ = native.Close()
		cleanup()
		return nil, ErrWindowClosed
	}

	// Show native window.
	if err := native.Show(); err != nil {
		_ = m.Close(instance.ID)
		return nil, err
	}

	return instance, nil
}

func (m *Manager) Get(id WindowID) (*Instance, bool) {
	if m == nil {
		return nil, false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	instance, ok := m.windows[id]
	return instance, ok
}

func (m *Manager) Instances() []*Instance {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	instances := make([]*Instance, 0, len(m.windows))

	for _, instance := range m.windows {
		instances = append(instances, instance)
	}

	return instances
}

func (m *Manager) removeInstance(instance *Instance) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.windows[instance.ID] == instance {
		delete(m.windows, instance.ID)
	}

	if instance.ParentID != "" {
		if parent := m.windows[instance.ParentID]; parent != nil {
			parent.RemoveChild(instance.ID)
		}
	}
}

func (m *Manager) Close(id WindowID) error {
	if m == nil {
		return ErrManagerClosed
	}

	m.mu.RLock()
	instance := m.windows[id]
	m.mu.RUnlock()

	if instance == nil {
		return ErrWindowNotFound
	}

	if !instance.markClosing() {
		return ErrWindowClosed
	}

	var firstErr error

	// Close children before the parent.
	for _, childID := range instance.Children() {
		if err := m.Close(childID); err != nil &&
			err != ErrWindowNotFound &&
			err != ErrWindowClosed &&
			firstErr == nil {
			firstErr = err
		}
	}

	// Close native window.
	if native := instance.nativeWindow(); native != nil {
		if err := native.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	instance.Results.Complete(nil, ErrWindowClosed)
	instance.markClosed()

	m.removeInstance(instance)

	return firstErr
}

func (m *Manager) CloseAll() error {
	if m == nil {
		return ErrManagerClosed
	}

	m.mu.Lock()

	m.closed = true

	// Collect root windows.
	roots := make([]WindowID, 0)

	for id, instance := range m.windows {
		if instance.ParentID == "" {
			roots = append(roots, id)
		}
	}

	m.mu.Unlock()

	var firstErr error

	for _, id := range roots {
		if err := m.Close(id); err != nil &&
			err != ErrWindowNotFound &&
			err != ErrWindowClosed &&
			firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
