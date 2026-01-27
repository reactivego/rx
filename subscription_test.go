package rx_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/reactivego/rx"
)

func TestSubscription(t *testing.T) {
	const msec = time.Millisecond

	t.Run("Context timeout before subscription completes", func(t *testing.T) {
		// Force the maybeLongOperation to always take a long time
		maybeLongOperation := func(i int) rx.Observable[int] {
			// Ensure it exceeds our timeout
			return rx.Of(i).Delay(2000 * msec)
		}

		isOne := func(i int) bool { return i == 1 }

		ctx, cancel := context.WithTimeout(context.Background(), 500*msec)
		defer cancel()

		s := rx.ConcatMap(rx.Of(1).Filter(isOne).Take(1), maybeLongOperation).Go(context.Background())
		defer s.Unsubscribe()

		select {
		case <-ctx.Done():
			// Success, timeout exceeded!
		case <-s.Done():
			t.Errorf("expected context timeout, but observable completed with error: %v", s.Err())
		}
	})

	t.Run("Subscription completes before context timeout", func(t *testing.T) {
		// Force the maybeLongOperation to always be quick
		maybeLongOperation := func(i int) rx.Observable[int] {
			// No sleep, complete immediately
			return rx.Empty[int]()
		}

		isOne := func(i int) bool { return i == 1 }

		ctx, cancel := context.WithTimeout(context.Background(), 2000*msec)
		defer cancel()

		s := rx.ConcatMap(rx.From(1).Filter(isOne).Take(1), maybeLongOperation).Go(context.Background())
		defer s.Unsubscribe()

		select {
		case <-ctx.Done():
			t.Error("expected success, got timeout")
		case <-s.Done():
			if err := s.Err(); err != nil {
				t.Errorf("expected nil error for successful completion, got: %v", err)
			}
		}
	})

	t.Run("Unsubscribe terminates the subscription", func(t *testing.T) {
		// Create a subscription that would take a long time
		s := rx.Of(1).Delay(600 * msec).Go(context.Background())

		// Immediately unsubscribe
		s.Unsubscribe()

		// Check that the Done channel closes quickly
		select {
		case <-s.Done():
			// Success, the subscription completed after unsubscribing
			if err := s.Err(); err != rx.ErrSubscriptionCanceled {
				t.Errorf("expected %v error, got: %v", rx.ErrSubscriptionCanceled, err)
			}
		case <-time.After(500 * msec):
			t.Error("subscription did not complete in time after unsubscribing")
		}
	})

	t.Run("Done channel closes after normal completion", func(t *testing.T) {
		s := rx.From(1, 2, 3).Take(3).Go(context.Background())

		// Set a timeout to detect if Done channel doesn't close by itself.
		select {
		case <-s.Done():
			// Success - the channel closed as expected
			if err := s.Err(); err != nil {
				t.Errorf("expected nil error for successful completion, got: %v", err)
			}
		case <-time.After(500 * msec):
			t.Error("subscription Done channel did not close within the expected timeframe")
		}
	})

	t.Run("Wait blocks until completion", func(t *testing.T) {
		s := rx.Of(1).Delay(350 * msec).Go(context.Background())

		// Wait should block until the observable completes
		start := time.Now()
		err := s.Wait()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}

		if elapsed < 300*msec {
			t.Errorf("wait returned too quickly, expected at least 300ms, got: %v", elapsed)
		}
	})

	t.Run("Wait returns immediately if already completed", func(t *testing.T) {
		s := rx.Of(1).Go(context.Background())

		// Make sure it's done
		<-s.Done()

		// Wait should return immediately
		start := time.Now()
		err := s.Wait()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}

		if elapsed > 50*msec {
			t.Errorf("wait took too long for completed subscription, got: %v", elapsed)
		}
	})

	t.Run("Err returns error from observable", func(t *testing.T) {
		expectedErr := errors.New("test error")
		s := rx.Throw[int](expectedErr).Go(context.Background())

		// Wait for completion
		<-s.Done()

		if err := s.Err(); err != expectedErr {
			t.Errorf("expected error %v, got: %v", expectedErr, err)
		}
	})

	t.Run("Err returns SubscriptionActive while active", func(t *testing.T) {
		ch := make(chan struct{})
		s := rx.Recv(ch).Go(context.Background())

		// Subscription should be active
		if err := s.Err(); err != rx.ErrSubscriptionActive {
			t.Errorf("expected SubscriptionActive, got: %v", err)
		}

		// Complete the observable
		close(ch)
		<-s.Done()
	})

	t.Run("Context cancellation terminates subscription", func(t *testing.T) {
		// Create a context that will be canceled
		ctx, cancel := context.WithCancel(context.Background())

		// Create a long-running subscription with the cancelable context
		ch := make(chan int)
		s := rx.Recv(ch).Go(ctx)

		// Cancel the context
		cancel()

		// The subscription should terminate due to context cancellation
		select {
		case <-s.Done():
			// Success - subscription terminated due to context cancellation
			if err := s.Err(); err != rx.ErrSubscriptionCanceled {
				t.Errorf("expected ErrSubscriptionCanceled, got: %v", err)
			}
		case <-time.After(500 * msec):
			t.Error("subscription did not terminate after context cancellation")
		}
	})

	t.Run("Unsubscribe cancels subscriber context", func(t *testing.T) {
		// Channel to capture the subscriber's context
		ctxChan := make(chan context.Context, 1)

		// Create an observable that captures its subscriber's context and runs async
		observable := func(observe rx.Observer[int], scheduler rx.Scheduler, subscriber rx.Subscriber) {
			ctxChan <- subscriber.Context()
			// Run the blocking work in a goroutine so Subscribe() returns
			go func() {
				<-subscriber.Context().Done()
				var zero int
				observe(zero, nil, true)
			}()
		}

		s := rx.Observable[int](observable).Go(context.Background())

		// Get the subscriber's context
		subscriberCtx := <-ctxChan

		// Context should not be done yet
		select {
		case <-subscriberCtx.Done():
			t.Error("subscriber context should not be done before unsubscribe")
		default:
			// Good, context is not done
		}

		// Unsubscribe
		s.Unsubscribe()

		// Now the subscriber's context should be canceled
		select {
		case <-subscriberCtx.Done():
			// Success - context was canceled by unsubscribe
		case <-time.After(500 * msec):
			t.Error("subscriber context was not canceled after unsubscribe")
		}
	})
}
