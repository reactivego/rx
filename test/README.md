# test

    import "github.com/reactivego/rx/test"

[![](../svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test?tab=subdirectories)
[![](../svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test#pkg-subdirectories)
[![](../svg/rx.svg)](http://reactivex.io/documentation/operators.html)

Package `test` provides tests for the generic rx package.

Every operator / data type has its own subdirectory named after it containing one or more tests that exercise its functionality.

Currently 91 operators have been implemented:
   
| A … C                   | D … M                   | M … S                 | S … T           | T … W          |
|:------------------------|:------------------------|:----------------------|:----------------|:---------------|
| [All]                   | [DebounceTime]          | [Merge]               | [Sample]        | [Timeout]      |
| [AsObservable]          | [Defer]                 | [MergeAll]            | [Scan]          | [Timer]        |
| [AuditTime]             | [Delay]                 | [MergeDelayError]     | [Serialize]     | [Timestamp]    |
| [AutoConnect]           | [Distinct]              | [MergeDelayErrorWith] | [Single]        | [ToChan]       |
| [Average]               | [Do]                    | [MergeMap]            | [Skip]          | [ToSingle]     |
| [Catch]                 | [DoOnComplete]          | [MergeWith]           | [SkipLast]      | [ToSlice]      |
| [CombineLatest]         | [DoOnError]             | [Min]                 | [Start]         | [Wait]         |
| [CombineLatestAll]      | [ElementAt]             | [Never]               | [Subject]       |
| [CombineLatestMap]      | [Empty]                 | [ObserveOn]           | [Subscribe]     |
| [CombineLatestMapTo]    | [Filter]                | [Of]                  | [SubscribeOn]   |
| [CombineLatestWith]     | [Finally]               | [Only]                | [Sum]           |
| [Concat]                | [First]                 | [Passthrough]         | [SwitchAll]     |
| [ConcatAll]             | [From]                  | [Println]             | [SwitchMap]     |
| [ConcatMap]             | [FromChan]              | [Publish]             | [Take]          |
| [ConcatMapTo]           | [IgnoreCompletion]      | [PublishReplay]       | [TakeLast]      |
| [ConcatWith]            | [IgnoreElements]        | [Range]               | [TakeUntil]     |
| [Connect]               | [Interval]              | [Reduce]              | [TakeWhile]     |
| [Count]                 | [Just]                  | [RefCount]            | [ThrottleTime]  |
| [Create]                | [Last]                  | [Repeat]              | [Throw]         |
| [CreateFutureRecursive] | [Map]                   | [ReplaySubject]       | [Ticker]        |
| [CreateRecursive]       | [Max]                   | [Retry]               | [TimeInterval]  |

[All]: All
[All]: All
[AsObservable]: AsObservable
[AuditTime]: AuditTime
[AutoConnect]: AutoConnect
[Average]: Average
[Catch]: Catch
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
