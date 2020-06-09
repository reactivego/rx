package AutoConnect

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/reactivego/rx/generic"
	. "github.com/reactivego/rx/test"
)

// This example shows how to use AutoConnect with a PublishReplay connectable.
// The first Subscribe call will cause Connect on PublishReplay so it
// subscribes to range1to9. The second Subscribe call will cause the sequence
// to be replayed without doing a subscribe to range1to9.
func Example_autoConnect() {
	// Concurrent scheduler required for AutoConnect.
	concurrent := GoroutineScheduler()

	range1to9 := DeferInt(func() ObservableInt {
		fmt.Println("subscribed")
		i := 1
		return CreateRecursiveInt(func(N NextInt, E Error, C Complete) {
			if i < 10 {
				N(i)
				i++
			} else {
				C()
			}
		})
	})

	source := range1to9.PublishReplay(10, 0).AutoConnect(1).SubscribeOn(concurrent)

	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			fmt.Print(next, ",")
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	fmt.Println("first...")
	source.Subscribe(observe).Wait()
	fmt.Println("second...")
	source.Subscribe(observe).Wait()

	// Output:
	// first...
	// subscribed
	// 1,2,3,4,5,6,7,8,9,complete
	// second...
	// 1,2,3,4,5,6,7,8,9,complete
}

// This example how to use a concurrent scheduler, so multiple subscribe
// calls can be made to an AutoConnect(2) operator that will only connect on
// the second subscription.
func Example_autoConnectMulti() {
	// Concurrent scheduler required for AutoConnect.
	concurrent := GoroutineScheduler()

	range1to99 := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
		fmt.Println("subscribed")
		for i := 1; i < 100; i++ {
			if X() {
				return
			}
			N(i)
		}
		C()
	})

	source := range1to99.PublishReplay(10, 0).AutoConnect(2).SubscribeOn(concurrent)

	observe := func(next int, err error, done bool) {
		switch {
		case !done:
			// ignore values.
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("complete")
		}
	}

	fmt.Println("first...")

	// Subcribe is asynchronous now
	sub1 := source.Subscribe(observe)

	fmt.Println("second...")

	// Also asynchronous
	sub2 := source.Subscribe(observe)

	fmt.Println("wait...")

	// We now need to wait for the subscriptions to terminate.
	sub1.Wait()
	sub2.Wait()

	fmt.Println("done")

	// Output:
	// first...
	// second...
	// wait...
	// subscribed
	// complete
	// complete
	// done
}

