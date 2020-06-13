package rx

//jig:template Observable<Foo> Average

// Average calculates the average of numbers emitted by an ObservableFoo and
// emits this average.
func (o ObservableFoo) Average() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var sum, count foo
		observer := func(next foo, err error, done bool) {
			if !done {
				sum += next
				count++
			} else {
				if count > 0 {
					observe(sum/count, nil, false)
				}
				var zero foo
				observe(zero, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Count
//jig:needs ObservableInt

// Count counts the number of items emitted by the source ObservableFoo and
// emits only this value.
func (o ObservableFoo) Count() ObservableInt {
	observable := func(observe IntObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var count int
		observer := func(next foo, err error, done bool) {
			if !done {
				count++
			} else {
				observe(count, nil, false)
				observe(0, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Max

// Max determines, and emits, the maximum-valued item emitted by an
// ObservableFoo.
func (o ObservableFoo) Max() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var started bool
		var max foo
		observer := func(next foo, err error, done bool) {
			if started {
				if !done {
					if max < next {
						max = next
					}
				} else {
					observe(max, nil, false)
					var zero foo
					observe(zero, err, done)
				}
			} else {
				if !done {
					max = next
					started = true
				} else {
					var zero foo
					observe(zero, err, done)
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Min

// Min determines, and emits, the minimum-valued item emitted by an
// ObservableFoo.
func (o ObservableFoo) Min() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var started bool
		var min foo
		observer := func(next foo, err error, done bool) {
			if started {
				if !done {
					if min > next {
						min = next
					}
				} else {
					observe(min, nil, false)
					var zero foo
					observe(zero, err, done)
				}
			} else {
				if !done {
					min = next
					started = true
				} else {
					var zero foo
					observe(zero, err, done)
				}

			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Reduce<Bar>

// ReduceBar applies a reducer function to each item emitted by an ObservableFoo
// and the previous reducer result. The operator accepts a seed argument that
// is passed to the reducer for the first item emitted by the ObservableFoo.
// ReduceBar emits only the final value.
func (o ObservableFoo) ReduceBar(reducer func(bar, foo) bar, seed bar) ObservableBar {
	observable := func(observe BarObserver, subscribeOn Scheduler, subscriber Subscriber) {
		state := seed
		observer := func(next foo, err error, done bool) {
			if !done {
				state = reducer(state, next)
			} else {
				if err == nil {
					observe(state, nil, false)
				}
				var zeroBar bar
				observe(zeroBar, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Sum

// Sum calculates the sum of numbers emitted by an ObservableFoo and emits this sum.
func (o ObservableFoo) Sum() ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		var sum foo
		observer := func(next foo, err error, done bool) {
			if !done {
				sum += next
			} else {
				observe(sum, nil, false)
				var zero foo
				observe(zero, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}
