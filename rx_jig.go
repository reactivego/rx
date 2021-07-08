// This file guides regeneration of the heterogeneous rx package in this folder.
// The [jig tool](https://github.com/reactivego/jig) will generate rx.go guided
// by the code used in the generate function.

// +build ignore

package rx

import "time"

func generate() {
	_ = NewSubscriber()
	_ = GoroutineScheduler()
	t := MakeTrampolineScheduler()

	/*
		Observable Types
	*/

	var o Observable
	var b ObservableBool
	var i ObservableInt
	var oo ObservableObservable
	var oa ObservableSlice

	/*
		Observable Create Functions
	*/

	Create(func(Next, Error, Complete, Canceled) {})
	CreateRecursive(func(Next, Error, Complete) {})
	CreateFutureRecursive(time.Millisecond, func(Next, Error, Complete) time.Duration { return time.Millisecond })
	Defer(func() Observable { return o })
	Empty()
	From(1, 2)
	FromChan(make(chan interface{}))
	Of(1, 2)
	Interval(time.Millisecond)
	Just(1)
	Never()
	Range(1, 2)
	Start(func() (interface{}, error) { return o, nil })
	Throw(RxError("sad"))
	Ticker(time.Millisecond)
	Timer(time.Millisecond)

	/*
		Observable Operator Methods
	*/
	CombineLatest(o)
	o.CombineLatestWith(o)
	o.CombineLatestMap(func(interface{}) Observable { return nil })
	o.CombineLatestMapTo(o)
	oo.CombineLatestAll()

	Concat(o)
	o.ConcatWith()
	o.ConcatMap(nil)
	// o.ConcatMapTo
	// o.ConcatScan
	// o.ConcatReduce
	oo.ConcatAll()

	o.SwitchMap(func(interface{}) Observable { return o })
	// o.SwitchMapTo
	// o.SwitchScan
	// o.SwitchReduce
	oo.SwitchAll()

	// ExhaustMap
	// ExhaustMapTo
	// ExhaustAll

	Merge(o)
	o.MergeWith()
	o.MergeMap(func(interface{}) Observable { return o })
	// o.MergeMapTo(o)
	oo.MergeAll()

	MergeDelayError(o)
	o.MergeDelayErrorWith()
	//o.MergeDelayErrorMap(func(interface{}) Observable { return o })
	// o.MergeDelayErrorMapTo(o)
	//oo.MergeDelayErrorAll()

	o.WithLatestFrom(o)
	oo.WithLatestFromAll()

	o.All(func(next interface{}) bool { return true })
	o.AsObservableBool()
	o.AsObservableInt()
	b.AsObservable()
	i.AsObservable()
	o.AuditTime(time.Millisecond)
	i.Average()
	o.Buffer(o)
	o.BufferTime(time.Millisecond)
	o.Catch(o)
	o.CatchError(func(error, Observable) Observable { return o })
	o.Count()
	o.DebounceTime(time.Millisecond)
	o.Delay(time.Millisecond)
	o.Distinct()
	o.DistinctUntilChanged()
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
	o.Map(func(interface{}) interface{} { return nil })
	i.MapObservable(func(int) Observable { return nil })
	o.MapTo(nil)
	i.Max()
	i.Min()
	o.ObserveOn(func(task func()) { task() })
	o.OnlyBool()
	o.OnlyInt()
	// Passthrough
	o.Reduce(func(acc interface{}, value interface{}) interface{} { return nil }, nil)
	o.Repeat(1)
	o.Retry()
	o.SampleTime(time.Millisecond)
	o.Scan(func(acc interface{}, value interface{}) interface{} { return nil }, nil)
	o.Serialize()
	o.Single()
	o.Skip(1)
	o.SkipLast(1)
	o.StartWith(1)
	o.SubscribeOn(t)
	i.Sum()
	o.Take(1)
	i.Take(1)
	o.TakeLast(1)
	o.TakeUntil(o)
	o.TakeWhile(func(interface{}) bool { return true })
	o.ThrottleTime(time.Millisecond)
	o.TimeInterval()
	o.Timeout(time.Millisecond)
	o.Timestamp()

	/*
		Observable Subscribe Methods
	*/

	// Println
	o.Println()
	b.Println()
	i.Println()
	oa.Println()
	o.Subscribe(func(interface{}, error, bool) {})
	b.Subscribe(func(bool, error, bool) {})
	i.Subscribe(func(int, error, bool) {})
	oa.Subscribe(func([]interface{}, error, bool) {})
	o.ToChan()
	b.ToChan()
	i.ToChan()
	oa.ToChan()
	o.ToSingle()
	b.ToSingle()
	i.ToSingle()
	oa.ToSingle()
	o.ToSlice()
	b.ToSlice()
	i.ToSlice()
	oa.ToSlice()
	o.Wait()
	b.Wait()
	i.Wait()
	oa.Wait()

	Println()

	/*
		Multicast Types, Operator Methods and Connect Method
	*/

	var m Multicaster
	m.Connect()
	m.AutoConnect(0)
	m.RefCount()
	o.Publish()
	o.PublishReplay(0, 0)
}
