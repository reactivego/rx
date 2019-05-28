package schedulers

import "sync"

// Goroutine scheduler schedules every task asynchronously to run concurrently
// as a new goroutine.
type Goroutine struct{}

// Schedule a task asynchronously to run concurrently as a new goroutine.
func (s Goroutine) Schedule(task func()) {
	go task()
}

// Dispatch a task asynchronously to run concurrently as a new goroutine.
func (s Goroutine) Dispatch(task func()) {
	go task()
}

// IsAsynchronous retruns true.
func (s Goroutine) IsAsynchronous() bool {
	return true
}

// IsConcurrent retruns true.
func (s Goroutine) IsConcurrent() bool {
	return true
}

// Wait will block until the function registered via onCancel() is called.
// Goroutines will run by themselves and so there is nothing to do but wait.
func (s Goroutine) Wait(onCancel func(func())) {
	var wg sync.WaitGroup
	wg.Add(1)
	onCancel(wg.Done)
	wg.Wait()
}
