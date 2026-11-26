package windows

import (
	"context"
	"fmt"
	"sync"
)

type Manager struct {
	mu sync.RWMutex

	registry *Registry
	host     WindowHost

	windows map[WindowID]*Instance
	closed  bool
}

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
	if m.host == nil {
		return nil, fmt.Errorf("window host is not configured")
	}

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil, ErrManagerClosed
	}

	var parent NativeWindow

	if options.Parent != "" {
		parentInstance := m.windows[options.Parent]
		if parentInstance == nil || parentInstance.IsClosed() {
			m.mu.Unlock()
			return nil, fmt.Errorf("parent: %w", ErrWindowNotFound)
		}
		parent = parentInstance.Native
	}

	instance := newInstance(screenID, options)
	m.windows[instance.ID] = instance
	m.mu.Unlock()

	factory, err := m.registry.Get(screenID)
	if err != nil {
		m.removeInstance(instance)
		return nil, err
	}

	screen, err := factory(params)
	if err != nil {
		m.removeInstance(instance)
		return nil, err
	}

	instance.Screen = screen

	native, err := m.host.Create(
		ctx,
		options,
		screen.View(),
		parent,
	)
	if err != nil {
		m.removeInstance(instance)
		return nil, err
	}

	instance.Native = native

	if options.Mode == WindowModal {
		if err := native.SetModal(true); err != nil {
			_ = native.Close()
			m.removeInstance(instance)
			return nil, err
		}
	}

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		_ = native.Close()
		m.removeInstance(instance)
		return nil, ErrManagerClosed
	}

	if options.Parent != "" {
		if parentInstance := m.windows[options.Parent]; parentInstance != nil {
			parentInstance.AddChild(instance.ID)
		}
	}
	m.mu.Unlock()

	if err := native.Show(); err != nil {
		_ = m.Close(instance.ID)
		return nil, err
	}

	return instance, nil
}

func (m *Manager) Get(id WindowID) (*Instance, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	instance, ok := m.windows[id]
	return instance, ok
}

func (m *Manager) Instances() []*Instance {
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

	delete(m.windows, instance.ID)

	if instance.ParentID != "" {
		if parent := m.windows[instance.ParentID]; parent != nil {
			parent.RemoveChild(instance.ID)
		}
	}
}

func (m *Manager) Close(id WindowID) error {
	m.mu.RLock()
	instance := m.windows[id]
	m.mu.RUnlock()

	if instance == nil {
		return ErrWindowNotFound
	}

	if !instance.markClosed() {
		return ErrWindowClosed
	}

	var closeErr error

	// Close children before closing their parent.
	for _, childID := range instance.Children() {
		if err := m.Close(childID); err != nil && closeErr == nil {
			closeErr = err
		}
	}

	if instance.Native != nil {
		if err := instance.Native.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}

	instance.Results.Complete(nil, ErrWindowClosed)
	m.removeInstance(instance)

	return closeErr
}

func (m *Manager) CloseAll() error {
	m.mu.Lock()
	m.closed = true

	var roots []WindowID
	for id, instance := range m.windows {
		if instance.ParentID == "" {
			roots = append(roots, id)
		}
	}
	m.mu.Unlock()

	var firstErr error

	for _, id := range roots {
		if err := m.Close(id); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
