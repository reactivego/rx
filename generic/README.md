# rx (generic)

    import _ "github.com/reactivego/rx/generic"

[![](../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/generic?tab=doc)
[![](../../assets/rx.svg?raw=true)](http://reactivex.io/intro.html)

Package `rx` provides *Reactive Extensions* for Go, an API for asynchronous programming with [Observables](#observables) and [Operators](#operators).

## Install
In order to use the package as a generic programming library for *Go 1*, first install the [*jig*](https://github.com/reactivego/jig) generator tool. It will generate source code from a generics library.

```bash
$ go get github.com/reactivego/jig
```
## Usage

To use the library, write a Go program that imports the library as follows:

```go
import _ "github.com/reactivego/rx/generic"
```
The `_` prefix stops Go from complaining about the imported package not being used. However, [*jig*](https://github.com/reactivego/jig) will be able to use the generics from the library via this import.

```go	
package main

import _ "github.com/reactivego/rx/generic"

func main() {
  JustString("Hello, World!").Println()
  // Output: Hello, World!
}
```
This will not build as-is, because `JustString` doesn't exist. Code for `JustString` must first be generated with the [*jig*](https://github.com/reactivego/jig) tool.

```bash
$ jig -v
found 149 templates in package "rx" (github.com/reactivego/rx/generic)
...
```
You run [*jig*](https://github.com/reactivego/jig) from the command line in the directory where your go file is stored. The [*jig*](https://github.com/reactivego/jig) tool then analyzes your code and determines what additional code is needed to make it build.

In the example above, [*jig*](https://github.com/reactivego/jig) takes templates from the `rx` library and specializes them on specific types. The generated code is written to the file `rx.go`. If all went well, your code will now build with the Go tool.

For a more in-depth look see the [Quick Start Guide](QUICKSTART.md).

## Observables

An Observable:

- is a stream of events.
- assumes zero to many values over time.
- pushes values
- can take any amount of time to complete (or may never)
- is cancellable
- is lazy (it doesn't do anything until you subscribe).

An Observer subscribes to an Observable and **reacts** to whatever it emits.

## Operators

Operators work on one or more Observables to transform, filter and combine them into new Observables.

Currently 100 operators are implemented:
    
| A … C                   | D … L                         | M … P                  | R … S              | T … W                   |
|:------------------------|:------------------------------|:-----------------------|:-------------------|:------------------------|
| [All]                   | [DebounceTime] :star:         | [Map] :star:           | [Range]            | [Take]                  |
| [AsObservable]          | [Defer]                       | [MapTo]                | [Reduce]           | [TakeLast]              |
| [AsyncSubject]          | [Delay]                       | [Max]                  | [RefCount]         | [TakeUntil]             |
| [AuditTime]             | [Distinct]                    | [Merge] :star:         | [Repeat]           | [TakeWhile]             |
| [AutoConnect]           | [DistinctUntilChanged] :star: | [MergeAll]             | [ReplaySubject]    | [ThrottleTime]          |
| [Average]               | [Do] :star:                   | [MergeDelayError]      | [Retry]            | [Throw]                 |
| [BehaviorSubject]       | [DoOnComplete]                | [MergeDelayErrorWith]  | [Sample]           | [Ticker]                |
| [Catch] :star:          | [DoOnError]                   | [MergeMap] :star:      | [Scan] :star:      | [TimeInterval]          |
| [CatchError] :star:     | [ElementAt]                   | [MergeMapTo]           | [Serialize]        | [Timeout]               |
| [CombineLatest] :star:  | [Empty]                       | [MergeWith] :star:     | [Single]           | [Timer]                 |
| [CombineLatestAll]      | [Filter] :star:               | [Min]                  | [Skip]             | [Timestamp]             |
| [CombineLatestMap]      | [Finally]                     | [Never]                | [SkipLast]         | [ToChan]                |
| [CombineLatestMapTo]    | [First]                       | [ObserveOn]            | [Start]            | [ToSingle]              |
| [CombineLatestWith]     | [From] :star:                 | [Of] :star:            | [StartWith] :star: | [ToSlice]               |
| [Concat] :star:         | [FromChan]                    | [Only]                 | [Subject]          | [Wait]                  |
| [ConcatAll]             | [IgnoreCompletion]            | [Passthrough]          | [Subscribe]        | [WithLatestFrom] :star: |
| [ConcatMap] :star:      | [IgnoreElements]              | [Println]              | [SubscribeOn]      | [WithLatestFromAll]     |
| [ConcatMapTo]           | [Interval]                    | [Publish] :star:       | [Sum]              |
| [ConcatWith] :star:     | [Just] :star:                 | [PublishReplay] :star: | [SwitchAll]        |
| [Connect]               | [Last]                        |                        | [SwitchMap] :star: 
| [Count]                 |
| [Create] :star:         |
| [CreateFutureRecursive] |
| [CreateRecursive]       |

:star: - commonly used

[All]: ../test/All/README.md
[AsObservable]: ../test/AsObservable/README.md
[AuditTime]: ../test/AuditTime/README.md
[AsyncSubject]: ../test/AsyncSubject/README.md
[AutoConnect]: ../test/AutoConnect/README.md
[Average]: ../test/Average/README.md
[BehaviorSubject]: ../test/BehaviorSubject/README.md
[Catch]: ../test/Catch/README.md
[CatchError]: ../test/CatchError/README.md
[CombineLatest]: ../test/CombineLatest/README.md
[CombineLatestAll]: ../test/CombineLatestAll/README.md
[CombineLatestMap]: ../test/CombineLatestMap/README.md
[CombineLatestMapTo]: ../test/CombineLatestMapTo/README.md
[CombineLatestWith]: ../test/CombineLatestWith/README.md
[Concat]: ../test/Concat/README.md
[ConcatAll]: ../test/ConcatAll/README.md
[ConcatMap]: ../test/ConcatMap/README.md
[ConcatMapTo]: ../test/ConcatMapTo/README.md
[ConcatWith]: ../test/ConcatWith/README.md
[Connect]: ../test/Connect/README.md
[Count]: ../test/Count/README.md
[Create]: ../test/Create/README.md
[CreateFutureRecursive]: ../test/CreateFutureRecursive/README.md
[CreateRecursive]: ../test/CreateRecursive/README.md
[DebounceTime]: ../test/DebounceTime/README.md
[Defer]: ../test/Defer/README.md
[Delay]: ../test/Delay/README.md
[Distinct]: ../test/Distinct/README.md
[DistinctUntilChanged]: ../test/DistinctUntilChanged/README.md
[Do]: ../test/Do/README.md
[DoOnComplete]: ../test/DoOnComplete/README.md
[DoOnError]: ../test/DoOnError/README.md
[ElementAt]: ../test/ElementAt/README.md
[Empty]: ../test/Empty/README.md
[Filter]: ../test/Filter/README.md
[Finally]: ../test/Finally/README.md
[First]: ../test/First/README.md
[From]: ../test/From/README.md
[FromChan]: ../test/FromChan/README.md
[IgnoreCompletion]: ../test/IgnoreCompletion/README.md
[IgnoreElements]: ../test/IgnoreElements/README.md
[Interval]: ../test/Interval/README.md
[Just]: ../test/Just/README.md
[Last]: ../test/Last/README.md
[Map]: ../test/Map/README.md
[MapTo]: ../test/MapTo/README.md
[Max]: ../test/Max/README.md
[Merge]: ../test/Merge/README.md
[MergeAll]: ../test/MergeAll/README.md
[MergeDelayError]: ../test/MergeDelayError/README.md
[MergeDelayErrorWith]: ../test/MergeDelayErrorWith/README.md
[MergeMap]: ../test/MergeMap/README.md
[MergeMapTo]: ../test/MergeMapTo/README.md
[MergeWith]: ../test/MergeWith/README.md
[Min]: ../test/Min/README.md
[Never]: ../test/Never/README.md
[ObserveOn]: ../test/ObserveOn/README.md
[Of]: ../test/Of/README.md
[Only]: ../test/Only/README.md
[Passthrough]: ../test/Passthrough/README.md
[Println]: ../test/Println/README.md
[Publish]: ../test/Publish/README.md
[PublishReplay]: ../test/PublishReplay/README.md
[Range]: ../test/Range/README.md
[Reduce]: ../test/Reduce/README.md
[RefCount]: ../test/RefCount/README.md
[Repeat]: ../test/Repeat/README.md
[ReplaySubject]: ../test/ReplaySubject/README.md
[Retry]: ../test/Retry/README.md
[Sample]: ../test/Sample/README.md
[Scan]: ../test/Scan/README.md
[Serialize]: ../test/Serialize/README.md
[Single]: ../test/Single/README.md
[Skip]: ../test/Skip/README.md
[SkipLast]: ../test/SkipLast/README.md
[Start]: ../test/Start/README.md
[StartWith]: ../test/StartWith/README.md
[Subject]: ../test/Subject/README.md
[Subscribe]: ../test/Subscribe/README.md
[SubscribeOn]: ../test/SubscribeOn/README.md
[Sum]: ../test/Sum/README.md
[SwitchAll]: ../test/SwitchAll/README.md
[SwitchMap]: ../test/SwitchMap/README.md
[Take]: ../test/Take/README.md
[TakeLast]: ../test/TakeLast/README.md
[TakeUntil]: ../test/TakeUntil/README.md
[TakeWhile]: ../test/TakeWhile/README.md
[ThrottleTime]: ../test/ThrottleTime/README.md
[Throw]: ../test/Throw/README.md
[Ticker]: ../test/Ticker/README.md
[TimeInterval]: ../test/TimeInterval/README.md
[Timeout]: ../test/Timeout/README.md
[Timer]: ../test/Timer/README.md
[Timestamp]: ../test/Timestamp/README.md
[ToChan]: ../test/ToChan/README.md
[ToSingle]: ../test/ToSingle/README.md
[ToSlice]: ../test/ToSlice/README.md
[Wait]: ../test/Wait/README.md
[WithLatestFrom]: ../test/WithLatestFrom/README.md
[WithLatestFromAll]: ../test/WithLatestFromAll/README.md
