package rx

//jig:template Observable Distinct

// Distinct suppress duplicate items emitted by an Observable
func (o Observable) Distinct() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		seen := map[interface{}]struct{}{}
		observer := func(next interface{}, err error, done bool) {
			if !done {
				if _, present := seen[next]; present {
					return
				}
				seen[next] = struct{}{}
			}
			observe(next, err, done)
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Distinct
//jig:needs Observable Distinct

// Distinct suppress duplicate items emitted by an Observable
func (o ObservableFoo) Distinct() ObservableFoo {
	return o.AsObservable().Distinct().AsObservableFoo()
}

//jig:template Observable ElementAt

// ElementAt emit only item n emitted by an Observable
func (o Observable) ElementAt(n int) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		observer := func(next interface{}, err error, done bool) {
			if done || i == n {
				observe(next, err, done)
			}
			i++
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> ElementAt
//jig:needs Observable ElementAt

// ElementAt emit only item n emitted by an Observable
func (o ObservableFoo) ElementAt(n int) ObservableFoo {
	return o.AsObservable().ElementAt(n).AsObservableFoo()
}

//jig:template Observable<Foo> Filter

// Filter emits only those items from an ObservableFoo that pass a predicate test.
func (o ObservableFoo) Filter(predicate func(next foo) bool) ObservableFoo {
	observable := func(observe FooObserver, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next foo, err error, done bool) {
			if done || predicate(next) {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable First

// First emits only the first item, or the first item that meets a condition, from an Observable.
func (o Observable) First() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		start := true
		observer := func(next interface{}, err error, done bool) {
			if done || start {
				observe(next, err, done)
			}
			start = false
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> First
//jig:needs Observable First

// First emits only the first item, or the first item that meets a condition, from an ObservableFoo.
func (o ObservableFoo) First() ObservableFoo {
	return o.AsObservable().First().AsObservableFoo()
}

//jig:template Observable IgnoreElements

// IgnoreElements does not emit any items from an Observable but mirrors its termination notification.
func (o Observable) IgnoreElements() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if done {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> IgnoreElements
//jig:needs Observable IgnoreElements

// IgnoreElements does not emit any items from an ObservableFoo but mirrors its termination notification.
func (o ObservableFoo) IgnoreElements() ObservableFoo {
	return o.AsObservable().IgnoreElements().AsObservableFoo()
}

//jig:template Observable IgnoreCompletion

// IgnoreCompletion only emits items and never completes, neither with Error nor with Complete.
func (o Observable) IgnoreCompletion() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		observer := func(next interface{}, err error, done bool) {
			if !done {
				observe(next, err, done)
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> IgnoreCompletion
//jig:needs Observable IgnoreCompletion

// IgnoreCompletion is unknown on the reactivex.io site do we need this?
func (o ObservableFoo) IgnoreCompletion() ObservableFoo {
	return o.AsObservable().IgnoreCompletion().AsObservableFoo()
}

//jig:template Observable Last

// Last emits only the last item emitted by an Observable.
func (o Observable) Last() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		have := false
		var last interface{}
		observer := func(next interface{}, err error, done bool) {
			if done {
				if have {
					observe(last, nil, false)
				}
				observe(nil, err, true)
			} else {
				last = next
				have = true
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Last
//jig:needs Observable Last

// Last emits only the last item emitted by an ObservableFoo.
func (o ObservableFoo) Last() ObservableFoo {
	return o.AsObservable().Last().AsObservableFoo()
}

//jig:template Observable Single
//jig:needs RxError

// Single enforces that the observable sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing, this is reported as an error to the observer.
func (o Observable) Single() Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		var (
			count  int
			latest interface{}
		)
		observer := func(next interface{}, err error, done bool) {
			if count < 2 {
				if done {
					if err != nil {
						observe(nil, err, true)
					} else {
						if count == 1 {
							observe(latest, nil, false)
							observe(nil, nil, true)
						} else {
							observe(nil, RxError("expected one value, got none"), true)
						}
					}
				} else {
					count++
					if count == 1 {
						latest = next
					} else {
						observe(nil, RxError("expected one value, got multiple"), true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Single

// Single enforces that the observableFoo sends exactly one data item and then
// completes. If the observable sends no data before completing or sends more
// than 1 item before completing  this reported as an error to the observer.
func (o ObservableFoo) Single() ObservableFoo {
	return o.AsObservable().Single().AsObservableFoo()
}

//jig:template Observable Skip

// Skip suppresses the first n items emitted by an Observable.
func (o Observable) Skip(n int) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		i := 0
		observer := func(next interface{}, err error, done bool) {
			if done || i >= n {
				observe(next, err, done)
			}
			i++
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Skip

// Skip suppresses the first n items emitted by an ObservableFoo.
func (o ObservableFoo) Skip(n int) ObservableFoo {
	return o.AsObservable().Skip(n).AsObservableFoo()
}

//jig:template Observable SkipLast

// SkipLast suppresses the last n items emitted by an Observable.
func (o Observable) SkipLast(n int) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		read := 0
		write := 0
		n++
		buffer := make([]interface{}, n)
		observer := func(next interface{}, err error, done bool) {
			if done {
				observe(nil, err, true)
			} else {
				buffer[write] = next
				write = (write + 1) % n
				if write == read {
					observe(buffer[read], nil, false)
					read = (read + 1) % n
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> SkipLast
//jig:needs Observable SkipLast

// SkipLast suppresses the last n items emitted by an ObservableFoo.
func (o ObservableFoo) SkipLast(n int) ObservableFoo {
	return o.AsObservable().SkipLast(n).AsObservableFoo()
}

//jig:template Observable Take

// Take emits only the first n items emitted by an Observable.
func (o Observable) Take(n int) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		taken := 0
		observer := func(next interface{}, err error, done bool) {
			if taken < n {
				observe(next, err, done)
				if !done {
					taken++
					if taken >= n {
						observe(nil, nil, true)
					}
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> Take
//jig:needs Observable Take

// Take emits only the first n items emitted by an ObservableFoo.
func (o ObservableFoo) Take(n int) ObservableFoo {
	return o.AsObservable().Take(n).AsObservableFoo()
}

//jig:template Observable TakeLast

// TakeLast emits only the last n items emitted by an Observable.
func (o Observable) TakeLast(n int) Observable {
	observable := func(observe Observer, subscribeOn Scheduler, subscriber Subscriber) {
		read := 0
		write := 0
		n++
		buffer := make([]interface{}, n)
		observer := func(next interface{}, err error, done bool) {
			if done {
				for read != write {
					observe(buffer[read], nil, false)
					read = (read + 1) % n
				}
				observe(nil, err, true)
			} else {
				buffer[write] = next
				write = (write + 1) % n
				if write == read {
					read = (read + 1) % n
				}
			}
		}
		o(observer, subscribeOn, subscriber)
	}
	return observable
}

//jig:template Observable<Foo> TakeLast
//jig:needs Observable TakeLast

// TakeLast emits only the last n items emitted by an ObservableFoo.
func (o ObservableFoo) TakeLast(n int) ObservableFoo {
	return o.AsObservable().TakeLast(n).AsObservableFoo()
}
