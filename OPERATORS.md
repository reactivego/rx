# Operators

Currently 94 operators have been implemented:
    
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

[All]: #all
[AsObservable]: #asobservable
[AuditTime]: #audittime
[AsyncSubject]: #asyncsubject
[AutoConnect]: #autoconnect
[Average]: #average
[BehaviorSubject]: #behaviorsubject
[Catch]: #catch
[CatchError]: #catcherror
[CombineLatest]: #combinelatest
[CombineLatestAll]: #combinelatestall
[CombineLatestMap]: #combinelatestmap
[CombineLatestMapTo]: #combinelatestmapto
[CombineLatestWith]: #combinelatestwith
[Concat]: #concat
[ConcatAll]: #concatall
[ConcatMap]: #concatmap
[ConcatMapTo]: #concatmapto
[ConcatWith]: #concatwith
[Connect]: #connect
[Count]: #count
[Create]: #create
[CreateFutureRecursive]: #createfuturerecursive
[CreateRecursive]: #createrecursive
[DebounceTime]: #debouncetime
[Defer]: #defer
[Delay]: #delay
[Distinct]: #distinct
[Do]: #do
[DoOnComplete]: #dooncomplete
[DoOnError]: #doonerror
[ElementAt]: #elementat
[Empty]: #empty
[Filter]: #filter
[Finally]: #finally
[First]: #first
[From]: #from
[FromChan]: #fromchan
[IgnoreCompletion]: #ignorecompletion
[IgnoreElements]: #ignoreelements
[Interval]: #interval
[Just]: #just
[Last]: #last
[Map]: #map
[Max]: #max
[Merge]: #merge
[MergeAll]: #mergeall
[MergeDelayError]: #mergedelayerror
[MergeDelayErrorWith]: #mergedelayerrorwith
[MergeMap]: #mergemap
[MergeMapTo]: #mergemap
[MergeWith]: #mergewith
[Min]: #min
[Never]: #never
[ObserveOn]: #observeon
[Of]: #of
[Only]: #only
[Passthrough]: #passthrough
[Println]: #println
[Publish]: #publish
[PublishReplay]: #publishreplay
[Range]: #range
[Reduce]: #reduce
[RefCount]: #refcount
[Repeat]: #repeat
[ReplaySubject]: #replaysubject
[Retry]: #retry
[Sample]: #sample
[Scan]: #scan
[Serialize]: #serialize
[Single]: #single
[Skip]: #skip
[SkipLast]: #skiplast
[Start]: #start
[Subject]: #subject
[Subscribe]: #subscribe
[SubscribeOn]: #subscribeon
[Sum]: #sum
[SwitchAll]: #switchall
[SwitchMap]: #switchmap
[Take]: #take
[TakeLast]: #takelast
[TakeUntil]: #takeuntil
[TakeWhile]: #takewhile
[ThrottleTime]: #throttletime
[Throw]: #throw
[Ticker]: #ticker
[TimeInterval]: #timeinterval
[Timeout]: #timeout
[Timer]: #timer
[Timestamp]: #timestamp
[ToChan]: #tochan
[ToSingle]: #tosingle
[ToSlice]: #toslice
[Wait]: #wait

