package ConcatAll

import (
	"testing"
	"time"

	. "github.com/reactivego/rx/test"
	_ "github.com/reactivego/rx"
)

func Example_concatAll() {
	Interval(time.Millisecond).Take(3).MapObservableInt(func(next int) ObservableInt {
		return RangeInt(next, 2)
	}).ConcatAll().Println()

	// Output:
	// 0
	// 1
	// 1
	// 2
	// 2
	// 3
}

func TestObservable_ConcatAll(e *testing.T) {
	Describ(e, "emits", func(t BDD) {
		const ms = time.Millisecond

		Contex(t, "with a fast source", func(t BDD) {

			Contex(t, "and a slow concat", func(t BDD) {

				I(t, "should complete before the first emit to observer", func(t BDD) {
					f := func(i int) ObservableInt {
						return JustInt(i).Delay(10 * ms)
					}

					var complete time.Time
					c := func() {
						complete = time.Now()
					}

					var next time.Time
					n := func(int) {
						if next.IsZero() {
							next = time.Now()
						}
					}

					expect := []int{1, 2, 3}
					actual, err := FromInt(expect...).DoOnComplete(c).MapObservableInt(f).ConcatAll().Do(n).ToSlice()

					Asser(t, NoError(err))
					Asser(t, Equal(actual, expect))
					Asser(t, Not(complete.IsZero(), "zero complete time"))
					Asser(t, Not(next.IsZero(), "zero next time"))
					Asser(t, Must(complete.Before(next), "complete.Before(next), but complete ", complete.Sub(next), " after next"))
				})

				I(t, "should emit all concatenated observables", func(t BDD) {
					f := func(i int) ObservableInt {
						return JustInt(i).Delay(10 * ms)
					}
					expect := []int{1, 2, 3}
					actual, err := FromInt(expect...).MapObservableInt(f).ConcatAll().ToSlice()

					Asser(t, NoError(err))
					Asser(t, Equal(actual, expect))
				})
			})
		})

		Contex(t, "with a slow source", func(t BDD) {

			Contex(t, "and a fast concat", func(t BDD) {

				I(t, "should start emitting to observer before source has completed", func(t BDD) {
					f := func(i int) ObservableInt {
						return JustInt(i)
					}

					var complete time.Time
					c := func() {
						complete = time.Now()
					}

					var next time.Time
					n := func(int) {
						if next.IsZero() {
							next = time.Now()
						}
					}

					Interval(10 * ms).Take(3).DoOnComplete(c).MapObservableInt(f).ConcatAll().Do(n).Wait()

					Asser(t, Not(complete.IsZero(), "complete time is zero"))
					Asser(t, Not(next.IsZero(), "next time is zero"))
					Asser(t, Must(next.Before(complete), "next.Before(complete) but next ", next.Sub(complete), " after complete"))
				})

				I(t, "should emit all concatenated observables", func(t BDD) {
					f := func(i int) ObservableInt {
						return JustInt(i)
					}
					expect := []int{0, 1, 2}
					actual, err := Interval(10 * ms).Take(3).MapObservableInt(f).ConcatAll().ToSlice()

					Asser(t, NoError(err))
					Asser(t, Equal(actual, expect))
				})
			})
		})
	})
}