# rx (generic)

    import _ "github.com/reactivego/rx/generic"

[![](../../assets/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/generic?tab=doc)
[![](../../assets/godoc.svg)](https://godoc.org/github.com/reactivego/rx/generic)
[![](../../assets/rx.svg)](http://reactivex.io/intro.html)

Package `rx` provides *Reactive Extensions* for Go, an API for asynchronous programming with observable streams.

Observables are combined with [operators](#operators) to form a program.

## Install
In order to use the package as a generic programming library for *Go 1*, first install the [*jig*](https://github.com/reactivego/jig) generator tool. It will generate source code from a generics library.

```bash
$ go get github.com/reactivego/jig
```

Installing the actual library is very simple. Just import the library with the blank identifier `_` in a Go program.

```go
import _ "github.com/reactivego/rx/generic"
```
The use of `_` stops the *go* tool from complaining about you not using any code from the package while at the same time allowing the [*jig*](https://github.com/reactivego/jig) tool to actually find the generics in the library.

## Usage

Write a program, placing the code in a file (e.g. `main.go`). In the code, import the `rx/generic` package.

```go	
package main

import _ "github.com/reactivego/rx/generic"

func main() {
  JustString("Hello, World!").Println()
  // Output: Hello, World!
}
```
This will not build as-is, because `JustString` doesn't exist.
Code for `JustString` must first be generated with the [*jig*](https://github.com/reactivego/jig) tool.

1. You run the [jig](https://github.com/reactivego/jig) command in the directory where your program source is located.
2. **Now [jig](https://github.com/reactivego/jig) analyzes your source code and determines what additional code is needed to make it build**.
3. *[jig](https://github.com/reactivego/jig)* takes templates from the `rx` library and specializes them on specific types.
4. Specializations are generated into the file `rx.go` alongside your own code.
5. If all went well, your code will now build.

This was a bird's-eye view of the whole process. Look in the [Quick Start Guide](QUICKSTART.md) for a more detailed explanation.

## Operators

[![](../../assets/rx.svg)](http://reactivex.io/documentation/operators.html) 

Currently 94 operators have been implemented:
    
| A … C                   | D … L                 | M … P                  | R … S              | T … W                   |
|:------------------------|:----------------------|:-----------------------|:-------------------|:------------------------|
| [All]                   | [DebounceTime] :star: | [Map] :star:           | [Range]            | [Take]                  |
| [AsObservable]          | [Defer]               | [Max]                  | [Reduce]           | [TakeLast]              |
| [AsyncSubject]          | [Delay]               | [Merge] :star:         | [RefCount]         | [TakeUntil]             |
| [AuditTime]             | [Distinct]            | [MergeAll]             | [Repeat]           | [TakeWhile]             |
| [AutoConnect]           | [Do] :star:           | [MergeDelayError]      | [ReplaySubject]    | [ThrottleTime]          |
| [Average]               | [DoOnComplete]        | [MergeDelayErrorWith]  | [Retry]            | [Throw]                 |
| [BehaviorSubject]       | [DoOnError]           | [MergeMap] :star:      | [Sample]           | [Ticker]                |
| [Catch] :star:          | [ElementAt]           | [MergeMapTo]           | [Scan] :star:      | [TimeInterval]          |
| [CatchError] :star:     | [Empty]               | [MergeWith] :star:     | [Serialize]        | [Timeout]               |
| [CombineLatest] :star:  | [Filter] :star:       | [Min]                  | [Single]           | [Timer]                 |
| [CombineLatestAll]      | [Finally]             | [Never]                | [Skip]             | [Timestamp]             |
| [CombineLatestMap]      | [First]               | [ObserveOn]            | [SkipLast]         | [ToChan]                |
| [CombineLatestMapTo]    | [From] :star:         | [Of] :star:            | [Start]            | [ToSingle]              |
| [CombineLatestWith]     | [FromChan]            | [Only]                 | [Subject]          | [ToSlice]               |
| [Concat] :star:         | [IgnoreCompletion]    | [Passthrough]          | [Subscribe]        | [Wait]                  |
| [ConcatAll]             | [IgnoreElements]      | [Println]              | [SubscribeOn]      | [WithLatestFrom] :star: |
| [ConcatMap] :star:      | [Interval]            | [Publish] :star:       | [Sum]              | [WithLatestFromAll]     |
| [ConcatMapTo]           | [Just] :star:         | [PublishReplay] :star: | [SwitchAll]        |
| [ConcatWith] :star:     | [Last]                |                        | [SwitchMap] :star: |
| [Connect]               |
| [Count]                 |
| [Create] :star:         |
| [CreateFutureRecursive] |
| [CreateRecursive]       |

:star: - commonly used

[All]: operators/All.md
[AsObservable]: operators/AsObservable.md
[AuditTime]: operators/AuditTime.md
[AsyncSubject]: operators/AsyncSubject.md
[AutoConnect]: operators/AutoConnect.md
[Average]: operators/Average.md
[BehaviorSubject]: operators/BehaviorSubject.md
[Catch]: operators/Catch.md
[CatchError]: operators/CatchError.md
[CombineLatest]: operators/CombineLatest.md
[CombineLatestAll]: operators/CombineLatestAll.md
[CombineLatestMap]: operators/CombineLatestMap.md
[CombineLatestMapTo]: operators/CombineLatestMapTo.md
[CombineLatestWith]: operators/CombineLatestWith.md
[Concat]: operators/Concat.md
[ConcatAll]: operators/ConcatAll.md
[ConcatMap]: operators/ConcatMap.md
[ConcatMapTo]: operators/ConcatMapTo.md
[ConcatWith]: operators/ConcatWith.md
[Connect]: operators/Connect.md
[Count]: operators/Count.md
[Create]: operators/Create.md
[CreateFutureRecursive]: operators/CreateFutureRecursive.md
[CreateRecursive]: operators/CreateRecursive.md
[DebounceTime]: operators/DebounceTime.md
[Defer]: operators/Defer.md
[Delay]: operators/Delay.md
[Distinct]: operators/Distinct.md
[Do]: operators/Do.md
[DoOnComplete]: operators/DoOnComplete.md
[DoOnError]: operators/DoOnError.md
[ElementAt]: operators/ElementAt.md
[Empty]: operators/Empty.md
[Filter]: operators/Filter.md
[Finally]: operators/Finally.md
[First]: operators/First.md
[From]: operators/From.md
[FromChan]: operators/FromChan.md
[IgnoreCompletion]: operators/IgnoreCompletion.md
[IgnoreElements]: operators/IgnoreElements.md
[Interval]: operators/Interval.md
[Just]: operators/Just.md
[Last]: operators/Last.md
[Map]: operators/Map.md
[Max]: operators/Max.md
[Merge]: operators/Merge.md
[MergeAll]: operators/MergeAll.md
[MergeDelayError]: operators/MergeDelayError.md
[MergeDelayErrorWith]: operators/MergeDelayErrorWith.md
[MergeMap]: operators/MergeMap.md
[MergeMapTo]: operators/MergeMapTo.md
[MergeWith]: operators/MergeWith.md
[Min]: operators/Min.md
[Never]: operators/Never.md
[ObserveOn]: operators/ObserveOn.md
[Of]: operators/Of.md
[Only]: operators/Only.md
[Passthrough]: operators/Passthrough.md
[Println]: operators/Println.md
[Publish]: operators/Publish.md
[PublishReplay]: operators/PublishReplay.md
[Range]: operators/Range.md
[Reduce]: operators/Reduce.md
[RefCount]: operators/RefCount.md
[Repeat]: operators/Repeat.md
[ReplaySubject]: operators/ReplaySubject.md
[Retry]: operators/Retry.md
[Sample]: operators/Sample.md
[Scan]: operators/Scan.md
[Serialize]: operators/Serialize.md
[Single]: operators/Single.md
[Skip]: operators/Skip.md
[SkipLast]: operators/SkipLast.md
[Start]: operators/Start.md
[Subject]: operators/Subject.md
[Subscribe]: operators/Subscribe.md
[SubscribeOn]: operators/SubscribeOn.md
[Sum]: operators/Sum.md
[SwitchAll]: operators/SwitchAll.md
[SwitchMap]: operators/SwitchMap.md
[Take]: operators/Take.md
[TakeLast]: operators/TakeLast.md
[TakeUntil]: operators/TakeUntil.md
[TakeWhile]: operators/TakeWhile.md
[ThrottleTime]: operators/ThrottleTime.md
[Throw]: operators/Throw.md
[Ticker]: operators/Ticker.md
[TimeInterval]: operators/TimeInterval.md
[Timeout]: operators/Timeout.md
[Timer]: operators/Timer.md
[Timestamp]: operators/Timestamp.md
[ToChan]: operators/ToChan.md
[ToSingle]: operators/ToSingle.md
[ToSlice]: operators/ToSlice.md
[Wait]: operators/Wait.md
[WithLatestFrom]: operators/WithLatestFrom.md
[WithLatestFromAll]: operators/WithLatestFromAll.md
