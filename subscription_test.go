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

		s := rx.ConcatMap(rx.Of(1).Filter(isOne).Take(1), maybeLongOperation).Go()
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

		s := rx.ConcatMap(rx.From(1).Filter(isOne).Take(1), maybeLongOperation).Go()
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
		s := rx.Of(1).Delay(600 * msec).Go()

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
		s := rx.From(1, 2, 3).Take(3).Go()

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
		s := rx.Of(1).Delay(350 * msec).Go()

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
		s := rx.Of(1).Go()

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
		s := rx.Throw[int](expectedErr).Go()

		// Wait for completion
		<-s.Done()

		if err := s.Err(); err != expectedErr {
			t.Errorf("expected error %v, got: %v", expectedErr, err)
		}
	})

	t.Run("Err returns SubscriptionActive while active", func(t *testing.T) {
		ch := make(chan struct{})
		s := rx.Recv(ch).Go()

		// Subscription should be active
		if err := s.Err(); err != rx.ErrSubscriptionActive {
			t.Errorf("expected SubscriptionActive, got: %v", err)
		}

		// Complete the observable
		close(ch)
		<-s.Done()
	})
}
