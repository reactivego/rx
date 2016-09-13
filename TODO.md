# TODO
Fix situations where Unsubscribe may cause goroutines to leak in filters.go
Implement ObserveOn
Integrate Scheduler based on SerialQueue implementation.
Test with ImmediateScheduler see if that allways works also when used together with Publish and Fork operators
Create test cases combining scheduler Go and Do with SubscribeOn and ObserveOn to look for nice conncurrency patterns.
Build a leak detector for goroutines, by never calling go directly, always via wrapper call that also intercepts panics.
Build test cases using early unsubscribe to test for proper cleanup of internal goroutines.
Switch to alternative interpretation of ObserverFunc using 'done' where 'done' indicates 'error or completed'.


I'm thinking about switching to an alternative interpretation of ObserverFunc
where I'm replacing the notion of 'completed' with 'done'. The done notion is true
both on error and when completed. The err argument is nil on completion and assigned
on error. This makes deciding between terminating and next much simpler.

func ObserverFunc(interface{},error,bool)

func (f ObserverFunc) Next(next interface{}) {
	f(next, nil, false)
}

func (f ObserverFunc) Done(err error) {
	f(nil, err, true)
}

func op1(next interface{}, err error, done bool)  {
	switch {
	case done:
		observer.Done(err)
	default:
		observer.Next(next)
	}
}

func op2(next interface{}, err error, done bool)  {
	if done {
		observer.Done(err)
	} else {
		observer.Next(next)
	}
}
