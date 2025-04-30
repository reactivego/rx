package rx_test

import (
	"testing"

	"github.com/reactivego/rx"
)

func TestPull(t *testing.T) {
	t.Run("Basic Pull with Seq", func(t *testing.T) {
		// Create a simple iterator sequence
		seq := func(yield func(int) bool) {
			for i := 1; i <= 5; i++ {
				if !yield(i) {
					break
				}
			}
		}

		// Convert to an Observable using Pull
		obs := rx.Pull(seq)

		// Check results
		var results []int
		err := obs.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
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

	t.Run("Pull with early termination", func(t *testing.T) {
		// Create a sequence that would yield many values
		seq := func(yield func(int) bool) {
			for i := 1; i <= 100; i++ {
				if !yield(i) {
					break
				}
			}
		}

		// Convert to an Observable using Pull and take only a few
		obs := rx.Pull(seq).Take(3)

		// Check results
		var results []int
		err := obs.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := []int{1, 2, 3}
		if len(results) != len(expected) {
			t.Errorf("Expected %v items, got %v", len(expected), len(results))
		}
	})

	t.Run("Empty sequence", func(t *testing.T) {
		// Create an empty sequence
		seq := func(yield func(int) bool) {
			// No yields
		}

		// Convert to an Observable using Pull
		obs := rx.Pull(seq)

		// Check results
		var results []int
		err := obs.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 items, got %v", len(results))
		}
	})
}

func TestPull2(t *testing.T) {
	t.Run("Basic Pull2 with Seq2", func(t *testing.T) {
		// Create a simple key-value iterator sequence
		seq := func(yield func(string, int) bool) {
			data := map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			}

			for k, v := range data {
				if !yield(k, v) {
					break
				}
			}
		}

		// Convert to an Observable using Pull2
		obs := rx.Pull2(seq)

		// Check results
		var results []rx.Tuple2[string, int]
		err := obs.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(results) != 3 {
			t.Errorf("Expected 3 items, got %v", len(results))
		}

		// Convert results to a map for easy checking
		resultMap := make(map[string]int)
		for _, tuple := range results {
			resultMap[tuple.First] = tuple.Second
		}

		expectedMap := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}

		for k, v := range expectedMap {
			if resultMap[k] != v {
				t.Errorf("Expected %v for key %v, got %v", v, k, resultMap[k])
			}
		}
	})

	t.Run("Pull2 with filtering", func(t *testing.T) {
		// Create a simple key-value iterator sequence
		seq := func(yield func(string, int) bool) {
			data := map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
				"four":  4,
				"five":  5,
			}

			for k, v := range data {
				if !yield(k, v) {
					break
				}
			}
		}

		// Convert to an Observable using Pull2 and filter even values
		obs := rx.Pull2(seq).Filter(func(tuple rx.Tuple2[string, int]) bool {
			return tuple.Second%2 == 0
		})

		// Check results
		var results []rx.Tuple2[string, int]
		err := obs.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Map entries should only contain "two" and "four"
		if len(results) != 2 {
			t.Errorf("Expected 2 items, got %v", len(results))
		}

		// Convert results to a map for easy checking
		resultMap := make(map[string]int)
		// Manually iterate over results since ToMap is not available
		for _, tuple := range results {
			resultMap[tuple.First] = tuple.Second
		}

		expectedMap := map[string]int{
			"two":  2,
			"four": 4,
		}

		for k, v := range expectedMap {
			if resultMap[k] != v {
				t.Errorf("Expected %v for key %v, got %v", v, k, resultMap[k])
			}
		}
	})

	t.Run("Empty sequence with Pull2", func(t *testing.T) {
		// Create an empty sequence
		seq := func(yield func(int, string) bool) {
			// No yields
		}

		// Convert to an Observable using Pull2
		obs := rx.Pull2(seq)

		// Check results
		var results []rx.Tuple2[int, string]
		err := obs.Append(&results).Wait()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 items, got %v", len(results))
		}
	})
}
