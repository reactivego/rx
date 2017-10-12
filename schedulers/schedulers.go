// Package schedulers contains Scheduler implementations.
package schedulers

// Trampoline scheduler schedules a task to occur after the currently
// running task completes. Trampoline is not supposed to be used from
// multiple goroutines. It should be used purely for recursive scheduling
// tasks from a single goroutine.
type Trampoline struct{ tasks []func() }

// Schedule the first task to run synchronously and any subsequent tasks
// asynchronously on a task queue. So when the first task eventually
// returns the queue of tasks is empty again.
func (s *Trampoline) Schedule(task func()) {
	s.tasks = append(s.tasks, task)
	if len(s.tasks) == 1 {
		for len(s.tasks) > 0 {
			s.tasks[0]()
			s.tasks = s.tasks[1:]
		}
	}
}

// Goroutine scheduler schedules every task asynchronously on its own
// goroutine.
type Goroutine struct{}

// Schedule a task asynchronously to run concurrently as a new goroutine.
func (s Goroutine) Schedule(task func()) { go task() }

// ScheduleFunc is the type of a function that can schedule tasks.
type ScheduleFunc func(task func())

// Schedule schedules a task using a schedule func.
func (s ScheduleFunc) Schedule(task func()) { s(task) }
