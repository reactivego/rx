package rx

//jig:template Observable<Foo> Map<Bar>

// MapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item.
func (o ObservableFoo) MapBar(project func(foo) bar) ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			var mapped bar
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> MapAs<Bar>

// MapAsBar transforms the items emitted by an ObservableFoo by applying a
// function to each item.
func (o ObservableFoo) MapAsBar(project func(foo) bar) ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			var mapped bar
			if !done {
				mapped = project(next)
			}
			observe(mapped, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> MapTo<Bar>

// MapToBar transforms the items emitted by an ObservableFoo. Emitted values
// are mapped to the same value every time.
func (o ObservableFoo) MapToBar(value bar) ObservableBar {
	return o.MapBar(func(foo) bar { return value })
}

//jig:template Observable<Foo> Scan<Bar>

// ScanBar applies a accumulator function to each item emitted by an
// ObservableFoo and the previous accumulator result. The operator accepts a
// seed argument that is passed to the accumulator for the first item emitted
// by the ObservableFoo. ScanBar emits every value, both intermediate and final.
func (o ObservableFoo) ScanBar(accumulator func(bar, foo) bar, seed bar) ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next foo, err error, done bool) {
			if !done {
				state = accumulator(state, next)
				observe(state, nil, false)
			} else {
				var zeroBar bar
				observe(zeroBar, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
