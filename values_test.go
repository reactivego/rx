package rx_test

import (
	"testing"

	"github.com/reactivego/rx"
)

func TestValues(t *testing.T) {
	t.Run("Basic values extraction", func(t *testing.T) {
		nums := rx.From(1, 2, 3, 4, 5)
		var results []int
		for v := range nums.Values() {
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
		for range nums.Values() {
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
		for v := range filtered.Values() {
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
		for v := range taken.Values() {
			results = append(results, v)
		}
		expected := []int{1, 2, 3}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}
	})
}
