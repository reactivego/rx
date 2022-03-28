package observable

import (
	"github.com/reactivego/scheduler"
)

type Scheduler = scheduler.Scheduler

type ConcurrentScheduler = scheduler.ConcurrentScheduler

var Goroutine = scheduler.Goroutine

var NewScheduler = scheduler.New
