package ObserveOn

import (
	"fmt"
	"time"
)

type TestScheduler struct { Tasks []func() }
func (s *TestScheduler) Now() time.Time { return time.Now() }
func (s *TestScheduler) Schedule(task func()) { s.Tasks = append(s.Tasks, task) }
func (s *TestScheduler) ScheduleRecursive(task func(self func())) {}
func (s *TestScheduler) ScheduleFuture(due time.Duration, task func()) {}
func (s *TestScheduler) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {}
func (s *TestScheduler) Cancel() {}
func (s *TestScheduler) IsAsynchronous() bool { return len(s.Tasks) > 0 }

func Example_observeOn() {
	testScheduler := &TestScheduler{}

	// Observe by parking all next calls and the complete call on a custom scheduler
	source := FromInts(1, 2, 3, 4, 5).ObserveOn(testScheduler)
	subscription := source.Subscribe(func(next int, err error, done bool) {
		if !done {
			fmt.Printf("task %d\n", next)
		} else {
			fmt.Println("complete")
		}
	})

	// Observable ran to completion but nothing happended yet, all tasks have been parked
	fmt.Printf("%d parked tasks\n", len(testScheduler.Tasks))
	if !subscription.Closed() {
		fmt.Println("subscribed") // still subscribed, as complete is not yet delivered.
	}

	// Now let's run those tasks
	fmt.Println("---Hey Ho Let's Go!---")
	for _, task := range testScheduler.Tasks {
		task()
	}
	fmt.Println("--------")

	if subscription.Closed() {
		fmt.Println("unsubscribed") // complete should have caused subscription to close.
	}

	// Output:
	// 6 parked tasks
	// subscribed
	// ---Hey Ho Let's Go!---
	// task 1
	// task 2
	// task 3
	// task 4
	// task 5
	// complete
	// --------
	// unsubscribed
}

