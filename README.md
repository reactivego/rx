# rx

    import "github.com/reactivego/rx"

[![](https://godoc.org/github.com/reactivego/rx?status.png)](http://godoc.org/github.com/reactivego/rx)

Package `rx` provides [Reactive Extensions](http://reactivex.io/), an API for asynchronous programming with Observables.

## Table of Contents

<!-- MarkdownTOC -->

- [Ways of using this library](#ways-of-using-this-library)
- [Observables](#observables)
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
- [Regenerating this Package](#regenerating-this-package)
- [Obligatory Dijkstra Quote](#obligatory-dijkstra-quote)
- [Acknowledgements](#acknowledgements)
- [License](#license)

<!-- /MarkdownTOC -->

## Ways of using this library

You can use this package directly as follows:

	import "github.com/reactivego/rx"

Then you just use the code directly from the library:

	rx.From(1,2,"hello").Println()	

Alternatively you can use the library as a generics library and use
a tool to generate statically typed  observables and operators:

	import _ "github.com/reactivego/rx/generic"

Then you use the code as follows:
	
	FromInt(1,2).Println()

You'll need to generate the observables and operators by running the [jig tool](https://github.com/reactivego/jig).

## Observables

The main focus of `rx` is on [Observables](http://reactivex.io/documentation/observable.html).

An Observable:

- is a stream of events.
- assumes zero to many values over time.
- pushes values
- can take any amount of time to complete (or may never)
- is cancellable
- is lazy (it doesn't do anything until you subscribe).

This package uses `interface{}` for entry types, so an observable can emit a
mix of differently typed entries. To create an observable that emits three
values of different types you could write the following little program.

```go
package main

import "github.com/reactivego/rx"

func main() {
	rx.From(1,"hi",2.3).Println()
}
```

The code above creates an observable from numbers and strings and then prints them.

Observables in `rx` are somewhat similar to Go channels but have much richer
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

## Operators 

The combination of Observables and a set of expressive Operators is the real strength of Reactive Extensions. 
Operators work on one or more Observables. They are the language you use to describe the way in which observables should be combined to form new Observables. Operators specify how Observables representing streams of values are e.g. merged, transformed, concatenated, split, multicasted, replayed, delayed and debounced.

This implementation takes most of its cues from [RxJS 6](https://github.com/ReactiveX/rxjs) and [RxJava 2](https://github.com/ReactiveX/RxJava). Both libaries have been pushing the envelope in evolving operator semantics.

Below is the list of implemented [operators](http://reactivex.io/documentation/operators.html). Operators with a :star: are the most commonly used ones.

### Creating Operators
Operators that originate new Observables.

- [**Create**](https://godoc.org/github.com/reactivego/rx/test/Create/)() :star: Observable
- [**CreateRecursive**](https://godoc.org/github.com/reactivego/rx/test/CreateRecursive/)() Observable
- [**CreateFutureRecursive**](https://godoc.org/github.com/reactivego/rx/test/CreateFutureRecursive/)() Observable
- [**Defer**](https://godoc.org/github.com/reactivego/rx/test/Defer/)() Observable
- [**Empty**](https://godoc.org/github.com/reactivego/rx/test/Empty/)() Observable
- [**From**](https://godoc.org/github.com/reactivego/rx/test/From/)() :star: Observable
- [**FromChan**](https://godoc.org/github.com/reactivego/rx/test/FromChan/)() Observable
- [**Interval**](https://godoc.org/github.com/reactivego/rx/test/Interval/)() ObservableInt
- [**Just**](https://godoc.org/github.com/reactivego/rx/test/Just/)() :star: Observable
- [**Never**](https://godoc.org/github.com/reactivego/rx/test/Never/)() Observable
- [**Of**](https://godoc.org/github.com/reactivego/rx/test/Of/)() :star: Observable
- [**Range**](https://godoc.org/github.com/reactivego/rx/test/Range/)() ObservableInt
- [**Start**](https://godoc.org/github.com/reactivego/rx/test/Start/)() Observable
- [**Throw**](https://godoc.org/github.com/reactivego/rx/test/Throw/)() Observable
- FromEventSource(ch chan interface{}, opts ...options.Option) Observable

### Transforming Operators
Operators that transform items that are emitted by an Observable.

- (Observable) [**Map**](https://godoc.org/github.com/reactivego/rx/test/Map/)() :star: Observable
- (Observable) [**Scan**](https://godoc.org/github.com/reactivego/rx/test/Scan/)() :star: Observable
- BufferWithCount(count, skip int) Observable
- BufferWithTime(timespan, timeshift Duration) Observable
- BufferWithTimeOrCount(timespan Duration, count int) Observable

### Filtering Operators
Operators that selectively emit items from a source Observable.

- (Observable) [**Debounce**](https://godoc.org/github.com/reactivego/rx/test/Debounce/)() Observable
- (Observable) [**Distinct**](https://godoc.org/github.com/reactivego/rx/test/Distinct/)() Observable
- (Observable) [**ElementAt**](https://godoc.org/github.com/reactivego/rx/test/ElementAt/)() Observable
- (Observable) [**Filter**](https://godoc.org/github.com/reactivego/rx/test/Filter/)() :star: Observable
- (Observable) [**First**](https://godoc.org/github.com/reactivego/rx/test/First/)() Observable
- (Observable) [**IgnoreElements**](https://godoc.org/github.com/reactivego/rx/test/IgnoreElements/)() Observable
- (Observable) [**IgnoreCompletion**](https://godoc.org/github.com/reactivego/rx/test/IgnoreCompletion/)() Observable
- (Observable) [**Last**](https://godoc.org/github.com/reactivego/rx/test/Last/)() Observable
- (Observable) [**Sample**](https://godoc.org/github.com/reactivego/rx/test/Sample/)() Observable
- (Observable) [**Single**](https://godoc.org/github.com/reactivego/rx/test/Single/)() Observable
- (Observable) [**Skip**](https://godoc.org/github.com/reactivego/rx/test/Skip/)() Observable
- (Observable) [**SkipLast**](https://godoc.org/github.com/reactivego/rx/test/SkipLast/)() Observable
- (Observable) [**Take**](https://godoc.org/github.com/reactivego/rx/test/Take/)() :star: Observable
- (Observable) [**TakeLast**](https://godoc.org/github.com/reactivego/rx/test/TakeLast/)() Observable
- (Observable) [**TakeUntil**](https://godoc.org/github.com/reactivego/rx/test/TakeUntil/)() :star: Observable
- (Observable) [**TakeWhile**](https://godoc.org/github.com/reactivego/rx/test/TakeWhile/)() :star: Observable
- DebounceTime :star:
- DistinctUntilChanged(apply Function) Observable :star:
- FirstOrDefault(defaultValue interface{}) Single
- LastOrDefault(defaultValue interface{}) Single
- SkipWhile(apply Predicate) Observable

### Combining Operators
Operators that work with multiple source observables to create a single observable. It looks like there is an underlying logic at play for naming the different kinds of combining operators. The matrix below guides the naming of the operators. Where operators don't make sense we've placed a `-` in the cell. If it is not known yet whether an operator makes sense, a `?` is placed in the cell.

| Function | Operator | Map | MapTo | All |
| -------- | -------- | --- | ----- | --- |
| [**Concat**](https://godoc.org/github.com/reactivego/rx/test/Concat/) :star: | [**ConcatWith**](https://godoc.org/github.com/reactivego/rx/test/ConcatWith/) :star: | [**ConcatMap**](https://godoc.org/github.com/reactivego/rx/test/ConcatMap/) | [**ConcatMapTo**](https://godoc.org/github.com/reactivego/rx/test/ConcatMapTo/) | [**ConcatAll**](https://godoc.org/github.com/reactivego/rx/test/ConcatAll/) |
| -             | -                 | [**SwitchMap**](https://godoc.org/github.com/reactivego/rx/test/SwitchMap/) :star: | SwitchMapTo        | [**SwitchAll**](https://godoc.org/github.com/reactivego/rx/test/SwitchAll/) |
| -             | -                 | ExhaustMap       | ExhaustMapTo       | ExhaustAll       |
| [**Merge**](https://godoc.org/github.com/reactivego/rx/test/Merge/) :star: | [**MergeWith**](https://godoc.org/github.com/reactivego/rx/test/MergeWith/) :star: | [**MergeMap**](https://godoc.org/github.com/reactivego/rx/test/MergeMap/) :star: | MergeMapTo         | [**MergeAll**](https://godoc.org/github.com/reactivego/rx/test/MergeAll/) |
| [**MergeDelayError**](https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/) | [**MergeDelayErrorWith**](https://godoc.org/github.com/reactivego/rx/test/MergeDelayErrorWith/) | MergeDelayErrorMap         | MergeDelayErrorMapTo         | MergeDelayErrorAll         |
| Race          | RaceWith          | RaceMap          | RaceMapTo          | RaceAll          |
| [**CombineLatest**](https://godoc.org/github.com/reactivego/rx/test/CombineLatest/) | [**CombineLatestWith**](https://godoc.org/github.com/reactivego/rx/test/CombineLatestWith/) | [**CombineLatestMap**](https://godoc.org/github.com/reactivego/rx/test/CombineLatestMap/) | [**CombineLatestMapTo**](https://godoc.org/github.com/reactivego/rx/test/CombineLatestMapTo/) | [**CombineLatestAll**](https://godoc.org/github.com/reactivego/rx/test/CombineLatestAll/) |
| Zip           | ZipWith           | ZipMap           | ZipMapTo           | ZipAll           |
| ?             | WithLatestFrom :star:| ?                | ?                  | ?                |
| ForkJoin      | ?                 | ?                | ?                  | ?                |

**Concat**, **Switch** and **Exhaust** are all operators that flatten the emissions of multiple observables into a single stream by subscribing to every observable *stricly in sequence*. Observables may be added while the flattening is already going on.

**Merge**, **MergeDelayError** and **Race** are operators that flatten the emissions of multiple observables into a single stream by subscribing to all observables *concurrently*. Here also, observables may be added while the flattening is already going on.

**CombineLatest**, **Zip**, **WithLatestFrom** and **ForkJoin** are operators that flatten the emissions of multiple observables into a single observable that emits slices of values. Differently from the previous two sets of operators, these operators only start emitting once the list of observables to flatten is complete. 

### Error Handling Operators
Operators that help to recover from error notifications from an Observable.

- (Observable) [**Catch**](https://godoc.org/github.com/reactivego/rx/test/Catch/)() :star: Observable
- (Observable) [**Retry**](https://godoc.org/github.com/reactivego/rx/test/Retry/)() Observable
- OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable
- OnErrorReturn(resumeFunc ErrorFunction) Observable
- OnErrorReturnItem(item interface{}) Observable

### Utility Operators
A toolbox of useful Operators for working with Observables.

- (Observable) [**Delay**](https://godoc.org/github.com/reactivego/rx/test/Delay/)() Observable
- (Observable) [**Do**](https://godoc.org/github.com/reactivego/rx/test/Do/)() :star: Observable
- (Observable) [**DoOnError**](https://godoc.org/github.com/reactivego/rx/test/DoOnError/)() Observable
- (Observable) [**DoOnComplete**](https://godoc.org/github.com/reactivego/rx/test/DoOnComplete/)() Observable
- (Observable) [**Finally**](https://godoc.org/github.com/reactivego/rx/test/Finally/)() Observable
- (Observable) [**Passthrough**](https://godoc.org/github.com/reactivego/rx/test/Passthrough/)() Observable
- (Observable) [**Repeat**](https://godoc.org/github.com/reactivego/rx/test/Repeat/)() Observable
- (Observable) [**Serialize**](https://godoc.org/github.com/reactivego/rx/test/Serialize/)() Observable
- (Observable) [**Timeout**](https://godoc.org/github.com/reactivego/rx/test/Timeout/)() Observable
- Repeat(count int64, frequency Duration) Observable
- StartWith :star:
- StartWithItems(item interface{}, items ...interface{}) Observable
- StartWithObservable(observable Observable) Observable
- EndWith

### Conditional and Boolean Operators
Operators that evaluate one or more Observables or items emitted by Observables.

- (Observable) [**All**](https://godoc.org/github.com/reactivego/rx/test/All/)() ObservableBool
- Amb(observable Observable, observables ...Observable) Observable
- Contains(equal Predicate) Single
- DefaultIfEmpty(defaultValue interface{}) Observable
- SequenceEqual(obs Observable) Single

### Aggregate Operators
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

- [**Subject**](https://godoc.org/github.com/reactivego/rx/test/Subject)() Subject
- [**ReplaySubject**](https://godoc.org/github.com/reactivego/rx/test/ReplaySubject)() Subject

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

## Regenerating this Package

This package is generated from the sub-folder generic by the [jig](http://github.com/reactivego/jig) tool.
You don't need to regenerate the package in order to use it. However, if you are
interested in regenerating it, then read on.

The jig tool provides the parametric polymorphism capability that Go 1 is missing.
It works by replacing place-holder types of generic functions and datatypes
with interface{} (it can also generate statically typed code though).

To regenerate, change the current working directory to the package directory
and run the jig tool as follows:

```bash
$ go get -d github.com/reactivego/jig
$ go run github.com/reactivego/jig -v
```

## Obligatory Dijkstra Quote

Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed. For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.

*Edsger W. Dijkstra*, March 1968

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file for copyright notice and exact wording.
