// ObserverObservable actually is an observer that is made observable. It is a
// buffering multicaster and has both an Observer and Observable side. These
// are then used as the core of any Subject implementation. The Observer side
// is used to pass items into the buffering multicaster. This then multicasts
// the items to every Observer that subscribes to the returned Observable.
package ObserverObservable

func (o Observer) Next(next interface{}) {
	o(next, nil, false)
}

func (o Observer) Error(err error) {
	o(nil, err, true)
}

func (o Observer) Complete() {
	o(nil, nil, true)
}
