package windows

import (
	"sync"
	"sync/atomic"
)

var windowSequence atomic.Uint64

func nextWindowSequence() uint64 { return windowSequence.Add(1) }

type instanceState uint8

const (
	instanceCreating instanceState = iota
	instanceOpen
	instanceClosing
	instanceClosed
)

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
	state    instanceState
}

func newInstance(screenID ScreenID, options WindowOptions) *Instance {
	return &Instance{
		ID:       newWindowID(),
		ScreenID: screenID,
		ParentID: options.Parent,
		Options:  options,
		Results:  newResultFuture(),
		children: make(map[WindowID]struct{}),
		state:    instanceCreating,
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
	return i.state == instanceClosing || i.state == instanceClosed
}

func (i *Instance) setScreen(screen Screen) {
	i.mu.Lock()
	i.Screen = screen
	i.mu.Unlock()
}

func (i *Instance) setNative(native NativeWindow) {
	i.mu.Lock()
	i.Native = native
	i.mu.Unlock()
}

func (i *Instance) nativeWindow() NativeWindow {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Native
}

func (i *Instance) markOpen() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.state != instanceCreating {
		return false
	}
	i.state = instanceOpen
	return true
}

func (i *Instance) markClosing() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.state == instanceClosing || i.state == instanceClosed {
		return false
	}
	i.state = instanceClosing
	return true
}

func (i *Instance) markClosed() {
	i.mu.Lock()
	i.state = instanceClosed
	i.mu.Unlock()
}

// Complete sends a result to every caller waiting on this instance.
// The first completion wins.
func (i *Instance) Complete(value any, err error) {
	i.Results.Complete(value, err)
}
