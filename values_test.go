package rx_test

import (
	"context"
	"testing"

	"github.com/reactivego/rx"
)

func TestValues(t *testing.T) {
	t.Run("Basic values extraction", func(t *testing.T) {
		nums := rx.From(1, 2, 3, 4, 5)
		var results []int
		for v := range nums.Values(context.Background()) {
			results = append(results, v)
		}
		expected := []int{1, 2, 3, 4, 5}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}

		for i, v := range expected {
			if results[i] != v {
				t.Errorf("Expected %v at index %v, got %v", v, i, results[i])
			}
		}
	})

	t.Run("Empty observable", func(t *testing.T) {
		nums := rx.Empty[int]()
		count := 0
		for range nums.Values(context.Background()) {
			count++
		}
		if count != 0 {
			t.Errorf("Expected 0 items, got %v", count)
		}
	})

	t.Run("Values with filtering", func(t *testing.T) {
		nums := rx.From(1, 2, 3, 4, 5, 6, 7)
		filtered := nums.Filter(func(i int) bool {
			return i%2 == 0
		})
		var results []int
		for v := range filtered.Values(context.Background()) {
			results = append(results, v)
		}
		expected := []int{2, 4, 6}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}

		for i, v := range expected {
			if results[i] != v {
				t.Errorf("Expected %v at index %v, got %v", v, i, results[i])
			}
		}
	})

	t.Run("Values with early termination", func(t *testing.T) {
		nums := rx.From(1, 2, 3, 4, 5)
		taken := nums.Take(3)
		var results []int
		for v := range taken.Values(context.Background()) {
			results = append(results, v)
		}
		expected := []int{1, 2, 3}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}
	})

	t.Run("Pull to Values round trip", func(t *testing.T) {
		// Create a sequence
		seq := func(yield func(int) bool) {
			for i := 1; i <= 5; i++ {
				if !yield(i * 10) {
					break
				}
			}
		}

		// Convert to Observable, then back to an iterator
		obs := rx.Pull(seq)

		// Check results by iterating with Values()
		var results []int
		for v, err := range obs.Values(context.Background()) {
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			results = append(results, v)
		}

		expected := []int{10, 20, 30, 40, 50}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}

		for i, v := range expected {
			if results[i] != v {
				t.Errorf("Expected %v at index %v, got %v", v, i, results[i])
			}
		}
	})
}
