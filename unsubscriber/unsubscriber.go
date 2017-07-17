package unsubscriber

import (
	"sync"
	"sync/atomic"
)

// Unsubscriber is an interface returned by a subscribe function
// that allows you to unsubscribe from the subscription that was
// created when subscribe was called.
type Unsubscriber interface {
	Unsubscribe()
	Unsubscribed() bool
	OnUnsubscribe(func())
	AddChild() Unsubscriber
	Wait()
}

////////////////////////////////////////////////////////
// unsubscriber
////////////////////////////////////////////////////////
// unsubscriber is the standard unsubscriber implementation
type unsubscriber struct {
	int32
	onUnsubscribe []func()
	sync.Mutex
}

func (r *unsubscriber) Unsubscribed() bool {
	return atomic.LoadInt32((*int32)(&r.int32)) == 1
}

func (r *unsubscriber) Unsubscribe() {
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

func (r *unsubscriber) OnUnsubscribe(f func()) {
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
func (r *unsubscriber) AddChild() Unsubscriber {
	r.Lock()
	defer r.Unlock()
	child := &unsubscriber{}
	if r.Unsubscribed() {
		child.Unsubscribe()
	} else {
		r.onUnsubscribe = append(r.onUnsubscribe, child.Unsubscribe)
	}
	return child
}

func (r *unsubscriber) Wait() {
	done := make(chan struct{})
	r.OnUnsubscribe(func() {
		close(done)
	})
	<-done
}

func New() Unsubscriber {
	return &unsubscriber{}
}
