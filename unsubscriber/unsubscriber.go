package unsubscriber

import (
	"sync"
	"sync/atomic"
)

type Unsubscriber interface {
	Unsubscribe()
	Unsubscribed() bool
}

////////////////////////////////////////////////////////
// Int32
////////////////////////////////////////////////////////

// Int32 is an Int32 value with atomic operations implementing the Unsubscriber interface
type Int32 int32

func (t *Int32) Unsubscribe() {
	atomic.StoreInt32((*int32)(t), 1)
}

func (t *Int32) Unsubscribed() bool {
	return atomic.LoadInt32((*int32)(t)) == 1
}

////////////////////////////////////////////////////////
// Collection
////////////////////////////////////////////////////////

// Collection is itself an Unsubscriber that will forward an Unsubscribe method
// call to all unsubscribers in the collection. Useful for managing a collection of
// concurrently operating subscribers in e.g. a merge operation.
type Collection struct {
	Int32
	sync.Mutex
	unsubscribers []Unsubscriber
	onUnsubscribe func()
}

// Add will append the unsubscriber to the list of unsubscribers already present.
func (c *Collection) Add(u Unsubscriber) bool {
	c.Lock()
	defer c.Unlock()
	if c.Unsubscribed() {
		u.Unsubscribe()
		return false
	}
	c.unsubscribers = append(c.unsubscribers, u)
	return true
}

// Set will replace the collection of unsubscribers with a single unsubscriber.
func (c *Collection) Set(u Unsubscriber) bool {
	c.Lock()
	defer c.Unlock()
	if c.Unsubscribed() {
		u.Unsubscribe()
		return false
	}
	c.unsubscribers = []Unsubscriber{u}
	return true
}

func (c *Collection) Unsubscribe() {
	c.Lock()
	defer c.Unlock()
	if c.Unsubscribed() {
		return
	}
	c.Int32.Unsubscribe()
	for _, unsubscriber := range c.unsubscribers {
		if !unsubscriber.Unsubscribed() {
			unsubscriber.Unsubscribe()
		}
	}
	if c.onUnsubscribe != nil {
		c.onUnsubscribe()
	}
}

func (c *Collection) OnUnsubscribe(f func()) {
	c.Lock()
	defer c.Unlock()
	c.onUnsubscribe = f
}
