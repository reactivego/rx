package scheduler

import "sync"

type SerialQueue struct {
	funcs chan func()
	count sync.WaitGroup
}

func StartSerialQueue() *SerialQueue {
	q := &SerialQueue{funcs: make(chan func())}
	q.count.Add(1)
	go func() {
		q.count.Wait()
		close(q.funcs)
	}()
	return q
}

func (q *SerialQueue) Go(f func()) {
	q.count.Add(1)
	go func() {
		defer q.count.Done()
		f()
	}()
}

func (q *SerialQueue) Do(f func()) {
	q.count.Add(1)
	done := make(chan struct{})
	q.funcs <- func() {
		defer q.count.Done()
		defer close(done)
		f()
	}
	<-done
}

func (q *SerialQueue) Wait() {
	q.count.Done()
	for f := range q.funcs {
		f()
	}
}
