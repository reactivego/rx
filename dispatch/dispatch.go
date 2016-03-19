package dispatch

import (
	"errors"
	"runtime"
	"sync"
)

// Russ Cox presented a good solution for scheduling code on the main thread
// [here](https://groups.google.com/d/msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ)
// [see also](https://github.com/golang/go/wiki/LockOSThread)
func init() {
	// Arrange that main.main runs on main OS thread under following assumptions:
	//	1. current goroutine is the one that will be running main.main
	//	2. current thread is the main thread.
	// According to Russ's implementation these assumptions hold true.
	runtime.LockOSThread()
}

// SerialQueue is a serial dispatch queue for processing tasks by a
// single goroutine. This was originally designed to allow the dispatch of tasks
// from multiple goroutines to be performed on the main thread.
// However, it can be used as a general purpose access control mechanism to
// force operations on resources to be performed by a specific
// goroutine, thus obviating the need to use other synchronization primitives
// for acces control.
// The design of SerialQueue is inspired by Apple's Grand Central Dispatch library.
type SerialQueue struct {
	funcs chan func()
	count sync.WaitGroup
}

var Main = NewSerialQueue()

// NewSerialQueue creates a new initiliazed SerialQueue. Before any queue
// other than the Main queue can be used, it must be started. Once a queue
// has been started it must be Stopped whenever the queue is no longer needed.
func NewSerialQueue() *SerialQueue {
	return new(SerialQueue).Init()
}

// Init will initiallize a SerialQueue instance and return a pointer to it.
// Only should be called once. Called internally by the NewSerialQueue
// function (that everybody should be using to create SerialQueue instances).
func (q *SerialQueue) Init() *SerialQueue {
	q.funcs = make(chan func())
	return q
}

func (q *SerialQueue) dispatch() {
	for f := range q.funcs {
		f()
	}
}

func (q *SerialQueue) close() {
	q.count.Wait()
	close(q.funcs)
}

// Start should be called on any queue other than the Main queue to allow it
// to start performing tasks. This will then start a goroutine and and wait
// for tasks to be dispatched on the queue. Start will return before any
// task has been performed. To wait for all tasks to complete call Stop.
func (q *SerialQueue) Start() {
	if q == Main {
		panic(errors.New("Cannot start the Main queue, call Main() instead"))
	}
	go q.dispatch()
}

// Stop will signal the queue to stop performing tasks when it detects
// that the queue is idle. Note that all the tasks presently scheduled
// to be performed will be allowed to finish first. Tasks may even be
// added to the queue after Stop has been called. Stop will only return
// when the queue has been stopped.
func (q *SerialQueue) Stop() {
	if q == Main {
		panic(errors.New("Cannot Stop the Main queue, call Main() instead"))
	}
	q.close()
}

// Main must be called directly from main.main on the Main queue to
// actually start performing the tasks that have been dispatched on
// the Main queue. This method returns once all scheduled tasks have
// been performed.
func (q *SerialQueue) Main() {
	if q != Main {
		panic(errors.New("Main can only be called on the Main queue"))
	}
	go q.close()
	q.dispatch()
}

// DispatchSync will synchronously schedule a task to be performed
// on the queue. This method will only return when the work has
// been done. Don't call this recursively from the same goroutine
// because that will cause a deadlock.
func (q *SerialQueue) DispatchSync(f func()) {
	q.count.Add(1)
	done := make(chan struct{})
	q.funcs <- func() {
		defer q.count.Done()
		defer close(done)
		f()
	}
	<-done
}

// DispatchAsync will asynchronously schedule a task to be performed on
// the queue. This method will return before the task has been
// performed. The task will have been performed though before the
// queue's Main method or Stop method returns.
// Note that multiple DispatchAsync calls will usually NOT be scheduled
// in the same order in which they where dispatched to the Queue.
// The only guarantee is that they will be performed by the same
// goroutine, thus preventing simultaneous access by code executing
// concurrently.
func (q *SerialQueue) DispatchAsync(f func()) {
	q.count.Add(1)
	go func() {
		q.funcs <- func() {
			defer q.count.Done()
			f()
		}
	}()
}

// Go will start the passed in function in a separate goroutine.
// So this is a goroutine separate from one that is performing the tasks in the Queue.
// Using this method instead of a bare go statement will ensure that the Queue
// will not terminate on Stop before this goroutine has finshed.
// Not using this will (potentially) lead to a situation where the queue already
// has been stopped before the bare goroutine had a chance to dispatch tasks to
// the queue.
// Use this to run code that will dispatch tasks to the
func (q *SerialQueue) Go(f func(q *SerialQueue)) {
	q.count.Add(1)
	go func() {
		defer q.count.Done()
		f(q)
	}()
}
