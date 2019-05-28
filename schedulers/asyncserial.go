package schedulers

// AsyncSerial scheduler schedules a task to occur after the currently
// running task completes. It has no synchronization mechanims to protect
// from concurrent access and is thus not supposed to be used from multiple
// goroutines. It should be used purely for scheduling tasks from a single
// goroutine.
type AsyncSerial struct{ tasks []func() }

// Schedule tasks asynchronously on a task queue. Call Wait to actually
// start running the scheduled tasks.
func (s *AsyncSerial) Schedule(task func()) {
	s.tasks = append(s.tasks, task)
}

// Dispatch tasks asynchronously on a task queue. Call Wait to actually
// start running the scheduled tasks.
func (s *AsyncSerial) Dispatch(task func()) {
	s.tasks = append(s.tasks, task)
}

// IsAsynchronous retruns true.
func (s AsyncSerial) IsAsynchronous() bool {
	return true
}

// IsConcurrent retruns false.
func (s AsyncSerial) IsConcurrent() bool {
	return false
}

// Wait will run the tasks from the queue by executing every task one by one.
// It will return when the function registered via onCancel() is called or
// when there are no more tasks remaining. Note, the currently running task
// may append additional tasks to the queue to run later.
func (s AsyncSerial) Wait(onCancel func(func())) {
	active := true
	onCancel(func() {
		active = false
	})
	for active && len(s.tasks) > 0 {
		s.tasks[0]()
		s.tasks = s.tasks[1:]
	}
}
