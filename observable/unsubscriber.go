package observable

import (
	"sync"
	"sync/atomic"
)

type Unsubscriber interface {
	Unsubscribe()
	Unsubscribed() bool
}

// UnsubscriberEvents provides lifecycle event callbacks for an Unsubscriber.
type UnsubscriberEvents interface {
	OnUnsubscribe(func())
}

////////////////////////////////////////////////////////
// UnsubscriberInt
////////////////////////////////////////////////////////

// UnsubscriberInt is an Int32 value with atomic operations implementing the Unsubscriber interface

type UnsubscriberInt int32

func (t *UnsubscriberInt) Unsubscribe() {
	atomic.StoreInt32((*int32)(t), 1)
}

func (t *UnsubscriberInt) Unsubscribed() bool {
	return atomic.LoadInt32((*int32)(t)) == 1
}

////////////////////////////////////////////////////////
// UnsubscribedUnsubscriber
////////////////////////////////////////////////////////

// UnsubscribedUnsubscriber is a unsubscriber that always reports that it is already unsubscribed.
type UnsubscribedUnsubscriber struct{}

func (UnsubscribedUnsubscriber) Unsubscribe() {
}

func (UnsubscribedUnsubscriber) Unsubscribed() bool {
	return true
}

////////////////////////////////////////////////////////
// UnsubscriberCollection
////////////////////////////////////////////////////////

// UnsubscriberCollection is itself an Unsubscriber that will forward an Unsubscribe method
// call to all unsubscribers in the collection. Useful for managing a collection of
// concurrently operating subscribers in e.g. a merge operation.
type UnsubscriberCollection struct {
	UnsubscriberInt
	lock          sync.Mutex
	unsubscribers []Unsubscriber
}

// Add will append the unsubscriber to the list of unsubscribers already present.
func (c *UnsubscriberCollection) Add(u Unsubscriber) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Unsubscribed() {
		u.Unsubscribe()
		return false
	}
	c.unsubscribers = append(c.unsubscribers, u)
	return true
}

// Set will replace the collection of unsubscribers with a single unsubscriber.
func (c *UnsubscriberCollection) Set(u Unsubscriber) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Unsubscribed() {
		u.Unsubscribe()
		return false
	}
	c.unsubscribers = []Unsubscriber{u}
	return true
}

func (c *UnsubscriberCollection) Unsubscribe() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Unsubscribed() {
		return
	}
	c.UnsubscriberInt.Unsubscribe()
	for _, unsubscriber := range c.unsubscribers {
		if !unsubscriber.Unsubscribed() {
			unsubscriber.Unsubscribe()
		}
	}
}

////////////////////////////////////////////////////////
// UnsubscriberChannel
////////////////////////////////////////////////////////

// UnsubscriberChannel is implemented with a channel which is closed when unsubscribed.
type UnsubscriberChannel chan struct{}

func NewUnsubscriberChannel() UnsubscriberChannel {
	return make(UnsubscriberChannel)
}

func (c UnsubscriberChannel) Unsubscribe() {
	defer recover()
	close(c)
}

func (c UnsubscriberChannel) Unsubscribed() bool {
	select {
	case _, ok := <-c:
		return !ok
	default:
		return false
	}
}

func (c UnsubscriberChannel) OnUnsubscribe(handler func()) {
	go func() {
		<-c
		handler()
	}()
}

////////////////////////////////////////////////////////
// UnsubscriberFunc
////////////////////////////////////////////////////////

// UnsubscriberFunc is implemented as a pointer to a callback function. Whenever Unsubscribe()
// is called the callback is invoked and then the pointer to the callback is set to nil.
// The method Unsubscribed() returns true when the pointer to the callback function has been
// set to nil.
type UnsubscriberFunc func()

func (c *UnsubscriberFunc) Unsubscribe() {
	if *c != nil {
		(*c)()
		*c = nil
	}
}

func (c *UnsubscriberFunc) Unsubscribed() bool {
	return *c == nil
}
