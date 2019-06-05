package ObserveOn

import (
	"fmt"

	"github.com/reactivego/scheduler"
)

func Example_observeOn() {
	// Simple custom task scheduler
	var tasks []func()
	taskScheduler := scheduler.ScheduleFunc(func(task func()) {
		tasks = append(tasks, task)
	})

	// Observe by parking all next calls and the complete call on a custom scheduler
	source := FromInts(1, 2, 3, 4, 5).ObserveOn(taskScheduler)
	subscription := source.Subscribe(func(next int, err error, done bool) {
		if !done {
			fmt.Printf("task %d\n", next)
		} else {
			fmt.Println("complete")
		}
	})

	// Observable ran to completion but nothing happended yet, all tasks have been parked
	fmt.Printf("%d parked tasks\n", len(tasks))
	if !subscription.Closed() {
		fmt.Println("subscribed") // still subscribed, as complete is not yet delivered.
	}

	// Now let's run those tasks
	fmt.Println("---Hey Ho Let's Go!---")
	for _, task := range tasks {
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
