/*
Package rx provides Reactive Extensions for Go, an API for asynchronous
programming with observable streams.

Below is the list of supported operators:

Creating Operators

Operators that originate new Observables.

	Create           https://godoc.org/github.com/reactivego/rx/test/Create/
	                 https://godoc.org/github.com/reactivego/rx/test/CreateRecursive/
	                 https://godoc.org/github.com/reactivego/rx/test/CreateFutureRecursive/
	Defer            https://godoc.org/github.com/reactivego/rx/test/Defer/
	Empty            https://godoc.org/github.com/reactivego/rx/test/Empty/
	From             https://godoc.org/github.com/reactivego/rx/test/From/
	FromChan         https://godoc.org/github.com/reactivego/rx/test/FromChan/
	Interval         https://godoc.org/github.com/reactivego/rx/test/Interval/
	Just             https://godoc.org/github.com/reactivego/rx/test/Just/
	Never            https://godoc.org/github.com/reactivego/rx/test/Never/
	Of               https://godoc.org/github.com/reactivego/rx/test/Of/
	Range            https://godoc.org/github.com/reactivego/rx/test/Range/
	Start            https://godoc.org/github.com/reactivego/rx/test/Start/
	Throw            https://godoc.org/github.com/reactivego/rx/test/Throw/

Transforming Operators

Operators that transform items that are emitted by an Observable.

	Map              https://godoc.org/github.com/reactivego/rx/test/Map/
	Scan             https://godoc.org/github.com/reactivego/rx/test/Scan/

Filtering Operators

Operators that selectively emit items from a source Observable.

	Debounce         https://godoc.org/github.com/reactivego/rx/test/Debounce/
	Distinct         https://godoc.org/github.com/reactivego/rx/test/Distinct/
	ElementAt        https://godoc.org/github.com/reactivego/rx/test/ElementAt/
	Filter           https://godoc.org/github.com/reactivego/rx/test/Filter/ 
	First            https://godoc.org/github.com/reactivego/rx/test/First/
	IgnoreElements   https://godoc.org/github.com/reactivego/rx/test/IgnoreElements/
	IgnoreCompletion https://godoc.org/github.com/reactivego/rx/test/IgnoreCompletion/
	Last             https://godoc.org/github.com/reactivego/rx/test/Last/
	Sample           https://godoc.org/github.com/reactivego/rx/test/Sample/
	Single           https://godoc.org/github.com/reactivego/rx/test/Single/
	Skip             https://godoc.org/github.com/reactivego/rx/test/Skip/
	SkipLast         https://godoc.org/github.com/reactivego/rx/test/SkipLast/
	Take             https://godoc.org/github.com/reactivego/rx/test/Take/ 
	TakeLast         https://godoc.org/github.com/reactivego/rx/test/TakeLast/
	TakeUntil        https://godoc.org/github.com/reactivego/rx/test/TakeUntil/ 
	TakeWhile        https://godoc.org/github.com/reactivego/rx/test/TakeWhile/ 

Combining Operators

Below are the operators that flatten the emissions of multiple observables into
a single stream by subscribing to every observable stricly in sequence.
Observables may be added while the flattening is already going on.

	Concat           https://godoc.org/github.com/reactivego/rx/test/Concat/ 
	ConcatWith       https://godoc.org/github.com/reactivego/rx/test/ConcatWith/
	ConcatMap        https://godoc.org/github.com/reactivego/rx/test/ConcatMap/
	ConcatMapTo      https://godoc.org/github.com/reactivego/rx/test/ConcatMapTo/
	ConcatAll        https://godoc.org/github.com/reactivego/rx/test/ConcatAll/
	SwitchMap        https://godoc.org/github.com/reactivego/rx/test/SwitchMap/
	SwitchAll        https://godoc.org/github.com/reactivego/rx/test/SwitchAll/

Below are operators that flatten the emissions of multiple observables into a
single stream by subscribing to all observables concurrently. Here also,
observables may be added while the flattening is already going on.

	Merge            https://godoc.org/github.com/reactivego/rx/test/Merge/
	MergeWith        https://godoc.org/github.com/reactivego/rx/test/MergeWith/
	MergeMap         https://godoc.org/github.com/reactivego/rx/test/MergeMap/
	MergeAll         https://godoc.org/github.com/reactivego/rx/test/MergeAll/
	MergeDelayError  https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/
	MergeDelayErrorWith https://godoc.org/github.com/reactivego/rx/test/MergeDelayErrorWith/

Below are operators that flatten the emissions of multiple observables into a single observable
that emits slices of values. Differently from the previous two sets of operators,
these operators only start emitting once the list of observables to flatten is complete. 

	CombineLatest      https://godoc.org/github.com/reactivego/rx/test/CombineLatest/
	CombineLatestWith  https://godoc.org/github.com/reactivego/rx/test/CombineLatestWith/
	CombineLatestMap   https://godoc.org/github.com/reactivego/rx/test/CombineLatestMap/
	CombineLatestMapTo https://godoc.org/github.com/reactivego/rx/test/CombineLatestMapTo/
	CombineLatestAll   https://godoc.org/github.com/reactivego/rx/test/CombineLatestAll/

Error Handling Operators

Operators that help to recover from error notifications from an Observable.

	Catch            https://godoc.org/github.com/reactivego/rx/test/Catch/ 
	Retry            https://godoc.org/github.com/reactivego/rx/test/Retry/

Utility Operators

A toolbox of useful Operators for working with Observables.

	Delay            https://godoc.org/github.com/reactivego/rx/test/Delay/
	Do               https://godoc.org/github.com/reactivego/rx/test/Do/ 
	DoOnError        https://godoc.org/github.com/reactivego/rx/test/DoOnError/
	DoOnComplete     https://godoc.org/github.com/reactivego/rx/test/DoOnComplete/
	Finally          https://godoc.org/github.com/reactivego/rx/test/Finally/
	Passthrough      https://godoc.org/github.com/reactivego/rx/test/Passthrough/
	Repeat           https://godoc.org/github.com/reactivego/rx/test/Repeat/
	Serialize        https://godoc.org/github.com/reactivego/rx/test/Serialize/
	Timeout          https://godoc.org/github.com/reactivego/rx/test/Timeout/

Conditional and Boolean Operators

Operators that evaluate one or more Observables or items emitted by Observables.

	All              https://godoc.org/github.com/reactivego/rx/test/All/ ObservableBool

Aggregate Operators

Operators that operate on the entire sequence of items emitted by an Observable.

	Average           https://godoc.org/github.com/reactivego/rx/test/Average/
	Count             https://godoc.org/github.com/reactivego/rx/test/Count/
	Max               https://godoc.org/github.com/reactivego/rx/test/Max/
	Min               https://godoc.org/github.com/reactivego/rx/test/Min/
	Reduce            https://godoc.org/github.com/reactivego/rx/test/Reduce/
	Sum               https://godoc.org/github.com/reactivego/rx/test/Sum/

Type Casting and Type Filtering Operators

Operators to type cast, type filter observables.

	AsObservable      https://godoc.org/github.com/reactivego/rx/test/AsObservable
	Only              https://godoc.org/github.com/reactivego/rx/test/Only/

Scheduling Operators

Change the scheduler for subscribing and observing.

	ObserveOn         https://godoc.org/github.com/reactivego/rx/test/ObserveOn/
	SubscribeOn       https://godoc.org/github.com/reactivego/rx/test/SubscribeOn/

Multicasting Operators

A Connectable is an Observable that can multicast to observers subscribed to it.
The Connectable itself will subscribe to the Observable when the Connect method is called on it.

	Publish           https://godoc.org/github.com/reactivego/rx/test/Publish/ Connectable
	PublishReplay     https://godoc.org/github.com/reactivego/rx/test/PublishReplay/ Connectable

Connectable supports different strategies for subscribing to the Observable from which it was created.

	RefCount          https://godoc.org/github.com/reactivego/rx/test/RefCount/
	AutoConnect       https://godoc.org/github.com/reactivego/rx/test/AutoConnect/
	Connect           https://godoc.org/github.com/reactivego/rx/test/Connect

Subjects

A Subject is both a multicasting Observable as well as an Observer.
The Observable side allows multiple simultaneous subscribers.
The Observer side allows you to directly feed it data or subscribe it to another Observable.

	Subject           https://godoc.org/github.com/reactivego/rx/test/Subject
	ReplaySubject     https://godoc.org/github.com/reactivego/rx/test/ReplaySubject

Subscribing

Subscribing breathes life into a chain of observables. An observable may be subscribed to many times. 

Println and Subscribe implement subscribing behavior directly.

	Println          https://godoc.org/github.com/reactivego/rx/test/Println
	Subscribe        https://godoc.org/github.com/reactivego/rx/test/Subscribe
	Connect          https://godoc.org/github.com/reactivego/rx/test/Connect
	ToChan           https://godoc.org/github.com/reactivego/rx/test/ToChan
	ToSingle         https://godoc.org/github.com/reactivego/rx/test/ToSingle
	ToSlice          https://godoc.org/github.com/reactivego/rx/test/ToSlice
	Wait             https://godoc.org/github.com/reactivego/rx/test/Wait

Connect is called internally by RefCount and AutoConnect.

	RefCount         https://godoc.org/github.com/reactivego/rx/test/RefCount
	AutoConnect      https://godoc.org/github.com/reactivego/rx/test/AutoConnect
*/
package rx

import _ "github.com/reactivego/rx/generic"

