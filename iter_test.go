package rx_test

import (
	"fmt"
	"testing"

	"github.com/reactivego/rx"
)

// Additional integration test that combines multiple iterator functions
func TestIter(t *testing.T) {
	t.Run("Pull to Values to Pull round trip", func(t *testing.T) {
		// Create a sequence
		seq := func(yield func(int) bool) {
			for i := 1; i <= 5; i++ {
				if !yield(i * 10) {
					break
				}
			}
		}

		// Convert to Observable, then back to an iterator, and then back to Observable
		obs1 := rx.Pull(seq)
		iter := obs1.Values()
		obs2 := rx.Pull(iter)

		// Check results
		var results []int
		err := obs2.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
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

	t.Run("All iterator with transformation", func(t *testing.T) {
		nums := rx.From("apple", "banana", "cherry")

		// Use the iterator to create a formatted string
		var formattedItems []string
		for i, v := range nums.All() {
			formattedItems = append(formattedItems, fmt.Sprintf("%d: %s", i, v))
		}

		expected := []string{
			"0: apple",
			"1: banana",
			"2: cherry",
		}

		if len(formattedItems) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(formattedItems))
		}

		for i, v := range expected {
			if formattedItems[i] != v {
				t.Errorf("Expected %v at index %v, got %v", v, i, formattedItems[i])
			}
		}
	})
}
