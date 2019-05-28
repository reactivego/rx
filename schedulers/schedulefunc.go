package schedulers

// ScheduleFunc scheduler is the type of a function that can schedule tasks.
type ScheduleFunc func(task func())

func (s ScheduleFunc) Schedule(task func()) {
	s(task)
}

func (s ScheduleFunc) Dispatch(task func()) {
	s(task)
}

func (s ScheduleFunc) IsAsynchronous() bool {
	return true
}

func (s ScheduleFunc) IsConcurrent() bool {
	return false
}

func (s ScheduleFunc) Wait(onCancel func(func())) {
}
