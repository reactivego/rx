package rx

//jig:template Observable<Foo> Map<Bar>

// MapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item.
func (o ObservableFoo) MapBar(project func(foo) bar) ObservableBar {
	observable := func(observe BarObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
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

//jig:template Observable<Foo> SwitchMap<Bar>

// SwitchMapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item an returning an ObservableBar. In doing so, it behaves much like
// MergeMap (previously FlatMap), except that whenever a new ObservableBar is emitted
// SwitchMap will unsubscribe from the previous ObservableBar and begin emitting items
// from the newly emitted one.
func (o ObservableFoo) SwitchMapBar(project func(foo) ObservableBar) ObservableBar {
	return o.MapObservableBar(project).SwitchAll()
}

//jig:template Observable<Foo> MergeMap<Bar>

// MergeMapBar transforms the items emitted by an ObservableFoo by applying a
// function to each item an returning an ObservableBar. The stream of ObservableBar
// items is then merged into a single stream of Bar items using the MergeAll operator.
func (o ObservableFoo) MergeMapBar(project func(foo) ObservableBar) ObservableBar {
	return o.MapObservableBar(project).MergeAll()
}

//jig:template Observable<Foo> Scan<Bar>

// ScanBar applies a accumulator function to each item emitted by an
// ObservableFoo and the previous accumulator result. The operator accepts a
// seed argument that is passed to the accumulator for the first item emitted
// by the ObservableFoo. ScanBar emits every value, both intermediate and final.
func (o ObservableFoo) ScanBar(accumulator func(bar, foo) bar, seed bar) ObservableBar {
	observable := func(observe BarObserveFunc, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next foo, err error, done bool) {
			if !done {
				state = accumulator(state, next)
				observe(state, nil, false)
			} else {
				observe(zeroBar, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
