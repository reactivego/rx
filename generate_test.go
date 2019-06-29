package rx

import (
	"time"
)

// This example guides code generation by explicitly using the generics we want expanded into this package.
func Example_generate() {
	/*
		Observable Types
	*/

	var o Observable
	var b ObservableBool
	var i ObservableInt
	var oo ObservableObservable

	/*
		Observable Operator Methods
	*/

	o.All(func(next interface{}) bool { return true })
	o.AsObservableBool()
	o.AsObservableInt()
	b.AsObservable()
	i.AsObservable()
	i.Average()
	o.Catch(o)
	oo.CombineAll()
	o.Concat()
	oo.ConcatAll()
	o.Count()
	o.Debounce(time.Millisecond)
	o.Delay(time.Millisecond)
	o.Distinct()
	o.Do(func(interface{}) {})
	o.DoOnComplete(func() {})
	o.DoOnError(func(error) {})
	o.ElementAt(0)
	o.Filter(func(interface{}) bool { return true })
	o.Finally(func() {})
	o.First()
	o.IgnoreCompletion()
	o.IgnoreElements()
	o.Last()
	o.Map(func(interface{}) interface{} { return zero })
	i.MapObservable(func(int) Observable { return zeroObservable })
	i.Max()
	o.Merge()
	oo.MergeAll()
	o.MergeDelayError()
	o.MergeMap(func(interface{}) Observable { return o })
	i.Min()
	o.ObserveOn(CurrentGoroutineScheduler())
	o.OnlyBool()
	o.OnlyInt()
	// Passthrough
	o.Reduce(func(acc interface{}, value interface{}) interface{} { return zero }, zero)
	o.Repeat(1)
	o.Retry()
	o.Sample(time.Millisecond)
	o.Scan(func(acc interface{}, value interface{}) interface{} { return zero }, zero)
	// Serialize
	o.Single()
	o.Skip(1)
	o.SkipLast(1)
	o.SubscribeOn(CurrentGoroutineScheduler())
	i.Sum()
	oo.SwitchAll()
	o.SwitchMap(func(interface{}) Observable { return o })
	o.Take(1)
	i.Take(1)
	o.TakeLast(1)
	o.TakeUntil(o)
	o.TakeWhile(func(interface{}) bool { return true })
	o.Timeout(time.Millisecond)

	/*
		Observable Create Functions
	*/

	Create(func(Observer) {})
	Defer(func() Observable { return o })
	Empty()
	Error(RxError("sad"))
	From(1, 2)
	FromChan(make(chan interface{}))
	FromSlice([]interface{}{1, 2})
	Froms(1, 2)
	Interval(time.Millisecond)
	Just(1)
	Concat(o)
	Merge(o)
	Never()
	Range(1, 2)
	Repeat(zero, 1)
	Start(func() (interface{}, error) { return o, nil })
	Throw(RxError("sad"))

	/*
		Observable Subscribe Methods
	*/

	// Println
	o.Println()
	b.Println()
	i.Println()
	o.Subscribe(func(interface{}, error, bool) {})
	b.Subscribe(func(bool, error, bool) {})
	i.Subscribe(func(int, error, bool) {})
	o.SubscribeNext(func(interface{}) {})
	b.SubscribeNext(func(bool) {})
	i.SubscribeNext(func(int) {})
	o.ToChan()
	b.ToChan()
	i.ToChan()
	o.ToSingle()
	b.ToSingle()
	i.ToSingle()
	o.ToSlice()
	b.ToSlice()
	i.ToSlice()
	o.Wait()
	b.Wait()
	i.Wait()

	/*
		Connectable Types, Operator Methods and Subscribe Methods
	*/

	// var c Connectable
	// o.Publish()
	// p.PublishReplay()
	// c.Connect()
	// c.AutoConnect()
	// c.RefCount()

	/*
		Subject Types
	*/

	// Subject
	// ReplaySubject
}
