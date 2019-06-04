package subscriber

import (
	"sync"
	"sync/atomic"
)

const (
	subscribed = iota
	unsubscribed
)

type subscription struct {
	state int32

	sync.Mutex
	callbacks []func()
	onWait    func()
}

func (s *subscription) Unsubscribe() {
	if atomic.CompareAndSwapInt32(&s.state, subscribed, unsubscribed) {
		s.Lock()
		for _, cb := range s.callbacks {
			cb()
		}
		s.callbacks = nil
		s.Unlock()
	}
}

func (s *subscription) Closed() bool {
	return atomic.LoadInt32(&s.state) != subscribed
}

func (s *subscription) Canceled() bool {
	return atomic.LoadInt32(&s.state) != subscribed
}

func (s *subscription) Wait() {
	s.Lock()
	wait := s.onWait
	s.Unlock()
	if wait != nil {
		wait()
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		s.Add(wg.Done)
		wg.Wait()		
	}
}

func (s *subscription) Add(callback func()) Subscriber {
	child := &subscription{callbacks: []func(){callback}}
	s.Lock()
	if atomic.LoadInt32(&s.state) != subscribed {
		child.Unsubscribe()
	} else {
		s.callbacks = append(s.callbacks, child.Unsubscribe)
	}
	s.Unlock()
	return child
}

func (s *subscription) AddChild() Subscriber {
	child := &subscription{}
	s.Lock()
	if atomic.LoadInt32(&s.state) != subscribed {
		child.Unsubscribe()
	} else {
		s.callbacks = append(s.callbacks, child.Unsubscribe)
	}
	s.Unlock()
	return child
}

func (s *subscription) OnUnsubscribe(callback func()) {
	if callback == nil {
		return
	}
	s.Lock()
	if atomic.LoadInt32(&s.state) == subscribed {
		s.callbacks = append(s.callbacks, callback)
	}
	s.Unlock()
}

func (s *subscription) OnWait(callback func()) {
	s.Lock()
	s.onWait = callback
	s.Unlock()
}
