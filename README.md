# rx

    import "github.com/reactivego/rx"

[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx)
[![](svg/rx.svg)](http://reactivex.io/intro.html)

Package `rx` provides *Reactive Extensions* for Go, an API for asynchronous programming with observable streams.

> Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed.
> For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.
>
> *Edsger W. Dijkstra*, March 1968

### Installation
Use the go tool to get the package:

```bash
$ go get github.com/reactivego/rx
```

Then import it in your program:

```go
import "github.com/reactivego/rx"
```
## Usage
The API is accessed directly through the `rx` prefix.
```go
rx.From(1,2,"hello").Println()
```
For documentation and examples, see the online [go.dev](https://pkg.go.dev/github.com/reactivego/rx?tab=doc) or [godoc](http://godoc.org/github.com/reactivego/rx) reference.

> NOTE : this package can also be used as a *Generics Library* for *Go 1* see [below](#generics-library)

## Observables
The main focus of `rx` is on [Observables](http://reactivex.io/documentation/observable.html).
In accordance with Dijkstra's observation that we are relatively bad at grasping dynamic processes evolving in time, Observables indeed shorten the conceptual gap between the static program and the dynamic process. Observables have dynamic behavior folded into their core as they are defined as values changing over time. Observables are combined using (lots of different) operators that compose well together. At the program level you are no longer dealing explicitly with dynamic processes, that is all defined and implied by the operators that were used to construct the program.

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

Concurrency flows naturally from the fact that an observable is an ever
changing stream of values. Every Observable conceptually has at its core a
concurrently running process that pushes out values.

## Operators 
Operators form a language in which programs featuring Observables can be expressed.
They work on one or more Observables to transform, filter and combine them into new Observables.

Currently 90 operators have been implemented:

| A … C                | C … I                   | I … P             | P … S           | S … W          |
|:---------------------|:------------------------|:------------------|:----------------|:---------------|
| [All]                | [Create]                | [IgnoreElements]  | [Publish]       | [Sum]          |
| [AsObservable]       | [CreateFutureRecursive] | [Interval]        | [PublishReplay] | [SwitchAll]    |
| [Audit]              | [CreateRecursive]       | [Just]            | [Range]         | [SwitchMap]    |
| [AutoConnect]        | [Debounce]              | [Last]            | [Reduce]        | [Take]         |
| [Average]            | [Defer]                 | [Map]             | [RefCount]      | [TakeLast]     |
| [Catch]              | [Delay]                 | [Max]             | [Repeat]        | [TakeUntil]    |
| [CombineLatest]      | [Distinct]              | [Merge]           | [ReplaySubject] | [TakeWhile]    |
| [CombineLatestAll]   | [Do]                    | [MergeAll]        | [Retry]         | [Throttle]     |
| [CombineLatestMap]   | [DoOnComplete]          | [MergeDelayError] | [Sample]        | [Throw]        |
| [CombineLatestMapTo] | [DoOnError]             | [MergeMap]        | [Scan]          | [Ticker]       |
| [CombineLatestWith]  | [ElementAt]             | [MergeWith]       | [Serialize]     | [TimeInterval] |
| [Concat]             | [Empty]                 | [Min]             | [Single]        | [Timeout]      |
| [ConcatAll]          | [Filter]                | [Never]           | [Skip]          | [Timer]        |
| [ConcatMap]          | [Finally]               | [ObserveOn]       | [SkipLast]      | [Timestamp]    |
| [ConcatMapTo]        | [First]                 | [Of]              | [Start]         | [ToChan]       |
| [ConcatWith]         | [From]                  | [Only]            | [Subject]       | [ToSingle]     |
| [Connect]            | [FromChan]              | [Passthrough]     | [Subscribe]     | [ToSlice]      |
| [Count]              | [IgnoreCompletion]      | [Println]         | [SubscribeOn]   | [Wait]         |
> See [test directory](test) for Operators details 

## Regenerating this Package
This package is generated from generics in the sub-folder generic by the [jig](http://github.com/reactivego/jig) tool.
You don't need to regenerate this package in order to use it. However, if you are
interested in regenerating it, then read on.

The [jig](http://github.com/reactivego/jig) tool provides the parametric polymorphism capability that Go 1 is missing.
It works by replacing place-holder types of generic functions and datatypes
with `interface{}` (it can also generate statically typed code though).

To regenerate, change the current working directory to the package directory
and run the [jig](http://github.com/reactivego/jig) tool as follows:

```bash
$ go get -d github.com/reactivego/jig
$ go run github.com/reactivego/jig -v
```
## Generics Library
This package can be used as a *Generics Library* for *Go 1*. It supports the writing of programs that use statically typed *Observables* and *Operators*. For example:

```go
package main

import _ "github.com/reactivego/rx"

func main() {
	FromString("Hello", "Gophers!").Println()
}
```
> Note that From**String** is statically typed.

Instructions on how to use the package this way, is in the [GENERIC](GENERIC.md) document.

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

This implementation takes most of its cues from [RxJS](https://github.com/ReactiveX/rxjs).
It is the ReactiveX incarnation that pushes the envelope in evolving operator semantics.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file for copyright notice and exact wording.

[All]: test/All
[All]: test/All
[AsObservable]: test/AsObservable
[Audit]: test/Audit
[AutoConnect]: test/AutoConnect
[Average]: test/Average
[Catch]: test/Catch
[CombineLatest]: test/CombineLatest
[CombineLatestAll]: test/CombineLatestAll
[CombineLatestMap]: test/CombineLatestMap
[CombineLatestMapTo]: test/CombineLatestMapTo
[CombineLatestWith]: test/CombineLatestWith
[Concat]: test/Concat
[ConcatAll]: test/ConcatAll
[ConcatMap]: test/ConcatMap
[ConcatMapTo]: test/ConcatMapTo
[ConcatWith]: test/ConcatWith
[Connect]: test/Connect
[Count]: test/Count
[Create]: test/Create
[CreateFutureRecursive]: test/CreateFutureRecursive
[CreateRecursive]: test/CreateRecursive
[Debounce]: test/Debounce
[Defer]: test/Defer
[Delay]: test/Delay
[Distinct]: test/Distinct
[Do]: test/Do
[DoOnComplete]: test/DoOnComplete
[DoOnError]: test/DoOnError
[ElementAt]: test/ElementAt
[Empty]: test/Empty
[Filter]: test/Filter
[Finally]: test/Finally
[First]: test/First
[From]: test/From
[FromChan]: test/FromChan
[IgnoreCompletion]: test/IgnoreCompletion
[IgnoreElements]: test/IgnoreElements
[Interval]: test/Interval
[Just]: test/Just
[Last]: test/Last
[Map]: test/Map
[Max]: test/Max
[Merge]: test/Merge
[MergeAll]: test/MergeAll
[MergeDelayError]: test/MergeDelayError
[MergeMap]: test/MergeMap
[MergeWith]: test/MergeWith
[Min]: test/Min
[Never]: test/Never
[ObserveOn]: test/ObserveOn
[Of]: test/Of
[Only]: test/Only
[Passthrough]: test/Passthrough
[Println]: test/Println
[Publish]: test/Publish
[PublishReplay]: test/PublishReplay
[Range]: test/Range
[Reduce]: test/Reduce
[RefCount]: test/RefCount
[Repeat]: test/Repeat
[ReplaySubject]: test/ReplaySubject
[Retry]: test/Retry
[Sample]: test/Sample
[Scan]: test/Scan
[Serialize]: test/Serialize
[Single]: test/Single
[Skip]: test/Skip
[SkipLast]: test/SkipLast
[Start]: test/Start
[Subject]: test/Subject
[Subscribe]: test/Subscribe
[SubscribeOn]: test/SubscribeOn
[Sum]: test/Sum
[SwitchAll]: test/SwitchAll
[SwitchMap]: test/SwitchMap
[Take]: test/Take
[TakeLast]: test/TakeLast
[TakeUntil]: test/TakeUntil
[TakeWhile]: test/TakeWhile
[Throttle]: test/Throttle
[Throw]: test/Throw
[Ticker]: test/Ticker
[TimeInterval]: test/TimeInterval
[Timeout]: test/Timeout
[Timer]: test/Timer
[Timestamp]: test/Timestamp
[ToChan]: test/ToChan
[ToSingle]: test/ToSingle
[ToSlice]: test/ToSlice
[Wait]: test/Wait
