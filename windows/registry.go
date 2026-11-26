package windows

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu        sync.RWMutex
	factories map[ScreenID]ScreenFactory
}

func NewRegistry() *Registry {
	return &Registry{factories: make(map[ScreenID]ScreenFactory)}
}

func (r *Registry) Register(id ScreenID, factory ScreenFactory) error {
	if id == "" {
		return fmt.Errorf("screen ID cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("screen factory cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[id]; exists {
		return fmt.Errorf("screen %q already registered", id)
	}
	r.factories[id] = factory
	return nil
}

func (r *Registry) Get(id ScreenID) (ScreenFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrScreenNotFound, id)
	}
	return factory, nil
}

func (r *Registry) MustRegister(id ScreenID, factory ScreenFactory) {
	if err := r.Register(id, factory); err != nil {
		panic(err)
	}
}
