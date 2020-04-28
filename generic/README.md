# rx

    import _ "github.com/reactivego/rx/generic"

[![](https://godoc.org/github.com/reactivego/rx?status.png)](http://godoc.org/github.com/reactivego/rx/generic)

Package `rx` offers [Reactive eXtensions](http://reactivex.io/) for [Go](https://golang.org/). It's a generics library for composing asynchronous and event-based programs using observable sequences. The library consists of more than a 100 templates to enable type-safe programming with observable streams. To use it, you will need the *jig* tool from [Just-in-time Generics for Go](https://github.com/reactivego/jig).

Using the library is very simple. Import the library with the blank identifier `_` as the package name. The side effect of this import is that generics from the library can now be accessed by the *jig* tool. Then start using generics from the library and run *jig* to generate code. The following is a minimal *Hello, World!* program:

```go
package main

import _ "github.com/reactivego/rx/jig"

func main() {
	FromStrings("You!", "Gophers!", "World!").
		MapString(func(x string) string {
			return "Hello, " + x
		}).Println()

	// Output:
	// Hello, You!
	// Hello, Gophers!
	// Hello, World!
}
```

## Table of Contents

<!-- MarkdownTOC -->

- [Why?](#why)
- [Operators](#operators)
	- [Creating Operators](#creating-operators)
	- [Transforming Operators](#transforming-operators)
	- [Filtering Operators](#filtering-operators)
	- [Combining Operators](#combining-operators)
	- [Error Handling Operators](#error-handling-operators)
	- [Utility Operators](#utility-operators)
	- [Conditional and Boolean Operators](#conditional-and-boolean-operators)
	- [Mathematical and Aggregate Operators](#mathematical-and-aggregate-operators)
	- [Type Casting, Type Filtering Operators](#type-casting-and-type-filtering-operators)
	- [Scheduling Operators](#scheduling-operators)
	- [Multicasting Operators](#multicasting-operators)
- [Subjects](#subjects)
- [Subscribing](#subscribing)
- [Generic Programming](#generic-programming)
- [Obligatory Dijkstra Quote](#obligatory-dijkstra-quote)
- [Acknowledgements](#acknowledgements)
- [License](#license)

<!-- /MarkdownTOC -->

## Why?
ReactiveX observables are somewhat similar to Go channels but have much richer semantics. Observables can be hot or cold, can complete normally or with an error, use subscriptions that can be cancelled from the subscriber side. Where a normal variable is just a place where you read and write values from, an observable captures how the value of this variable changes over time. Concurrency follows naturally from the fact that an observable is an ever changing stream of values.

`rx` is a library of operators that work on one or more observables. The way in which observables can be combined using operators to form new observables is the real strength of ReactiveX. Operators specify how observables representing streams of values are e.g. merged, transformed, concatenated, split, multicasted, replayed, delayed and debounced.

This implemenation takes cues from boh [RxJS 6](https://github.com/ReactiveX/rxjs) and [RxJava 2](https://github.com/ReactiveX/RxJava) that have been pushing the envelope in evolving ReactiveX operator semantics.

## Operators 

Folowing is a list of [ReactiveX operators](http://reactivex.io/documentation/operators.html) that have been implemented. Operators that are most commonly used got a :star:.

### Creating Operators
Operators that originate new Observables.

- [**Create**](https://godoc.org/github.com/reactivego/rx/generic/test/Create/)() :star: Observable
- [**Defer**](https://godoc.org/github.com/reactivego/rx/generic/test/Defer/)() Observable
- [**Empty**](https://godoc.org/github.com/reactivego/rx/generic/test/Empty/)() Observable
- [**Error**](https://godoc.org/github.com/reactivego/rx/generic/test/Error/)() Observable
- [**FromChan**](https://godoc.org/github.com/reactivego/rx/generic/test/From/)() Observable
- FromEventSource
- FromIterable
- FromIterator
- [**FromSlice**](https://godoc.org/github.com/reactivego/rx/generic/test/From/)() Observable
- [**Froms**](https://godoc.org/github.com/reactivego/rx/generic/test/From/)() Observable
- [**From**](https://godoc.org/github.com/reactivego/rx/generic/test/From/)() :star: Observable
- [**Interval**](https://godoc.org/github.com/reactivego/rx/generic/test/Interval/)() ObservableInt
- [**Just**](https://godoc.org/github.com/reactivego/rx/generic/test/Just/)() :star: Observable
- [**Never**](https://godoc.org/github.com/reactivego/rx/generic/test/Never/)() Observable
- [**Range**](https://godoc.org/github.com/reactivego/rx/generic/test/Range/)() ObservableInt
- [**Repeat**](https://godoc.org/github.com/reactivego/rx/generic/test/Repeat/)() Observable
- (Observable) [**Repeat**](https://godoc.org/github.com/reactivego/rx/generic/test/Repeat/)() Observable
- [**Start**](https://godoc.org/github.com/reactivego/rx/generic/test/Start/)() Observable
- [**Throw**](https://godoc.org/github.com/reactivego/rx/generic/test/Throw/)() Observable

### Transforming Operators
Operators that transform items that are emitted by an Observable.

- BufferWithCount
- BufferWithTime :star:
- BufferWithTimeOrCount
- ConcatMap :star:
- (Observable) [**Map**](https://godoc.org/github.com/reactivego/rx/generic/test/Map/)() :star: Observable
- (Observable) [**MergeMap**](https://godoc.org/github.com/reactivego/rx/generic/test/MergeMap/)() :star: Observable
- (Observable) [**Scan**](https://godoc.org/github.com/reactivego/rx/generic/test/Scan/)() :star: Observable
- (Observable) [**SwitchMap**](https://godoc.org/github.com/reactivego/rx/generic/test/SwitchMap/)() :star: Observable

### Filtering Operators
Operators that selectively emit items from a source Observable.

- (Observable) [**Debounce**](https://godoc.org/github.com/reactivego/rx/generic/test/Debounce/)() Observable
- DebounceTime :star:
- (Observable) [**Distinct**](https://godoc.org/github.com/reactivego/rx/generic/test/Distinct/)() Observable
- DistinctUntilChanged :star:
- (Observable) [**ElementAt**](https://godoc.org/github.com/reactivego/rx/generic/test/ElementAt/)() Observable
- (Observable) [**Filter**](https://godoc.org/github.com/reactivego/rx/generic/test/Filter/)() :star: Observable
- (Observable) [**First**](https://godoc.org/github.com/reactivego/rx/generic/test/First/)() Observable
- FirstOrDefault
- (Observable) [**IgnoreElements**](https://godoc.org/github.com/reactivego/rx/generic/test/IgnoreElements/)() Observable
- (Observable) [**IgnoreCompletion**](https://godoc.org/github.com/reactivego/rx/generic/test/IgnoreCompletion/)() Observable
- (Observable) [**Last**](https://godoc.org/github.com/reactivego/rx/generic/test/Last/)() Observable
- LastOrDefault
- (Observable) [**Sample**](https://godoc.org/github.com/reactivego/rx/generic/test/Sample/)() Observable
- (Observable) [**Single**](https://godoc.org/github.com/reactivego/rx/generic/test/Single/)() Observable
- (Observable) [**Skip**](https://godoc.org/github.com/reactivego/rx/generic/test/Skip/)() Observable
- (Observable) [**SkipLast**](https://godoc.org/github.com/reactivego/rx/generic/test/SkipLast/)() Observable
- SkipWhile
- (Observable) [**Take**](https://godoc.org/github.com/reactivego/rx/generic/test/Take/)() :star: Observable
- (Observable) [**TakeLast**](https://godoc.org/github.com/reactivego/rx/generic/test/TakeLast/)() Observable
- (Observable) [**TakeUntil**](https://godoc.org/github.com/reactivego/rx/generic/test/TakeUntil/)() :star: Observable
- (Observable) [**TakeWhile**](https://godoc.org/github.com/reactivego/rx/generic/test/TakeWhile/)() :star: Observable

### Combining Operators
Operators that work with multiple source Observables to create a single Observable.

- (Observable<sup>2</sup>) [**ConbineAll**](https://godoc.org/github.com/reactivego/rx/generic/test/CombineAll/)() Observable
- CombineLatest :star:
- [**Concat**](https://godoc.org/github.com/reactivego/rx/generic/test/Concat/)() :star: Observable
- (Observable) [**Concat**](https://godoc.org/github.com/reactivego/rx/generic/test/Concat/)() :star: Observable
- (Observable<sup>2</sup>) [**ConcatAll**](https://godoc.org/github.com/reactivego/rx/generic/test/ConcatAll/)() Observable
- [**Merge**](https://godoc.org/github.com/reactivego/rx/generic/test/Merge/)() Observable
- (Observable) [**Merge**](https://godoc.org/github.com/reactivego/rx/generic/test/Merge/)() :star: Observable
- (Observable<sup>2</sup>) [**MergeAll**](https://godoc.org/github.com/reactivego/rx/generic/test/MergeAll/)() Observable
- [**MergeDelayError**](https://godoc.org/github.com/reactivego/rx/generic/test/MergeDelayError/)() Observable
- (Observable) [**MergeDelayError**](https://godoc.org/github.com/reactivego/rx/generic/test/MergeDelayError/)() Observable
- StartWith :star:
- StartWithItems
- StartWithIterable
- StartWithObservable
- (Observable<sup>2</sup>) [**SwitchAll**](https://godoc.org/github.com/reactivego/rx/generic/test/SwitchAll/)() Observable
- WithLatestFrom :star:
- ZipFromObservable

### Error Handling Operators
Operators that help to recover from error notifications from an Observable.

- (Observable) [**Catch**](https://godoc.org/github.com/reactivego/rx/generic/test/Catch/)() :star: Observable
- OnErrorResumeNext
- OnErrorReturn
- OnErrorReturnItem
- (Observable) [**Retry**](https://godoc.org/github.com/reactivego/rx/generic/test/Retry/)() Observable

### Utility Operators
A toolbox of useful Operators for working with Observables.

- (Observable) [**Delay**](https://godoc.org/github.com/reactivego/rx/generic/test/Delay/)() Observable
- (Observable) [**Do**](https://godoc.org/github.com/reactivego/rx/generic/test/Do/)() :star: Observable
- (Observable) [**DoOnError**](https://godoc.org/github.com/reactivego/rx/generic/test/Do/)() Observable
- (Observable) [**DoOnComplete**](https://godoc.org/github.com/reactivego/rx/generic/test/Do/)() Observable
- (Observable) [**Finally**](https://godoc.org/github.com/reactivego/rx/generic/test/Do/)() Observable
- (Observable) [**Passthrough**](https://godoc.org/github.com/reactivego/rx/generic/test/Passthrough/)() Observable
- (Observable) [**Serialize**](https://godoc.org/github.com/reactivego/rx/generic/test/Serialize/)() Observable
- (Observable) [**Timeout**](https://godoc.org/github.com/reactivego/rx/generic/test/Timeout/)() Observable

### Conditional and Boolean Operators
Operators that evaluate one or more Observables or items emitted by Observables.

- (Observable) [**All**](https://godoc.org/github.com/reactivego/rx/generic/test/All/)() ObservableBool
- Amb
- Contains
- DefaultIfEmpty
- SequenceEqual

### Mathematical and Aggregate Operators
Operators that operate on the entire sequence of items emitted by an Observable.

- (Observable) [**Average**](https://godoc.org/github.com/reactivego/rx/generic/test/Average/)() Observable
- (Observable) [**Count**](https://godoc.org/github.com/reactivego/rx/generic/test/Count/)() ObservableInt
- (Observable) [**Max**](https://godoc.org/github.com/reactivego/rx/generic/test/Max/)() Observable
- (Observable) [**Min**](https://godoc.org/github.com/reactivego/rx/generic/test/Min/)() Observable
- (Observable) [**Reduce**](https://godoc.org/github.com/reactivego/rx/generic/test/Reduce/)() Observable
- (Observable) [**Sum**](https://godoc.org/github.com/reactivego/rx/generic/test/Sum/)() Observable

### Type Casting and Type Filtering Operators
Operators to type cast, type filter observables.

- (Observable) [**AsObservable**](https://godoc.org/github.com/reactivego/rx/generic/test/AsObservable)() Observable
- (Observable) [**Only**](https://godoc.org/github.com/reactivego/rx/generic/test/Only/)() Observable

### Scheduling Operators
Change the scheduler for subscribing and observing.

- (Observable) [**ObserveOn**](https://godoc.org/github.com/reactivego/rx/generic/test/ObserveOn/)() Observable
- (Observable) [**SubscribeOn**](https://godoc.org/github.com/reactivego/rx/generic/test/SubscribeOn/)() Observable

### Multicasting Operators
A *Connectable* is an *Observable* that can multicast to observers subscribed to it. The *Connectable* itself will subscribe to the *Observable* when the *Connect* method is called on it.

- (Observable) [**Publish**](https://godoc.org/github.com/reactivego/rx/generic/test/Publish/)() Connectable
- (Observable) [**PublishReplay**](https://godoc.org/github.com/reactivego/rx/generic/test/PublishReplay/)() Connectable
- PublishLast
- PublishBehavior

*Connectable* supports different strategies for subscribing to the *Observable* from which it was created.

- (Connectable) [**RefCount**](https://godoc.org/github.com/reactivego/rx/generic/test/RefCount/)() Observable
- (Connectable) [**AutoConnect**](https://godoc.org/github.com/reactivego/rx/generic/test/AutoConnect/)() Observable
- (Connectable) [**Connect**](https://godoc.org/github.com/reactivego/rx/generic/test/Connect)() Subscription

## Subjects
A *Subject* is both a multicasting *Observable* as well as an *Observer*. The *Observable* side allows multiple simultaneous subscribers. The *Observer* side allows you to directly feed it data or subscribe it to another *Observable*.

- [**NewSubject**](https://godoc.org/github.com/reactivego/rx/generic/test/Subject)() Subject
- [**NewReplaySubject**](https://godoc.org/github.com/reactivego/rx/generic/test/ReplaySubject)() Subject

## Subscribing
Subscribing breathes life into a chain of observables. An observable may be subscribed to many times. 

**Println** and **Subscribe** implement subscribing behavior directly.

- (Observable) [**Println**](https://godoc.org/github.com/reactivego/rx/generic/test/Println)() error
- (Observable) [**Subscribe**](https://godoc.org/github.com/reactivego/rx/generic/test/Subscribe)() Subscription
- (Connectable) [**Connect**](https://godoc.org/github.com/reactivego/rx/generic/test/Connect)() Subscription
- (Observable) [**ToChan**](https://godoc.org/github.com/reactivego/rx/generic/test/ToChan)() chan foo
- (Observable) [**ToSingle**](https://godoc.org/github.com/reactivego/rx/generic/test/ToSingle)() (foo, error)
- (Observable) [**ToSlice**](https://godoc.org/github.com/reactivego/rx/generic/test/ToSlice)() ([]foo, error)
- (Observable) [**Wait**](https://godoc.org/github.com/reactivego/rx/generic/test/Wait)() error


**Connect** is called internally by **RefCount** and **AutoConnect**.

- (Connectable) [**RefCount**](https://godoc.org/github.com/reactivego/rx/generic/test/RefCount)() Observable
- (Connectable) [**AutoConnect**](https://godoc.org/github.com/reactivego/rx/generic/test/AutoConnect)() Observable

## Generic Programming

`rx` is a generics library. Generics like `Map<T>` use a place-holder type name `T` between angle brackets. Go does not natively support this syntax. Instead I use the syntax `MapT` that Go does support. A special comment prefix then specifies what part of the name is generic. See the example below for what a template written in this syntax will look like:

```go
//jig:template Observable<Foo> Map<Bar>

func (o ObservableFoo) MapBar(project func(foo) bar) ObservableBar {
	...
	// generic implementation using `foo` and `bar` for the real type and `Foo` and `Bar` in identifiers
	...
}
```
I have choosen to use metasyntactic type names like *Foo* and *Bar* for my template libraries. So e.g. `MapBar` instead of `Map<T>`. I settled on using the words *Foo*, *Bar* and *Baz*, but you can choose other words you like better.

Using generics is easy. Just reference the generic in your code and specify a concrete type instead of the place-holder type.
So, to specialize on `int` write e.g. `MapInt`. See the following code for how that works: 

```go
FromInts(1, 2).MapInt(func(x int) int {
	return x + 10
}).Println()

// Output:
// 11
// 12
```
This code will not compile by itself, because MapInt is not known. Running the [jig](https://github.com/reactivego/jig) command will specialize all generic templates into compilable code.

## Obligatory Dijkstra Quote

Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed. For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.

*Edsger W. Dijkstra*, March 1968

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](../LICENSE) file in this repository for copyright notice and exact wording.
