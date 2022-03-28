package observable

import (
	"sync"
)

// Subscription is an interface that allows code to monitor and control a
// subscription it received.
type Subscription interface {
	Subscribable

	// Wait will by default block the calling goroutine and wait for the
	// Unsubscribe method to be called on this subscription.
	// However, when OnWait was called with a callback wait function it will
	// call that instead. Calling Wait on a subscription that has already been
	// canceled will return immediately. If the subscriber was canceled by
	// calling Unsubscribe, then the error returned is ErrUnsubscribed.
	// If the subscriber was terminated by calling Done, then the error
	// returned here is the one passed to Done.
	Wait() error
}

type subscription struct {
	subscriber
	wait func()
	err  error
}

// SubscriptionCanceled is the error returned by Wait when the Unsubscribe method is
// called on an active subscription.
const SubscriptionCanceled = Error("subscription canceled")

func newSubscription(wait func()) *subscription {
	return &subscription{wait: wait, err: SubscriptionCanceled}
}

// done will set the error internally and then cancel the subscription by
// calling the Unsubscribe method. A nil value for error indicates success.
func (s *subscription) done(err error) {
	s.Lock()
	s.err = err
	s.Unlock()
	s.Unsubscribe()
}

func (s *subscription) Wait() error {
	if wait := s.wait; wait != nil {
		wait()
	}
	if s.Subscribed() {
		var wg sync.WaitGroup
		wg.Add(1)
		s.OnUnsubscribe(wg.Done)
		wg.Wait()
	}
	s.Lock()
	err := s.err
	s.Unlock()
	return err
}