## All
[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.All)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx#Observable.All)
[![](svg/rx.svg)](http://reactivex.io/documentation/operators/all.html)

**All** determines whether all items emitted by an Observable meet some
criteria.

Pass a predicate function to the **All** operator that accepts an item emitted
by the source Observable and returns a boolean value based on an
evaluation of that item. **All** returns an ObservableBool that emits a single
boolean value: true if and only if the source Observable terminates
normally and every item emitted by the source Observable evaluated as
true according to the predicate; false if any item emitted by the source
Observable evaluates as false according to the predicate.

![All](svg/All.svg)

Code:
```go
// Setup All to produce true only when all source values are less than 5
lessthan5 := func(i interface{}) bool {
	return i.(int) < 5
}

result, err := rx.From(1, 2, 5, 2, 1).All(lessthan5).ToSingle()

fmt.Println("All values less than 5?", result, err)

result, err = rx.From(4, 1, 0, -1, 2, 3, 4).All(lessthan5).ToSingle()

fmt.Println("All values less than 5?", result, err)
```
Output:
```
All values less than 5? false <nil>
All values less than 5? true <nil>
```
## AsObservable
[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.AsObservable)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx#Observable.AsObservable)

`AsObservableInt` or `AsObservableBool` type asserts an `Observable` to an observable of type `int` or `bool`.
Also `AsObservable` can be called on an `ObservableInt` and `ObservableBool` to convert to an observable of type `interface{}`.


## AuditTime
[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.AuditTime)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx#Observable.AuditTime)
[![](svg/rx.svg)](https://rxjs.dev/api/operators/auditTime)

**AuditTime** waits until the source emits and then starts a timer. When the timer
expires, **AuditTime** will emit the last value received from the source during the
time period when the timer was active.

## AsyncSubject
See [AsyncSubject](test/AsyncSubject) in test folder.
## AutoConnect
See [AutoConnect](test/AutoConnect) in test folder.
## Average
See [Average](test/Average) in test folder.
## BehaviorSubject
See [BehaviorSubject](test/BehaviorSubject) in test folder.
## Catch
See [Catch](test/Catch) in test folder.
## CatchError
See [CatchError](test/CatchError) in test folder.
## CombineLatest
See [CombineLatest](test/CombineLatest) in test folder.
## CombineLatestAll
See [CombineLatestAll](test/CombineLatestAll) in test folder.
## CombineLatestMap
See [CombineLatestMap](test/CombineLatestMap) in test folder.
## CombineLatestMapTo
See [CombineLatestMapTo](test/CombineLatestMapTo) in test folder.
## CombineLatestWith
See [CombineLatestWith](test/CombineLatestWith) in test folder.
## Concat
See [Concat](test/Concat) in test folder.
## ConcatAll
See [ConcatAll](test/ConcatAll) in test folder.
## ConcatMap
See [ConcatMap](test/ConcatMap) in test folder.
## ConcatMapTo
See [ConcatMapTo](test/ConcatMapTo) in test folder.
## ConcatWith
See [ConcatWith](test/ConcatWith) in test folder.
## Connect
See [Connect](test/Connect) in test folder.
## Count
See [Count](test/Count) in test folder.
## Create
See [Create](test/Create) in test folder.
## CreateFutureRecursive
See [CreateFutureRecursive](test/CreateFutureRecursive) in test folder.
## CreateRecursive
See [CreateRecursive](test/CreateRecursive) in test folder.
## DebounceTime
See [DebounceTime](test/DebounceTime) in test folder.
## Defer
See [Defer](test/Defer) in test folder.
## Delay
See [Delay](test/Delay) in test folder.
## Distinct
See [Distinct](test/Distinct) in test folder.
## Do
See [Do](test/Do) in test folder.
## DoOnComplete
See [DoOnComplete](test/DoOnComplete) in test folder.
## DoOnError
See [DoOnError](test/DoOnError) in test folder.
## ElementAt
See [ElementAt](test/ElementAt) in test folder.
## Empty
See [Empty](test/Empty) in test folder.
## Filter
See [Filter](test/Filter) in test folder.
## Finally
See [Finally](test/Finally) in test folder.
## First
See [First](test/First) in test folder.
## From
See [From](test/From) in test folder.
## FromChan
See [FromChan](test/FromChan) in test folder.
## IgnoreCompletion
See [IgnoreCompletion](test/IgnoreCompletion) in test folder.
## IgnoreElements
See [IgnoreElements](test/IgnoreElements) in test folder.
## Interval
See [Interval](test/Interval) in test folder.
## Just
See [Just](test/Just) in test folder.
## Last
See [Last](test/Last) in test folder.
## Map
See [Map](test/Map) in test folder.
## Max
See [Max](test/Max) in test folder.
## Merge
See [Merge](test/Merge) in test folder.
## MergeAll
See [MergeAll](test/MergeAll) in test folder.
## MergeDelayError
See [MergeDelayError](test/MergeDelayError) in test folder.
## MergeDelayErrorWith
See [MergeDelayErrorWith](test/MergeDelayErrorWith) in test folder.
## MergeMap
See [MergeMap](test/MergeMap) in test folder.
## MergeMapTo
See [MergeMapTo](test/MergeMapTo) in test folder.
## MergeWith
See [MergeWith](test/MergeWith) in test folder.
## Min
See [Min](test/Min) in test folder.
## Never
See [Never](test/Never) in test folder.
## ObserveOn
See [ObserveOn](test/ObserveOn) in test folder.
## Of
See [Of](test/Of) in test folder.
## Only
See [Only](test/Only) in test folder.
## Passthrough
See [Passthrough](test/Passthrough) in test folder.
## Println
See [Println](test/Println) in test folder.
## Publish
See [Publish](test/Publish) in test folder.
## PublishReplay
See [PublishReplay](test/PublishReplay) in test folder.
## Range
See [Range](test/Range) in test folder.
## Reduce
See [Reduce](test/Reduce) in test folder.
## RefCount
See [RefCount](test/RefCount) in test folder.
## Repeat
See [Repeat](test/Repeat) in test folder.
## ReplaySubject
See [ReplaySubject](test/ReplaySubject) in test folder.
## Retry
See [Retry](test/Retry) in test folder.
## SampleTime
See [SampleTime](test/SampleTime) in test folder.
## Scan
See [Scan](test/Scan) in test folder.
## Serialize
See [Serialize](test/Serialize) in test folder.
## Single
See [Single](test/Single) in test folder.
## Skip
See [Skip](test/Skip) in test folder.
## SkipLast
See [SkipLast](test/SkipLast) in test folder.
## Start
See [Start](test/Start) in test folder.
## Subject
See [Subject](test/Subject) in test folder.
## Subscribe
See [Subscribe](test/Subscribe) in test folder.
## SubscribeOn
[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.SubscribeOn)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx#Observable.SubscribeOn)
[![](svg/rx.svg)](http://reactivex.io/documentation/operators/subscribeon.html)

**SubscribeOn** specifies which [scheduler](https://github.com/reactivego/scheduler) an Observable should use when it is subscribed to.

### Example (SubscribeOn Trampoline)
Code:
```go
trampoline := rx.MakeTrampolineScheduler()
observer := func(next interface{}, err error, done bool) {
	switch {
	case !done:
		fmt.Println(trampoline, "print", next)
	case err != nil:
		fmt.Println(trampoline, "print", err)
	default:
		fmt.Println(trampoline, "print", "complete")
	}
}
fmt.Println(trampoline, "SUBSCRIBING...")
subscription := rx.From(1, 2, 3).SubscribeOn(trampoline).Subscribe(observer)
fmt.Println(trampoline, "WAITING...")
subscription.Wait()
fmt.Println(trampoline, "DONE")
```
Output:
```
Trampoline{ tasks = 0 } SUBSCRIBING...
Trampoline{ tasks = 1 } WAITING...
Trampoline{ tasks = 1 } print 1
Trampoline{ tasks = 1 } print 2
Trampoline{ tasks = 1 } print 3
Trampoline{ tasks = 1 } print complete
Trampoline{ tasks = 0 } DONE
```

### Example (SubscribeOn Goroutine)
Code:
```go
const ms = time.Millisecond
goroutine := rx.GoroutineScheduler()
observer := func(next interface{}, err error, done bool) {
	switch {
	case !done:
		fmt.Println(goroutine, "print", next)
	case err != nil:
		fmt.Println(goroutine, "print", err)
	default:
		fmt.Println(goroutine, "print", "complete")
	}
}
fmt.Println(goroutine, "SUBSCRIBING...")
subscription := rx.From(1, 2, 3).Delay(10 * ms).SubscribeOn(goroutine).Subscribe(observer)
// Note that without a Delay the next Println lands at a random spot in the output.
fmt.Println("WAITING...")
subscription.Wait()
fmt.Println(goroutine, "DONE")
```
Output:
```
Output:
Goroutine{ tasks = 0 } SUBSCRIBING...
WAITING...
Goroutine{ tasks = 1 } print 1
Goroutine{ tasks = 1 } print 2
Goroutine{ tasks = 1 } print 3
Goroutine{ tasks = 1 } print complete
Goroutine{ tasks = 0 } DONE
```

## Sum
See [Sum](test/Sum) in test folder.
## SwitchAll
See [SwitchAll](test/SwitchAll) in test folder.
## SwitchMap
See [SwitchMap](test/SwitchMap) in test folder.
## Take
See [Take](test/Take) in test folder.
## TakeLast
See [TakeLast](test/TakeLast) in test folder.
## TakeUntil
See [TakeUntil](test/TakeUntil) in test folder.
## TakeWhile
See [TakeWhile](test/TakeWhile) in test folder.
## ThrottleTime
See [ThrottleTime](test/ThrottleTime) in test folder.
## Throw
See [Throw](test/Throw) in test folder.
## Ticker
See [Ticker](test/Ticker) in test folder.
## TimeInterval
See [TimeInterval](test/TimeInterval) in test folder.
## Timeout
See [Timeout](test/Timeout) in test folder.
## Timer
See [Timer](test/Timer) in test folder.
## Timestamp
See [Timestamp](test/Timestamp) in test folder.
## ToChan
See [ToChan](test/ToChan) in test folder.
## ToSingle
See [ToSingle](test/ToSingle) in test folder.
## ToSlice
See [ToSlice](test/ToSlice) in test folder.
## Wait
See [Wait](test/Wait) in test folder.