func TestObservable(e *testing.T) {
	const ms = time.Millisecond

	Describ(e, "AutoConnect(0)", func(t T) {
		I(t, "should report error because count is less than 1", func(t T) {
			err := FromInt(1, 2, 3).Publish().AutoConnect(0).Wait()
			Asser(t).Error(err)
			Asser(t).Equal(err, ErrAutoConnectInvalidCount)
		})
	})
	Describ(e, "AutoConnect(1)", func(t T) {
		Contex(t, "trampoline scheduler", func(t T) {
			I(t, "should report error on subscribe", func(t T) {
				err := Interval(10 * ms).Publish().AutoConnect(1).Wait()
				Asser(t).Error(err)
				Asser(t).Equal(err, ErrSubjectNeedsConcurrentScheduler)
			})

		})
		Contex(t, "goroutine scheduler", func(t T) {
			Contex(t, "source stays active", func(t T) {
				I(t, "makes a multicaster behave like an observable", func(t T) {
					// AutoConnect makes a FooMulticaster behave like an ordinary
					// ObservableFoo that automatically connects the multicaster to its
					// source when the specified number of observers have subscribed to it.
					value, err := FromInt(42).Publish().AutoConnect(1).SubscribeOn(GoroutineScheduler()).ToSingle()

					Asser(t).NoError(err)
					Asser(t).Equal(value, 42, "= value")
				})
			})
			Contex(t, "source terminates with error", func(t T) {

				I(t, "reconnects when observer subscribes second time", func(t T) {
					concurrent := GoroutineScheduler()
					// When subsequently the next observer subscribes, AutoConnect will
					// connect to the source only when it was previously canceled or because
					// the source terminated with an error.
					subscriptions := 0
					observable := DeferInt(func() ObservableInt {
						subscriptions++
						if subscriptions == 1 {
							return ThrowInt(RxError("kaboom"))
						} else {
							return FromInt(subscriptions)
						}
					})

					autopub := observable.Publish().AutoConnect(1).SubscribeOn(concurrent)

					sub := autopub.Subscribe(func(int, error, bool) {})

					Asser(t).Equal(sub.Wait(), RxError("kaboom"), "= sub.Wait()")

					value, err := autopub.ToSingle()

					Asser(t).NoError(err)
					Asser(t).Equal(value, 2, "= value")

					Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
				})
			})
			Contex(t, "source completes", func(t T) {

				Contex(t, "publish", func(t T) {

					I(t, "does not reconnect when second observer subscribes", func(t T) {
						concurrent := GoroutineScheduler()
						// So it will not reconnect when the source completed succesfully. This
						// specific behavior allows for implementing a caching observable that
						// can be retried until it succeeds.
						subscriptions := 0
						observable := DeferInt(func() ObservableInt {
							subscriptions++
							return FromInt(subscriptions)
						})
						published := observable.Publish()
						autopub := published.AutoConnect(1).SubscribeOn(concurrent)

						value, err := autopub.ToSingle()
						Asser(t).Equal(subscriptions, 1, "= subscriptions")
						Asser(t).NoError(err)
						Asser(t).Equal(value, 1, "= value")

						// Second observer will not cause reconnect to source because first time
						// it completed.
						value, err = autopub.ToSingle()
						Asser(t).Equal(subscriptions, 1, "= subscriptions")
						Asser(t).Error(err)
						Asser(t).Equal(err, ErrSingleNoValue)
						Asser(t).Equal(value, 0, "= value")

						//Asser(t).Equal(concurrent.Count(), 0, "= concurrent.Count()")
						Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
					})
				})
				Contex(t, "publish replay", func(t T) {

					I(t, "replays value when second observer subscribes", func(t T) {
						concurrent := GoroutineScheduler()
						// So it will not reconnect when the source completed succesfully. This
						// specific behavior allows for implementing a caching observable that
						// can be retried until it succeeds.
						subscriptions := 0
						observable := DeferInt(func() ObservableInt {
							subscriptions++
							return FromInt(subscriptions)
						})
						published := observable.PublishReplay(1, 0)
						autopub := published.AutoConnect(1).SubscribeOn(concurrent)

						value, err := autopub.ToSingle()
						Asser(t).Equal(subscriptions, 1, "= subscriptions")
						Asser(t).NoError(err)
						Asser(t).Equal(value, 1, "= value")

						// Second observer will not cause reconnect to source because first time
						// it completed.
						value, err = autopub.ToSingle()
						Asser(t).Equal(subscriptions, 1, "= subscriptions")
						Asser(t).NoError(err)
						Asser(t).Equal(value, 1, "= value")

						Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
					})
				})
			})
		})
	})
	Describ(e, "AutoConnect(2)", func(t T) {
		I(t, "does not connect after first observer subscribes", func(t T) {
			concurrent := GoroutineScheduler()
			subscriptions := 0
			observable := DeferInt(func() ObservableInt {
				subscriptions++
				return FromInt(subscriptions)
			})

			published := observable.Publish().AutoConnect(2).SubscribeOn(concurrent)
			subscription1 := published.Subscribe(func(int, error, bool) {})

			time.Sleep(10 * ms)

			subscription1.Unsubscribe()
			subscription1.Wait()

			Asser(t).Equal(subscriptions, 0, "= subscriptions")

			concurrent.Wait()
			Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
		})
		I(t, "connects after second observer subscribes", func(t T) {
			concurrent := GoroutineScheduler()

			subscriptions := 0
			observable := DeferInt(func() ObservableInt {
				subscriptions++
				return FromInt(subscriptions)
			})
			published := observable.Publish().AutoConnect(2).SubscribeOn(concurrent)

			subscription1 := published.Subscribe(func(int, error, bool) {})
			subscription2 := published.Subscribe(func(int, error, bool) {})

			subscription1.Wait()
			subscription2.Wait()

			Asser(t).Equal(subscriptions, 1, "= subscriptions")

			concurrent.Wait()
			Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
		})
		I(t, "receives the same values on both observers", func(t T) {
			concurrent := GoroutineScheduler()

			expect := []int{42, 33, 51}
			published := FromInt(expect...).Publish().AutoConnect(2).SubscribeOn(concurrent)

			var actual1 []int
			subscription1 := published.Subscribe(func(next int, err error, done bool) {
				if !done {
					actual1 = append(actual1, next)
				}
			})
			defer subscription1.Unsubscribe()

			var actual2 []int
			subscription2 := published.Subscribe(func(next int, err error, done bool) {
				if !done {
					actual2 = append(actual2, next)
				}
			})
			defer subscription2.Unsubscribe()

			subscription1.Wait()
			subscription2.Wait()

			Asser(t).Equal(actual1, expect, "= actual1")
			Asser(t).Equal(actual2, expect, "= actual2")

			concurrent.Wait()
			Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
		})

		I(t, "should reconnect after timeout", func(t T) {
			concurrent := GoroutineScheduler()

			subscriptions := 0
			observable := DeferInt(func() ObservableInt {
				subscriptions++
				if subscriptions == 1 {
					return NeverInt()
				} else {
					return FromInt(subscriptions)
				}
			})

			published := observable.Timeout(500 * ms).Publish().AutoConnect(2).SubscribeOn(concurrent)

			subscription1 := published.Subscribe(func(int, error, bool) {})
			subscription2 := published.Subscribe(func(int, error, bool) {})

			err1 := subscription1.Wait()
			err2 := subscription2.Wait()

			Asser(t).Error(err1)
			Asser(t).Error(err2)

			Asser(t).Equal(err1, ErrTimeout)
			Asser(t).Equal(err2, ErrTimeout)

			Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")

			subscription1 = published.Subscribe(func(int, error, bool) {})
			subscription2 = published.Subscribe(func(int, error, bool) {})

			err1 = subscription1.Wait()
			err2 = subscription2.Wait()

			Asser(t).NoError(err1)
			Asser(t).NoError(err2)

			concurrent.Wait()
			Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
		})
		I(t, "cancels connect subscription if all observers unsubscribe", func(t T) {
			// Another thing to notice is that AutoConnect will disconnect an active
			// connection when the number of observers drops to zero. The reason for this is
			// that not doing so would leak a task and leave it hanging in the scheduler.
			concurrent := GoroutineScheduler()
			published := Timer(10 * ms).Publish()
			autopub := published.AutoConnect(2).SubscribeOn(concurrent)

			obs1 := autopub.Subscribe(func(int, error, bool) {})
			obs2 := autopub.Subscribe(func(int, error, bool) {})

			con := published.Connect()
			Asser(t).Must(con.Subscribed(), "con.Subscribed()")

			obs1.Unsubscribe()
			obs1.Wait()
			Asser(t).Must(con.Subscribed(), "con.Subscribed()")

			obs2.Unsubscribe()
			obs2.Wait()
			Asser(t).Not(con.Subscribed(), "con.Subscribed()")

			con.Wait()

			concurrent.Wait()
			Asser(t).Equal(fmt.Sprint(concurrent), "Goroutine{ tasks = 0 }")
		})
	})
}
