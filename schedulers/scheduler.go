package schedulers

type Scheduler interface {
	Schedule(task func())
}

type SchedulerFunc func(task func())

func (s SchedulerFunc) Schedule(task func()) {
	s(task)
}

var ImmediateScheduler SchedulerFunc = func(task func()) { task() }

var GoroutineScheduler SchedulerFunc = func(task func()) { go task() }
