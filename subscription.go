package rx

import (
	"errors"
	"sync"
)

// Subscription is an interface that allows monitoring and controlling a subscription.
// It provides methods for tracking the subscription's lifecycle.
type Subscription interface {
	// Subscribed returns true until Unsubscribe is called.
	Subscribed() bool

	// Unsubscribe will change the state to unsubscribed.
	Unsubscribe()

	// Done returns a channel that is closed when the subscription state changes to unsubscribed.
	// This channel can be used with select statements to react to subscription termination events.
	// If the scheduler is not concurrent, it will spawn a goroutine to wait for the scheduler.
	Done() <-chan struct{}

	// Err returns the subscription's terminal state:
	// - nil if the observable completed successfully
	// - the observable's error if it terminated with an error
	// - SubscriptionCanceled if the subscription was manually unsubscribed
	// - SubscriptionActive if the subscription is still active
	Err() error

	// Wait blocks until the subscription state becomes unsubscribed.
	// If the subscription is already unsubscribed, it returns immediately.
	// If the scheduler is not concurrent, it will wait for the scheduler to complete.
	// Returns:
	// - nil if the observable completed successfully
	// - the observable's error if it terminated with an error
	// - SubscriptionCanceled if the subscription was manually unsubscribed
	Wait() error
}

type subscription struct {
	subscriber
	scheduler Scheduler
	err       error
}

// ErrSubscriptionActive is the error returned by Err() when the
// subscription is still active and has not yet completed or been canceled.
var ErrSubscriptionActive = errors.Join(Err, errors.New("subscription active"))

// ErrSubscriptionCanceled is the error returned by Wait() and Err() when the
// subscription was canceled by calling Unsubscribe() on the Subscription.
// This indicates the subscription was terminated by the subscriber rather than
// by the observable completing normally or with an error.
var ErrSubscriptionCanceled = errors.Join(Err, errors.New("subscription canceled"))

func newSubscription(scheduler Scheduler) *subscription {
	return &subscription{scheduler: scheduler, err: ErrSubscriptionCanceled}
}

func (s *subscription) done(err error) {
	s.Lock()
	s.err = err
	s.Unlock()

	s.Unsubscribe()
}

func (s *subscription) Done() <-chan struct{} {
	cancel := make(chan struct{})
	s.OnUnsubscribe(func() { close(cancel) })
	if !s.scheduler.IsConcurrent() {
		go s.scheduler.Wait()
	}
	return cancel
}

func (s *subscription) Err() (err error) {
	if s.Subscribed() {
		return ErrSubscriptionActive
	}
	s.Lock()
	err = s.err
	s.Unlock()
	return
}

func (s *subscription) Wait() error {
	if !s.scheduler.IsConcurrent() {
		s.scheduler.Wait()
	}

	if s.Subscribed() {
		var wg sync.WaitGroup
		wg.Add(1)
		s.OnUnsubscribe(wg.Done)
		wg.Wait()
	}

	return s.Err()
}
