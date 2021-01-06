# test

    import "github.com/reactivego/rx/test"

[![](../../assets/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test?tab=subdirectories)
[![](../../assets/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test#pkg-subdirectories)
[![](../../assets/rx.svg)](http://reactivex.io/documentation/operators.html)

Package `test` provides tests for the generic rx package.

Every operator / data type has its own subdirectory named after it containing one or more tests that exercise its functionality.
    
| A … C                   | D … L                 | M … P                  | R … S              | T … W          |
|:------------------------|:----------------------|:-----------------------|:-------------------|:---------------|
| [All]                   | [DebounceTime] :star: | [Map] :star:           | [Range]            | [Take]         |
| [AsObservable]          | [Defer]               | [Max]                  | [Reduce]           | [TakeLast]     |
| [AsyncSubject]          | [Delay]               | [Merge] :star:         | [RefCount]         | [TakeUntil]    |
| [AuditTime]             | [Distinct]            | [MergeAll]             | [Repeat]           | [TakeWhile]    |
| [AutoConnect]           | [Do] :star:           | [MergeDelayError]      | [ReplaySubject]    | [ThrottleTime] |
| [Average]               | [DoOnComplete]        | [MergeDelayErrorWith]  | [Retry]            | [Throw]        |
| [BehaviorSubject]       | [DoOnError]           | [MergeMap] :star:      | [Sample]           | [Ticker]       |
| [Catch] :star:          | [ElementAt]           | [MergeMapTo]           | [Scan] :star:      | [TimeInterval] |
| [CatchError] :star:     | [Empty]               | [MergeWith] :star:     | [Serialize]        | [Timeout]      |
| [CombineLatest] :star:  | [Filter] :star:       | [Min]                  | [Single]           | [Timer]        |
| [CombineLatestAll]      | [Finally]             | [Never]                | [Skip]             | [Timestamp]    |
| [CombineLatestMap]      | [First]               | [ObserveOn]            | [SkipLast]         | [ToChan]       |
| [CombineLatestMapTo]    | [From] :star:         | [Of] :star:            | [Start]            | [ToSingle]     |
| [CombineLatestWith]     | [FromChan]            | [Only]                 | [Subject]          | [ToSlice]      |
| [Concat] :star:         | [IgnoreCompletion]    | [Passthrough]          | [Subscribe]        | [Wait]         |
| [ConcatAll]             | [IgnoreElements]      | [Println]              | [SubscribeOn]      |
| [ConcatMap] :star:      | [Interval]            | [Publish] :star:       | [Sum]              |
| [ConcatMapTo]           | [Just] :star:         | [PublishReplay] :star: | [SwitchAll]        |
| [ConcatWith] :star:     | [Last]                |                        | [SwitchMap] :star: |
| [Connect]               |
| [Count]                 |
| [Create] :star:         |
| [CreateFutureRecursive] |
| [CreateRecursive]       |

:star: - commonly used

[All]: All
[AsObservable]: AsObservable
[AuditTime]: AuditTime
[AsyncSubject]: AsyncSubject
[AutoConnect]: AutoConnect
[Average]: Average
[BehaviorSubject]: BehaviorSubject
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
[MergeMapTo]: MergeMapTo
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
