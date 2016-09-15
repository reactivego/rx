package scheduler

// Scheduler

type Scheduler interface {
	Asynchronous() bool
	Schedule(task func())
}

// Immediate

type immediate struct{}

func (s immediate) Asynchronous() bool {
	return false
}

func (s immediate) Schedule(task func()) {
	task()
}

var Immediate immediate

// SchedulerFunc

type SchedulerFunc func(task func())

func (s SchedulerFunc) Asynchronous() bool {
	return true
}

func (s SchedulerFunc) Schedule(task func()) {
	s(task)
}

// CurrentGoroutine
// Tasks are asynchronously scheduled on the current goroutine via dispatch queue

// NewGoroutine
// Tasks are asynchronously scheduled on a single new goroutine via dispatch queue

// Goroutines
// Tasks are scheduled on new goroutines

var Goroutines SchedulerFunc = func(task func()) { go task() }
