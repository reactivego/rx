package rx_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/reactivego/rx"
)

func TestZipAll(t *testing.T) {
	t.Run("Basic zipping behavior", func(t *testing.T) {
		// Create two source observables
		nums := rx.From(1, 2, 3)
		chars := rx.From("a", "b", "c")

		// Create an observable of observables
		sources := rx.From(
			nums.AsObservable(),
			chars.AsObservable(),
		)

		// Use ZipAll to combine them
		zipped := rx.ZipAll[any](sources)

		// Collect results
		var results [][]any

		// Use non-concurrent scheduler for predictable test execution
		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, scheduler).Wait()

		// Verify expected results
		expected := [][]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	t.Run("Different length sources", func(t *testing.T) {
		// Source 1 has more elements than source 2
		nums := rx.From(1, 2, 3, 4, 5)
		chars := rx.From("a", "b", "c")

		sources := rx.From(
			nums.AsObservable(),
			chars.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any
		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, scheduler).Wait()

		// ZipAll should stop when the shortest source completes
		expected := [][]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	t.Run("Empty source observable", func(t *testing.T) {
		// No source observables provided
		var sources rx.Observable[rx.Observable[any]]
		sources = rx.From[rx.Observable[any]]() // Empty source
		zipped := rx.ZipAll[any](sources)

		completed := false
		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if done && err == nil {
				completed = true
			}
		}, scheduler).Wait()

		if !completed {
			t.Error("Expected observable to complete when source is empty")
		}
	})

	t.Run("Error in source observable", func(t *testing.T) {
		nums := rx.From(1, 2, 3)

		// Observable that emits an error
		errObs := rx.Create(func(index int) (string, error, bool) {
			switch index {
			case 0:
				return "a", nil, false
			case 1:
				return "", errors.New("test error"), true
			default:
				return "", nil, true
			}
		})

		sources := rx.From(
			nums.AsObservable(),
			errObs.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any
		var receivedErr error

		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			} else if err != nil {
				receivedErr = err
			}
		}, scheduler).Wait()

		expected := [][]any{
			{1, "a"},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}

		if receivedErr == nil || receivedErr.Error() != "test error" {
			t.Errorf("Expected 'test error', got %v", receivedErr)
		}
	})

	t.Run("One source completes early with empty buffer", func(t *testing.T) {
		// Create a source that completes after emitting just one item
		earlyComplete := rx.Create(func(index int) (int, error, bool) {
			switch index {
			case 0:
				return 1, nil, false
			default:
				return 0, nil, true
			}
		})

		// Create a source that emits multiple items
		multiValue := rx.From(10, 20, 30, 40)

		sources := rx.From(
			earlyComplete.AsObservable(),
			multiValue.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any

		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, scheduler).Wait()

		// Should only emit one pair, then complete
		expected := [][]any{
			{1, 10},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	t.Run("Multiple sources with different completion times", func(t *testing.T) {
		// Three sources with different numbers of elements
		src1 := rx.From(1, 2)
		src2 := rx.From("a", "b", "c", "d")
		src3 := rx.From(true, false, true)

		sources := rx.From(
			src1.AsObservable(),
			src2.AsObservable(),
			src3.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any

		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, scheduler).Wait()

		// Should only emit based on the shortest source (src1)
		expected := [][]any{
			{1, "a", true},
			{2, "b", false},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	t.Run("Unsubscribe during emission", func(t *testing.T) {
		// Create sources
		src1 := rx.From(1, 2, 3, 4, 5)
		src2 := rx.From("a", "b", "c", "d", "e")

		sources := rx.From(
			src1.AsObservable(),
			src2.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any
		var subscription rx.Subscription

		scheduler := rx.NewScheduler()
		subscription = zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
				// Unsubscribe after receiving 2 emissions
				if len(results) >= 2 {
					subscription.Unsubscribe()
				}
			}
		}, scheduler)

		subscription.Wait()

		// Should only emit two pairs due to unsubscription
		expected := [][]any{
			{1, "a"},
			{2, "b"},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	t.Run("Parallel execution test", func(t *testing.T) {
		// Using goroutine scheduler for true parallelism

		// Create sources
		src1 := rx.From(1, 2, 3)
		src2 := rx.From("a", "b", "c")

		sources := rx.From(
			src1.AsObservable(),
			src2.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any

		// Use Goroutine scheduler for parallel execution
		scheduler := rx.Goroutine
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, scheduler).Wait()

		// Results should still be zipped in order despite parallel execution
		expected := [][]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	t.Run("Interleaved emissions test", func(t *testing.T) {
		// Create source that alternates between emitting and not emitting
		alternate := rx.Create(func(index int) (int, error, bool) {
			if index >= 6 {
				return 0, nil, true
			}
			// Only emit on even indices, to create "delays" between emissions
			if index%2 == 0 {
				return index/2 + 1, nil, false
			}
			// Skip odd indices (simulating delay)
			return 0, nil, false
		}).Filter(func(i int) bool { return i != 0 })

		// Create source that emits on every call
		regular := rx.Create(func(index int) (string, error, bool) {
			if index >= 3 {
				return "", nil, true
			}
			return string('a' + byte(index)), nil, false
		})

		sources := rx.From(
			alternate.AsObservable(),
			regular.AsObservable(),
		)

		zipped := rx.ZipAll[any](sources)

		var results [][]any

		scheduler := rx.NewScheduler()
		zipped.Subscribe(func(next []any, err error, done bool) {
			if !done {
				results = append(results, next)
			}
		}, scheduler).Wait()

		expected := [][]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})
}
