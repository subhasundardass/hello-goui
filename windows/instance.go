package windows

import (
	"sync"
	"sync/atomic"
)

var windowSequence atomic.Uint64

func nextWindowSequence() uint64 {
	return windowSequence.Add(1)
}

type Instance struct {
	ID       WindowID
	ScreenID ScreenID
	ParentID WindowID
	Options  WindowOptions

	Screen  Screen
	Native  NativeWindow
	Results *ResultFuture

	mu       sync.RWMutex
	children map[WindowID]struct{}
	closed   bool
}

func newInstance(
	screenID ScreenID,
	options WindowOptions,
) *Instance {
	id := newWindowID()

	return &Instance{
		ID:       id,
		ScreenID: screenID,
		ParentID: options.Parent,
		Options:  options,
		Results:  newResultFuture(),
		children: make(map[WindowID]struct{}),
	}
}

func (i *Instance) AddChild(id WindowID) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.children[id] = struct{}{}
}

func (i *Instance) RemoveChild(id WindowID) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.children, id)
}

func (i *Instance) Children() []WindowID {
	i.mu.RLock()
	defer i.mu.RUnlock()

	children := make([]WindowID, 0, len(i.children))
	for id := range i.children {
		children = append(children, id)
	}

	return children
}

func (i *Instance) IsClosed() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.closed
}

func (i *Instance) markClosed() bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.closed {
		return false
	}

	i.closed = true
	return true
}

// Complete sends a result back to the caller.
func (i *Instance) Complete(value any, err error) {
	i.Results.Complete(value, err)
}
