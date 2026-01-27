package rx_test

import (
	"fmt"
	"testing"

	"github.com/reactivego/rx"
)

func TestAll(t *testing.T) {
	t.Run("Basic indexed values", func(t *testing.T) {
		nums := rx.From("a", "b", "c")
		var indexes []int
		var values []string
		for next, err := range nums.All() {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			indexes = append(indexes, next.First)
			values = append(values, next.Second)
		}
		expectedIndexes := []int{0, 1, 2}
		expectedValues := []string{"a", "b", "c"}

		if len(indexes) != len(expectedIndexes) {
			t.Errorf("Expected %v indexes, got %v", len(expectedIndexes), len(indexes))
		}

		if len(values) != len(expectedValues) {
			t.Errorf("Expected %v values, got %v", len(expectedValues), len(values))
		}

		for i, v := range expectedIndexes {
			if indexes[i] != v {
				t.Errorf("Expected index %v at position %v, got %v", v, i, indexes[i])
			}
		}

		for i, v := range expectedValues {
			if values[i] != v {
				t.Errorf("Expected value %v at position %v, got %v", v, i, values[i])
			}
		}
	})

	t.Run("Empty observable", func(t *testing.T) {
		nums := rx.Empty[string]()
		count := 0
		for _, err := range nums.All() {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			count++
		}
		if count != 0 {
			t.Errorf("Expected 0 items, got %v", count)
		}
	})

	t.Run("All with filtering", func(t *testing.T) {
		nums := rx.From(10, 20, 30, 40, 50)
		filtered := nums.Filter(func(i int) bool {
			return i > 25
		})
		// Collect indexes and values using direct iteration
		var indexes []int
		var values []int
		seq := filtered.All()
		seq(func(next rx.Tuple2[int, int], err error) bool {
			indexes = append(indexes, next.First)
			values = append(values, next.Second)
			return true
		})
		// Note: indexes should be 0, 1, 2 (not 2, 3, 4) because All reindexes
		expectedIndexes := []int{0, 1, 2}
		expectedValues := []int{30, 40, 50}

		if len(indexes) != len(expectedIndexes) {
			t.Errorf("Expected %v indexes, got %v", len(expectedIndexes), len(indexes))
		}

		for i, v := range expectedIndexes {
			if indexes[i] != v {
				t.Errorf("Expected index %v at position %v, got %v", v, i, indexes[i])
			}
		}

		for i, v := range expectedValues {
			if values[i] != v {
				t.Errorf("Expected value %v at position %v, got %v", v, i, values[i])
			}
		}
	})

	t.Run("All iterator with transformation", func(t *testing.T) {
		nums := rx.From("apple", "banana", "cherry")

		// Use the iterator to create a formatted string
		var formattedItems []string
		for next, err := range nums.All() {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			formattedItems = append(formattedItems, fmt.Sprintf("%d: %s", next.First, next.Second))
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
