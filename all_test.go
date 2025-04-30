package rx_test

import (
	"testing"

	"github.com/reactivego/rx"
)

func TestAll2(t *testing.T) {
	t.Run("Basic pair iteration", func(t *testing.T) {
		tuples := rx.From([]rx.Tuple2[string, int]{
			{"a", 1},
			{"b", 2},
			{"c", 3},
		}...)

		var firsts []string
		var seconds []int

		for f, s := range rx.All2(tuples) {
			firsts = append(firsts, f)
			seconds = append(seconds, s)
		}

		expectedFirsts := []string{"a", "b", "c"}
		expectedSeconds := []int{1, 2, 3}

		if len(firsts) != len(expectedFirsts) {
			t.Errorf("Expected %v firsts, got %v", len(expectedFirsts), len(firsts))
		}

		if len(seconds) != len(expectedSeconds) {
			t.Errorf("Expected %v seconds, got %v", len(expectedSeconds), len(seconds))
		}

		for i, v := range expectedFirsts {
			if firsts[i] != v {
				t.Errorf("Expected first %v at position %v, got %v", v, i, firsts[i])
			}
		}

		for i, v := range expectedSeconds {
			if seconds[i] != v {
				t.Errorf("Expected second %v at position %v, got %v", v, i, seconds[i])
			}
		}
	})

	t.Run("Empty tuple observable", func(t *testing.T) {
		tuples := rx.Empty[rx.Tuple2[string, int]]()
		count := 0

		for range rx.All2(tuples) {
			count++
		}

		if count != 0 {
			t.Errorf("Expected 0 tuples, got %v", count)
		}
	})

	t.Run("Direct sequence access", func(t *testing.T) {
		tuples := rx.From([]rx.Tuple2[string, int]{
			{"x", 10},
			{"y", 20},
		}...)

		var firsts []string
		var seconds []int

		for f, s := range rx.All2(tuples) {
			firsts = append(firsts, f)
			seconds = append(seconds, s)
		}

		expectedFirsts := []string{"x", "y"}
		expectedSeconds := []int{10, 20}

		if len(firsts) != len(expectedFirsts) || len(seconds) != len(expectedSeconds) {
			t.Errorf("Expected %v pairs, got %v", len(expectedFirsts), len(firsts))
		}

		for i := range expectedFirsts {
			if firsts[i] != expectedFirsts[i] || seconds[i] != expectedSeconds[i] {
				t.Errorf("At position %v, expected (%v, %v), got (%v, %v)",
					i, expectedFirsts[i], expectedSeconds[i], firsts[i], seconds[i])
			}
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("Basic indexed values", func(t *testing.T) {
		nums := rx.From("a", "b", "c")
		var indexes []int
		var values []string
		for i, v := range nums.All() {
			indexes = append(indexes, i)
			values = append(values, v)
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
		for i, v := range nums.All() {
			count++
			_, _ = i, v // Avoid unused variable warnings
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
		seq(func(i int, v int) bool {
			indexes = append(indexes, i)
			values = append(values, v)
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
}
