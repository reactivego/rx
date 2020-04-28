# rx

    import "github.com/reactivego/rx"

[![](https://godoc.org/github.com/reactivego/rx?status.png)](http://godoc.org/github.com/reactivego/rx)

Package `rx` offers Reactive Extensions ([ReactiveX](http://reactivex.io/)) for [Go](https://golang.org/), an API for asynchronous programming with observable streams (Observables).

## What ReactiveX is

Observables are the main focus of ReactiveX and they are (sort of) sets of events.

- They assume zero to many values over time.
- They push values.
- They can take any amount of time to complete (or may never).
- They are cancellable.
- They are lazy; they don't do anything until you subscribe.

This implementation of rx uses `interface{}` for value types, so you can
mix different types of values in function and method calls. To create an
observable of three values with different types you might write the following:

```go
package main

import "github.com/reactivego/rx"

func main() {
	rx.From(1,"hi",2.3).Println()
}
```

The call above creates an observable from numbers and strings and then prints
them.

## Table of Contents

<!-- MarkdownTOC -->

- [Why you should use ReactiveX](#why-you-should-use-reactivex)
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
- [Obligatory Dijkstra Quote](#obligatory-dijkstra-quote)
- [Acknowledgements](#acknowledgements)
- [License](#license)

<!-- /MarkdownTOC -->

## Why you should use ReactiveX
ReactiveX observables are somewhat similar to Go channels but have much richer semantics. Observables can be hot or cold, can complete normally or with an error, use subscriptions that can be cancelled from the subscriber side. Where a normal variable is just a place where you read and write values from, an observable captures how the value of this variable changes over time. Concurrency follows naturally from the fact that an observable is an ever changing stream of values.

`rx` is a library of operators that work on one or more observables. The way in which observables can be combined using operators to form new observables is the real strength of ReactiveX. Operators specify how observables representing streams of values are e.g. merged, transformed, concatenated, split, multicasted, replayed, delayed and debounced.

This implemenation takes cues from boh [RxJS 6](https://github.com/ReactiveX/rxjs) and [RxJava 2](https://github.com/ReactiveX/RxJava) that have been pushing the envelope in evolving ReactiveX operator semantics.

## Operators 

Folowing is a list of [ReactiveX operators](http://reactivex.io/documentation/operators.html) that have been implemented.

Note, the operators that are used most commonly have been marked with a :star:.

### Creating Operators
Operators that originate new Observables.

- [**Create**](https://godoc.org/github.com/reactivego/rx/test/Create/)() :star: Observable
- [**Defer**](https://godoc.org/github.com/reactivego/rx/test/Defer/)() Observable
- [**Empty**](https://godoc.org/github.com/reactivego/rx/test/Empty/)() Observable
- [**Error**](https://godoc.org/github.com/reactivego/rx/test/Error/)() Observable
- [**FromChan**](https://godoc.org/github.com/reactivego/rx/test/From/)() Observable
- FromEventSource
- FromIterable
- FromIterator
- [**FromSlice**](https://godoc.org/github.com/reactivego/rx/test/From/)() Observable
- [**Froms**](https://godoc.org/github.com/reactivego/rx/test/From/)() Observable
- [**From**](https://godoc.org/github.com/reactivego/rx/test/From/)() :star: Observable
- [**Interval**](https://godoc.org/github.com/reactivego/rx/test/Interval/)() ObservableInt
- [**Just**](https://godoc.org/github.com/reactivego/rx/test/Just/)() :star: Observable
- [**Never**](https://godoc.org/github.com/reactivego/rx/test/Never/)() Observable
- [**Range**](https://godoc.org/github.com/reactivego/rx/test/Range/)() ObservableInt
- [**Repeat**](https://godoc.org/github.com/reactivego/rx/test/Repeat/)() Observable
- (Observable) [**Repeat**](https://godoc.org/github.com/reactivego/rx/test/Repeat/)() Observable
- [**Start**](https://godoc.org/github.com/reactivego/rx/test/Start/)() Observable
- [**Throw**](https://godoc.org/github.com/reactivego/rx/test/Throw/)() Observable

### Transforming Operators
Operators that transform items that are emitted by an Observable.

- BufferWithCount
- BufferWithTime :star:
- BufferWithTimeOrCount
- ConcatMap :star:
- (Observable) [**Map**](https://godoc.org/github.com/reactivego/rx/test/Map/)() :star: Observable
- (Observable) [**MergeMap**](https://godoc.org/github.com/reactivego/rx/test/MergeMap/)() :star: Observable
- (Observable) [**Scan**](https://godoc.org/github.com/reactivego/rx/test/Scan/)() :star: Observable
- (Observable) [**SwitchMap**](https://godoc.org/github.com/reactivego/rx/test/SwitchMap/)() :star: Observable

### Filtering Operators
Operators that selectively emit items from a source Observable.

- (Observable) [**Debounce**](https://godoc.org/github.com/reactivego/rx/test/Debounce/)() Observable
- DebounceTime :star:
- (Observable) [**Distinct**](https://godoc.org/github.com/reactivego/rx/test/Distinct/)() Observable
- DistinctUntilChanged :star:
- (Observable) [**ElementAt**](https://godoc.org/github.com/reactivego/rx/test/ElementAt/)() Observable
- (Observable) [**Filter**](https://godoc.org/github.com/reactivego/rx/test/Filter/)() :star: Observable
- (Observable) [**First**](https://godoc.org/github.com/reactivego/rx/test/First/)() Observable
- FirstOrDefault
- (Observable) [**IgnoreElements**](https://godoc.org/github.com/reactivego/rx/test/IgnoreElements/)() Observable
- (Observable) [**IgnoreCompletion**](https://godoc.org/github.com/reactivego/rx/test/IgnoreCompletion/)() Observable
- (Observable) [**Last**](https://godoc.org/github.com/reactivego/rx/test/Last/)() Observable
- LastOrDefault
- (Observable) [**Sample**](https://godoc.org/github.com/reactivego/rx/test/Sample/)() Observable
- (Observable) [**Single**](https://godoc.org/github.com/reactivego/rx/test/Single/)() Observable
- (Observable) [**Skip**](https://godoc.org/github.com/reactivego/rx/test/Skip/)() Observable
- (Observable) [**SkipLast**](https://godoc.org/github.com/reactivego/rx/test/SkipLast/)() Observable
- SkipWhile
- (Observable) [**Take**](https://godoc.org/github.com/reactivego/rx/test/Take/)() :star: Observable
- (Observable) [**TakeLast**](https://godoc.org/github.com/reactivego/rx/test/TakeLast/)() Observable
- (Observable) [**TakeUntil**](https://godoc.org/github.com/reactivego/rx/test/TakeUntil/)() :star: Observable
- (Observable) [**TakeWhile**](https://godoc.org/github.com/reactivego/rx/test/TakeWhile/)() :star: Observable

### Combining Operators
Operators that work with multiple source Observables to create a single Observable.

- (Observable<sup>2</sup>) [**ConbineAll**](https://godoc.org/github.com/reactivego/rx/test/CombineAll/)() Observable
- CombineLatest :star:
- [**Concat**](https://godoc.org/github.com/reactivego/rx/test/Concat/)() :star: Observable
- (Observable) [**Concat**](https://godoc.org/github.com/reactivego/rx/test/Concat/)() :star: Observable
- (Observable<sup>2</sup>) [**ConcatAll**](https://godoc.org/github.com/reactivego/rx/test/ConcatAll/)() Observable
- [**Merge**](https://godoc.org/github.com/reactivego/rx/test/Merge/)() Observable
- (Observable) [**Merge**](https://godoc.org/github.com/reactivego/rx/test/Merge/)() :star: Observable
- (Observable<sup>2</sup>) [**MergeAll**](https://godoc.org/github.com/reactivego/rx/test/MergeAll/)() Observable
- [**MergeDelayError**](https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/)() Observable
- (Observable) [**MergeDelayError**](https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/)() Observable
- StartWith :star:
- StartWithItems
- StartWithIterable
- StartWithObservable
- (Observable<sup>2</sup>) [**SwitchAll**](https://godoc.org/github.com/reactivego/rx/test/SwitchAll/)() Observable
- WithLatestFrom :star:
- ZipFromObservable

### Error Handling Operators
Operators that help to recover from error notifications from an Observable.

- (Observable) [**Catch**](https://godoc.org/github.com/reactivego/rx/test/Catch/)() :star: Observable
- OnErrorResumeNext
- OnErrorReturn
- OnErrorReturnItem
- (Observable) [**Retry**](https://godoc.org/github.com/reactivego/rx/test/Retry/)() Observable

### Utility Operators
A toolbox of useful Operators for working with Observables.

- (Observable) [**Delay**](https://godoc.org/github.com/reactivego/rx/test/Delay/)() Observable
- (Observable) [**Do**](https://godoc.org/github.com/reactivego/rx/test/Do/)() :star: Observable
- (Observable) [**DoOnError**](https://godoc.org/github.com/reactivego/rx/test/Do/)() Observable
- (Observable) [**DoOnComplete**](https://godoc.org/github.com/reactivego/rx/test/Do/)() Observable
- (Observable) [**Finally**](https://godoc.org/github.com/reactivego/rx/test/Do/)() Observable
- (Observable) [**Passthrough**](https://godoc.org/github.com/reactivego/rx/test/Passthrough/)() Observable
- (Observable) [**Serialize**](https://godoc.org/github.com/reactivego/rx/test/Serialize/)() Observable
- (Observable) [**Timeout**](https://godoc.org/github.com/reactivego/rx/test/Timeout/)() Observable

### Conditional and Boolean Operators
Operators that evaluate one or more Observables or items emitted by Observables.

- (Observable) [**All**](https://godoc.org/github.com/reactivego/rx/test/All/)() ObservableBool
- Amb
- Contains
- DefaultIfEmpty
- SequenceEqual

### Mathematical and Aggregate Operators
Operators that operate on the entire sequence of items emitted by an Observable.

- (Observable) [**Average**](https://godoc.org/github.com/reactivego/rx/test/Average/)() Observable
- (Observable) [**Count**](https://godoc.org/github.com/reactivego/rx/test/Count/)() ObservableInt
- (Observable) [**Max**](https://godoc.org/github.com/reactivego/rx/test/Max/)() Observable
- (Observable) [**Min**](https://godoc.org/github.com/reactivego/rx/test/Min/)() Observable
- (Observable) [**Reduce**](https://godoc.org/github.com/reactivego/rx/test/Reduce/)() Observable
- (Observable) [**Sum**](https://godoc.org/github.com/reactivego/rx/test/Sum/)() Observable

### Type Casting and Type Filtering Operators
Operators to type cast, type filter observables.

- (Observable) [**AsObservable**](https://godoc.org/github.com/reactivego/rx/test/AsObservable)() Observable
- (Observable) [**Only**](https://godoc.org/github.com/reactivego/rx/test/Only/)() Observable

### Scheduling Operators
Change the scheduler for subscribing and observing.

- (Observable) [**ObserveOn**](https://godoc.org/github.com/reactivego/rx/test/ObserveOn/)() Observable
- (Observable) [**SubscribeOn**](https://godoc.org/github.com/reactivego/rx/test/SubscribeOn/)() Observable

### Multicasting Operators
A *Connectable* is an *Observable* that can multicast to observers subscribed to it. The *Connectable* itself will subscribe to the *Observable* when the *Connect* method is called on it.

- (Observable) [**Publish**](https://godoc.org/github.com/reactivego/rx/test/Publish/)() Connectable
- (Observable) [**PublishReplay**](https://godoc.org/github.com/reactivego/rx/test/PublishReplay/)() Connectable
- PublishLast
- PublishBehavior

*Connectable* supports different strategies for subscribing to the *Observable* from which it was created.

- (Connectable) [**RefCount**](https://godoc.org/github.com/reactivego/rx/test/RefCount/)() Observable
- (Connectable) [**AutoConnect**](https://godoc.org/github.com/reactivego/rx/test/AutoConnect/)() Observable
- (Connectable) [**Connect**](https://godoc.org/github.com/reactivego/rx/test/Connect)() Subscription

## Subjects
A *Subject* is both a multicasting *Observable* as well as an *Observer*. The *Observable* side allows multiple simultaneous subscribers. The *Observer* side allows you to directly feed it data or subscribe it to another *Observable*.

- [**NewSubject**](https://godoc.org/github.com/reactivego/rx/test/Subject)() Subject
- [**NewReplaySubject**](https://godoc.org/github.com/reactivego/rx/test/ReplaySubject)() Subject

## Subscribing
Subscribing breathes life into a chain of observables. An observable may be subscribed to many times. 

**Println** and **Subscribe** implement subscribing behavior directly.

- (Observable) [**Println**](https://godoc.org/github.com/reactivego/rx/test/Println)() error
- (Observable) [**Subscribe**](https://godoc.org/github.com/reactivego/rx/test/Subscribe)() Subscription
- (Connectable) [**Connect**](https://godoc.org/github.com/reactivego/rx/test/Connect)() Subscription
- (Observable) [**ToChan**](https://godoc.org/github.com/reactivego/rx/test/ToChan)() chan foo
- (Observable) [**ToSingle**](https://godoc.org/github.com/reactivego/rx/test/ToSingle)() (foo, error)
- (Observable) [**ToSlice**](https://godoc.org/github.com/reactivego/rx/test/ToSlice)() ([]foo, error)
- (Observable) [**Wait**](https://godoc.org/github.com/reactivego/rx/test/Wait)() error


**Connect** is called internally by **RefCount** and **AutoConnect**.

- (Connectable) [**RefCount**](https://godoc.org/github.com/reactivego/rx/test/RefCount)() Observable
- (Connectable) [**AutoConnect**](https://godoc.org/github.com/reactivego/rx/test/AutoConnect)() Observable

## Obligatory Dijkstra Quote

Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed. For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.

*Edsger W. Dijkstra*, March 1968

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](../LICENSE) file in this repository for copyright notice and exact wording.

