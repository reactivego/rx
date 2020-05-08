/*
Package rx provides Go with Reactive Extensions, an API for asynchronous
programming with Observables.

Regenerating this Package

This package is generated from the sub-folder generic by the jig tool.
You don't need to regenerate the package in order to use it. However, if you
are interested in regenerating it, then read on.

The jig tool provides the parametric polymorphism capability that Go 1 is
missing. It works by replacing place-holder types of generic functions and
datatypes with interface{} (it can also generate statically typed code though).

To regenerate, change the current working directory to the package directory
and run the jig tool as follows:

	$ go get -d github.com/reactivego/generics/cmd/jig
	$ go run github.com/reactivego/generics/cmd/jig -v

Observables

The main focus of rx is on Observables.

An Observable; is a stream of events, assumes zero to many values over time,
pushes values, can take any amount of time to complete (or may never), is
cancellable, is lazy (it doesn't do anything until you subscribe).

This package uses `interface{}` for entry types, so an observable can emit a
mix of differently typed entries. To create an observable that emits three
values of different types you could write the following little program.

	package main

	import "github.com/reactivego/rx"

	func main() {
		rx.From(1,"hi",2.3).Println()
	}

The code above creates an observable from numbers and strings and then prints
them.

Observables in rx are somewhat similar to Go channels but have much richer
semantics:

Observables can be hot or cold. A hot observable will try to emit values even
when nobody is subscribed. As long as there are no subscribers the values of
a hot observable are lost. The position of a mouse pointer or the current time
are examples of hot observables. 

A cold observable will only start emitting values when somebody subscribes.
The contents of a file or a database are examples of cold observables.

An observable can complete normally or with an error, it uses subscriptions
that can be canceled from the subscriber side. Where a normal variable is
just a place where you read and write values from, an observable captures how
the value of this variable changes over time.

Concurrency follows naturally from the fact that an observable is an ever
changing stream of values.

Operators 

The combination of Observables and a set of expressive Operators is the real
strength of Reactive Extensions. Operators work on one or more Observables.
They are the language you use to describe the way in which observables should
be combined to form new Observables. Operators specify how Observables
representing streams of values are e.g. merged, transformed, concatenated,
split, multicasted, replayed, delayed and debounced.

This implementation takes most of its cues from
RxJS 6 (https://github.com/ReactiveX/rxjs) and
RxJava 2 (https://github.com/ReactiveX/RxJava).
Both libaries have been pushing the envelope in evolving operator semantics.

Below is the list of implemented operators.

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
	Repeat           https://godoc.org/github.com/reactivego/rx/test/Repeat/
	Start            https://godoc.org/github.com/reactivego/rx/test/Start/
	Throw            https://godoc.org/github.com/reactivego/rx/test/Throw/

Transforming Operators

Operators that transform items that are emitted by an Observable.

	ConcatMap 
	Map              https://godoc.org/github.com/reactivego/rx/test/Map/
	MergeMap         https://godoc.org/github.com/reactivego/rx/test/MergeMap/
	Scan             https://godoc.org/github.com/reactivego/rx/test/Scan/
	SwitchMap        https://godoc.org/github.com/reactivego/rx/test/SwitchMap/

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

Operators that work with multiple source Observables to create a single Observable.

	CombineAll       https://godoc.org/github.com/reactivego/rx/test/CombineAll/
	CombineLatest     
	Concat           https://godoc.org/github.com/reactivego/rx/test/Concat/ 
	ConcatAll        https://godoc.org/github.com/reactivego/rx/test/ConcatAll/
	Merge            https://godoc.org/github.com/reactivego/rx/test/Merge/
	MergeAll         https://godoc.org/github.com/reactivego/rx/test/MergeAll/
	MergeDelayError	 https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/
	MergeDelayError	 https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/
	SwitchAll        https://godoc.org/github.com/reactivego/rx/test/SwitchAll/

Error Handling Operators

Operators that help to recover from error notifications from an Observable.

	Catch            https://godoc.org/github.com/reactivego/rx/test/Catch/ 
	Retry            https://godoc.org/github.com/reactivego/rx/test/Retry/

Utility Operators

A toolbox of useful Operators for working with Observables.

	Delay            https://godoc.org/github.com/reactivego/rx/test/Delay/
	Do               https://godoc.org/github.com/reactivego/rx/test/Do/ 
	DoOnError        https://godoc.org/github.com/reactivego/rx/test/Do/
	DoOnComplete     https://godoc.org/github.com/reactivego/rx/test/Do/
	Finally          https://godoc.org/github.com/reactivego/rx/test/Do/
	Passthrough      https://godoc.org/github.com/reactivego/rx/test/Passthrough/
	Serialize        https://godoc.org/github.com/reactivego/rx/test/Serialize/
	Timeout          https://godoc.org/github.com/reactivego/rx/test/Timeout/

Conditional and Boolean Operators

Operators that evaluate one or more Observables or items emitted by Observables.

	All              https://godoc.org/github.com/reactivego/rx/test/All/ ObservableBool

Mathematical and Aggregate Operators

Operators that operate on the entire sequence of items emitted by an Observable.

	Average           https://godoc.org/github.com/reactivego/rx/test/Average/
	Count             https://godoc.org/github.com/reactivego/rx/test/Count/ ObservableInt
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

A *Connectable* is an *Observable* that can multicast to observers subscribed to it. The *Connectable* itself will subscribe to the *Observable* when the *Connect* method is called on it.

	Publish           https://godoc.org/github.com/reactivego/rx/test/Publish/ Connectable
	PublishReplay     https://godoc.org/github.com/reactivego/rx/test/PublishReplay/ Connectable

*Connectable* supports different strategies for subscribing to the *Observable* from which it was created.

	RefCount          https://godoc.org/github.com/reactivego/rx/test/RefCount/
	AutoConnect       https://godoc.org/github.com/reactivego/rx/test/AutoConnect/
	Connect           https://godoc.org/github.com/reactivego/rx/test/Connect

Subjects

A *Subject* is both a multicasting *Observable* as well as an *Observer*. The *Observable* side allows multiple simultaneous subscribers. The *Observer* side allows you to directly feed it data or subscribe it to another *Observable*.

	Subject           https://godoc.org/github.com/reactivego/rx/test/Subject
	ReplaySubject     https://godoc.org/github.com/reactivego/rx/test/ReplaySubject

Subscribing

Subscribing breathes life into a chain of observables. An observable may be subscribed to many times. 

**Println** and **Subscribe** implement subscribing behavior directly.

	Println          https://godoc.org/github.com/reactivego/rx/test/Println
	Subscribe        https://godoc.org/github.com/reactivego/rx/test/Subscribe
	Connect          https://godoc.org/github.com/reactivego/rx/test/Connect
	ToChan           https://godoc.org/github.com/reactivego/rx/test/ToChan
	ToSingle         https://godoc.org/github.com/reactivego/rx/test/ToSingle
	ToSlice          https://godoc.org/github.com/reactivego/rx/test/ToSlice
	Wait             https://godoc.org/github.com/reactivego/rx/test/Wait

**Connect** is called internally by **RefCount** and **AutoConnect**.

	- RefCount		https://godoc.org/github.com/reactivego/rx/test/RefCount
	- AutoConnect	https://godoc.org/github.com/reactivego/rx/test/AutoConnect

Acknowledgements

This library started life as Reactive eXtensions for Go (https://github.com/alecthomas/gorx) by *Alec Thomas*.
Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present
in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

License

This library is licensed under the terms of the MIT License. See LICENSE file for copyright notice and exact wording.
*/
package rx

import _ "github.com/reactivego/rx/generic"

