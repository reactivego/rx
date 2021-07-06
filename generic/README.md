# rx (generic)

    import _ "github.com/reactivego/rx/generic"

[![](../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/generic#section-documentation)
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

Currently 102 operators are implemented:
    
| A … C                   | D … L                         | M … P                  | R … S              | T … W                   |
|:------------------------|:------------------------------|:-----------------------|:-------------------|:------------------------|
| [All]                   | [DebounceTime] :star:         | [Map] :star:           | [Range]            | [Take] :star:           |
| [AsObservable]          | [Defer]                       | [MapTo]                | [Reduce]           | [TakeLast]              |
| [AsyncSubject]          | [Delay]                       | [Max]                  | [RefCount]         | [TakeUntil] :star:      |
| [AuditTime]             | [Distinct]                    | [Merge] :star:         | [Repeat]           | [TakeWhile]             |
| [AutoConnect]           | [DistinctUntilChanged] :star: | [MergeAll]             | [ReplaySubject]    | [ThrottleTime]          |
| [Average]               | [Do] :star:                   | [MergeDelayError]      | [Retry]            | [Throw]                 |
| [BehaviorSubject]       | [DoOnComplete]                | [MergeDelayErrorWith]  | [SampleTime]       | [Ticker]                |
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
| [ConcatWith] :star:     | [Just] :star:                 | [PublishBehavior]      | [SwitchAll]        |
| [Connect]               | [Last]                        | [PublishLast]          | [SwitchMap] :star: |
| [Count]                 |                               | [PublishReplay] :star: |
| [Create] :star:         |
| [CreateFutureRecursive] |
| [CreateRecursive]       |

:star: - commonly used

[All]: ../test/All/README.md
[AsyncSubject]: ../test/AsyncSubject/README.md
[AuditTime]: ../test/AuditTime/README.md
[BehaviorSubject]: ../test/BehaviorSubject/README.md
[Catch]: ../test/Catch/README.md
[CatchError]: ../test/CatchError/README.md
[ConcatWith]: ../test/ConcatWith/README.md
[DebounceTime]: ../test/DebounceTime/README.md
[Distinct]: ../test/Distinct/README.md
[DistinctUntilChanged]: ../test/DistinctUntilChanged/README.md
[Map]: ../test/Map/README.md
[MapTo]: ../test/MapTo/README.md
[Merge]: ../test/Merge/README.md
[MergeDelayError]: ../test/MergeDelayError/README.md
[MergeDelayErrorWith]: ../test/MergeDelayErrorWith/README.md
[MergeMap]: ../test/MergeMap/README.md
[MergeMapTo]: ../test/MergeMapTo/README.md
[MergeWith]: ../test/MergeWith/README.md
[Publish]: ../test/Publish/README.md
[PublishBehavior]: ../test/PublishBehavior/README.md
[PublishLast]: ../test/PublishLast/README.md
[PublishReplay]: ../test/PublishReplay/README.md
[ReplaySubject]: ../test/ReplaySubject/README.md
[Retry]: ../test/Retry/README.md
[SampleTime]: ../test/SampleTime/README.md
[StartWith]: ../test/StartWith/README.md
[Subject]: ../test/Subject/README.md
[ThrottleTime]: ../test/ThrottleTime/README.md
[Timer]: ../test/Timer/README.md
[WithLatestFrom]: ../test/WithLatestFrom/README.md
[WithLatestFromAll]: ../test/WithLatestFromAll/README.md
