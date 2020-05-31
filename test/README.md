# test

    import "github.com/reactivego/rx/test"

[![](../svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test?tab=subdirectories)
[![](../svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test#pkg-subdirectories)
[![](../svg/rx.svg)](http://reactivex.io/documentation/operators.html)

Package `test` provides tests for the generic rx package.

Every operator / data type has its own subdirectory named after it containing one or more tests that exercise its functionality.

Currently 92 operators have been implemented:
    
| A … C                   | D … L                 | M … P                  | R … S              | T … W          |
|:------------------------|:----------------------|:-----------------------|:-------------------|:---------------|
| [All]                   | [DebounceTime] :star: | [Map] :star:           | [Range]            | [Take]         |
| [AsObservable]          | [Defer]               | [Max]                  | [Reduce]           | [TakeLast]     |
| [AuditTime]             | [Delay]               | [Merge] :star:         | [RefCount]         | [TakeUntil]    |
| [AutoConnect]           | [Distinct]            | [MergeAll]             | [Repeat]           | [TakeWhile]    |
| [Average]               | [Do] :star:           | [MergeDelayError]      | [ReplaySubject]    | [ThrottleTime] |
| [Catch] :star:          | [DoOnComplete]        | [MergeDelayErrorWith]  | [Retry]            | [Throw]        |
| [CatchError] :star:     | [DoOnError]           | [MergeMap] :star:      | [Sample]           | [Ticker]       |
| [CombineLatest] :star:  | [ElementAt]           | [MergeWith] :star:     | [Scan] :star:      | [TimeInterval] |
| [CombineLatestAll]      | [Empty]               | [Min]                  | [Serialize]        | [Timeout]      |
| [CombineLatestMap]      | [Filter] :star:       | [Never]                | [Single]           | [Timer]        |
| [CombineLatestMapTo]    | [Finally]             | [ObserveOn]            | [Skip]             | [Timestamp]    |
| [CombineLatestWith]     | [First]               | [Of] :star:            | [SkipLast]         | [ToChan]       |
| [Concat] :star:         | [From] :star:         | [Only]                 | [Start]            | [ToSingle]     |
| [ConcatAll]             | [FromChan]            | [Passthrough]          | [Subject]          | [ToSlice]      |
| [ConcatMap] :star:      | [IgnoreCompletion]    | [Println]              | [Subscribe]        | [Wait]         |
| [ConcatMapTo]           | [IgnoreElements]      | [Publish] :star:       | [SubscribeOn]      |
| [ConcatWith] :star:     | [Interval]            | [PublishReplay] :star: | [Sum]              |
| [Connect]               | [Just] :star:         |                        | [SwitchAll]        |
| [Count]                 | [Last]                |                        | [SwitchMap] :star: |
| [Create] :star:         |                       |                        |
| [CreateFutureRecursive] |                       |                        |
| [CreateRecursive]       |                       |                        |
|                         |                       |                        |

:star: - commonly used

[All]: All
[All]: All
[AsObservable]: AsObservable
[AuditTime]: AuditTime
[AutoConnect]: AutoConnect
[Average]: Average
[Catch]: Catch
[CatchError]: CatchError
[CombineLatest]: CombineLatest
[CombineLatestAll]: CombineLatestAll
[CombineLatestMap]: CombineLatestMap
[CombineLatestMapTo]: CombineLatestMapTo
[CombineLatestWith]: CombineLatestWith
[Concat]: Concat
[ConcatAll]: ConcatAll
[ConcatMap]: ConcatMap
[ConcatMapTo]: ConcatMapTo
[ConcatWith]: ConcatWith
[Connect]: Connect
[Count]: Count
[Create]: Create
[CreateFutureRecursive]: CreateFutureRecursive
[CreateRecursive]: CreateRecursive
[DebounceTime]: DebounceTime
[Defer]: Defer
[Delay]: Delay
[Distinct]: Distinct
[Do]: Do
[DoOnComplete]: DoOnComplete
[DoOnError]: DoOnError
[ElementAt]: ElementAt
[Empty]: Empty
[Filter]: Filter
[Finally]: Finally
[First]: First
[From]: From
[FromChan]: FromChan
[IgnoreCompletion]: IgnoreCompletion
[IgnoreElements]: IgnoreElements
[Interval]: Interval
[Just]: Just
[Last]: Last
[Map]: Map
[Max]: Max
[Merge]: Merge
[MergeAll]: MergeAll
[MergeDelayError]: MergeDelayError
[MergeDelayErrorWith]: MergeDelayErrorWith
[MergeMap]: MergeMap
[MergeWith]: MergeWith
[Min]: Min
[Never]: Never
[ObserveOn]: ObserveOn
[Of]: Of
[Only]: Only
[Passthrough]: Passthrough
[Println]: Println
[Publish]: Publish
[PublishReplay]: PublishReplay
[Range]: Range
[Reduce]: Reduce
[RefCount]: RefCount
[Repeat]: Repeat
[ReplaySubject]: ReplaySubject
[Retry]: Retry
[Sample]: Sample
[Scan]: Scan
[Serialize]: Serialize
[Single]: Single
[Skip]: Skip
[SkipLast]: SkipLast
[Start]: Start
[Subject]: Subject
[Subscribe]: Subscribe
[SubscribeOn]: SubscribeOn
[Sum]: Sum
[SwitchAll]: SwitchAll
[SwitchMap]: SwitchMap
[Take]: Take
[TakeLast]: TakeLast
[TakeUntil]: TakeUntil
[TakeWhile]: TakeWhile
[ThrottleTime]: ThrottleTime
[Throw]: Throw
[Ticker]: Ticker
[TimeInterval]: TimeInterval
[Timeout]: Timeout
[Timer]: Timer
[Timestamp]: Timestamp
[ToChan]: ToChan
[ToSingle]: ToSingle
[ToSlice]: ToSlice
[Wait]: Wait
