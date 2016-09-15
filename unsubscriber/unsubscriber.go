package unsubscriber

import (
	"sync"
	"sync/atomic"
)

type Unsubscriber interface {
	Unsubscribe()
	Unsubscribed() bool
	OnUnsubscribe(func())
	AddChild() Unsubscriber
	Wait()
}

////////////////////////////////////////////////////////
// Standard
////////////////////////////////////////////////////////
// Standard is the standard unsubscriber implementation
// that can also call back to a function when unsubscribed.
type Standard struct {
	int32
	onUnsubscribe []func()
	sync.Mutex
}

func (r *Standard) Unsubscribed() bool {
	return atomic.LoadInt32((*int32)(&r.int32)) == 1
}

func (r *Standard) Unsubscribe() {
	r.Lock()
	defer r.Unlock()
	if r.Unsubscribed() {
		return
	}
	atomic.StoreInt32((*int32)(&r.int32), 1)
	for _, f := range r.onUnsubscribe {
		f()
	}
}

func (r *Standard) OnUnsubscribe(f func()) {
	r.Lock()
	defer r.Unlock()
	if r.Unsubscribed() {
		f()
	} else {
		r.onUnsubscribe = append(r.onUnsubscribe, f)
	}
}

// AddChild will create an Unsubscriber who's Unsubscribe
// method will be called when the parent's (the current object)
// Unsubscribe method is called. Calling the Unsubscribe method
// on the child will NOT propagate to the parent!
func (r *Standard) AddChild() Unsubscriber {
	r.Lock()
	defer r.Unlock()
	child := &Standard{}
	if r.Unsubscribed() {
		child.Unsubscribe()
	} else {
		r.onUnsubscribe = append(r.onUnsubscribe, child.Unsubscribe)
	}
	return child
}

func (r *Standard) Wait() {
	done := make(chan struct{})
	r.OnUnsubscribe(func() {
		close(done)
	})
	<-done
}

func New() Unsubscriber {
	return &Standard{}
}
