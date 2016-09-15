# ObserverFunc

Switch to alternative interpretation of ObserverFunc using 'done' instead of 'completed'.
Using done makes the logic inside operations simpler because there is no need 
for "err!=nil || completed" checking anymore.
The done notion is true both on error and when completed. The err argument is nil on completion
and assigned on error.

func ObserverFunc(interface{},error,bool)

func (f ObserverFunc) Next(next interface{}) {
	f(next, nil, false)
}

func (f ObserverFunc) Done(err error) {
	f(nil, err, true)
}

func (f ObserverFunc) Completed() {
	f(nil, nil, true)
}

/*
// Previous implementation
func (f ObserverFunc) Error(err error) {
	f(nil, err, false)
}

// Previous implementation
func (f ObserverFunc) Completed() {
	f(nil, nil, true)
}
*/

func op(next interface{}, err error, done bool)  {
	if done {
		observer.Done(err)
	} else {
		observer.Next(next)
	}
}
