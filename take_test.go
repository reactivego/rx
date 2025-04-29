package rx_test

import (
	"testing"
	"time"

	"github.com/reactivego/rx"
)

func TestTake(t *testing.T) {
	t.Run("Basic take behavior", func(t *testing.T) {
		nums := rx.From(1, 2, 3, 4, 5)
		taken := nums.Take(3)

		var results []int
		taken.Collect(&results).Wait()

		expected := []int{1, 2, 3}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}

		for i, v := range expected {
			if results[i] != v {
				t.Errorf("Expected %v at index %v, got %v", v, i, results[i])
			}
		}
	})

	t.Run("Take more than available", func(t *testing.T) {
		nums := rx.From(1, 2, 3)
		taken := nums.Take(5)

		var results []int
		taken.Collect(&results).Wait()

		expected := []int{1, 2, 3}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}
	})

	t.Run("Take zero", func(t *testing.T) {
		nums := rx.From(1, 2, 3)
		taken := nums.Take(0)

		var results []int
		taken.Collect(&results).Wait()

		if len(results) != 0 {
			t.Errorf("Expected 0 items, got %v", len(results))
		}
	})

	t.Run("Take from channel", func(t *testing.T) {
		ch := make(chan int, 3)
		observable := rx.Recv(ch)

		completed := false
		var results []int

		// Use goroutine scheduler for async processing
		s := observable.Take(1).Subscribe(func(next int, err error, done bool) {
			if !done {
				results = append(results, next)
			} else if err == nil {
				completed = true
			}
		}, rx.Goroutine)

		// Send a value through the channel
		ch <- 1

		// Wait for subscription to complete
		err := s.Wait()

		// Verify results
		if err != nil {
			t.Errorf("Expected completion, got %v", err)
		}

		if !completed {
			t.Error("Observable didn't complete")
		}

		if len(results) != 1 || results[0] != 1 {
			t.Errorf("Expected [1], got %v", results)
		}
	})

	t.Run("Take zero from channel", func(t *testing.T) {
		ch := make(chan int)
		observable := rx.Recv(ch)
		var results []int

		err := observable.Take(0).Collect(&results).Wait()

		if err != nil {
			t.Errorf("Expected completion, got %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 items, got %v", len(results))
		}
	})

	t.Run("Take with multiple values from channel", func(t *testing.T) {
		ch := make(chan int, 5)
		observable := rx.Recv(ch)

		var results []int

		// Take only 3 values
		s := observable.Take(3).Subscribe(func(next int, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, rx.Goroutine)

		// Send values
		ch <- 10
		ch <- 20
		ch <- 30
		ch <- 40 // This should not be received due to Take(3)
		ch <- 50 // This should not be received due to Take(3)

		// Allow some time for processing
		time.Sleep(10 * time.Millisecond)

		// Wait for subscription to complete
		s.Wait()

		// Verify results
		expected := []int{10, 20, 30}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}

		for i, v := range expected {
			if results[i] != v {
				t.Errorf("Expected %v at index %v, got %v", v, i, results[i])
			}
		}
	})

	t.Run("Take with closed channel", func(t *testing.T) {
		ch := make(chan int, 1)
		observable := rx.Recv(ch)

		completed := false
		var results []int

		s := observable.Take(3).Subscribe(func(next int, err error, done bool) {
			if !done {
				results = append(results, next)
			} else if err == nil {
				completed = true
			}
		}, rx.Goroutine)

		// Send one value and close the channel
		ch <- 100
		close(ch)

		// Wait for subscription to complete
		s.Wait()

		// Verify results
		if !completed {
			t.Error("Observable didn't complete")
		}

		if len(results) != 1 || results[0] != 100 {
			t.Errorf("Expected [100], got %v", results)
		}
	})

	t.Run("Take with Go method", func(t *testing.T) {
		ch := make(chan int, 1)
		observable := rx.Recv(ch)

		completed := false
		var results []int

		// Use Go method as in the provided main function
		s := observable.Take(1).Subscribe(func(next int, err error, done bool) {
			if !done {
				results = append(results, next)
			} else if err == nil {
				completed = true
			}
		}, rx.Goroutine)

		// Send a value through the channel
		ch <- 1

		// Wait for subscription to complete
		s.Wait()

		// Verify results
		if !completed {
			t.Error("Observable didn't complete")
		}

		if len(results) != 1 || results[0] != 1 {
			t.Errorf("Expected [1], got %v", results)
		}
	})
}
