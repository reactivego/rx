# Operators

Currently 98 operators have been implemented:
    
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
| [CombineLatestWith]     | [FromChan]            | [Only]                 | [StartWith] :star: | [ToSlice]               |
| [Concat] :star:         | [IgnoreCompletion]    | [Passthrough]          | [Subject]          | [Wait]                  |
| [ConcatAll]             | [IgnoreElements]      | [Println]              | [Subscribe]        | [WithLatestFrom] :star: |
| [ConcatMap] :star:      | [Interval]            | [Publish] :star:       | [SubscribeOn]      | [WithLatestFromAll]     |
| [ConcatMapTo]           | [Just] :star:         | [PublishReplay] :star: | [Sum]              |
| [ConcatWith] :star:     | [Last]                |                        | [SwitchAll]        |
| [Connect]               |                       |                        | [SwitchMap] :star: 
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
[StartWith]: #startwith
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
[WithLatestFrom]: #withlatestfrom
[WithLatestFromAll]: #withlatestfromall

## All
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.All)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx#Observable.All)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/all.html)

**All** determines whether all items emitted by an Observable meet some
criteria.

Pass a predicate function to the **All** operator that accepts an item emitted
by the source Observable and returns a boolean value based on an
evaluation of that item. **All** returns an ObservableBool that emits a single
boolean value: true if and only if the source Observable terminates
normally and every item emitted by the source Observable evaluated as
true according to the predicate; false if any item emitted by the source
Observable evaluates as false according to the predicate.

![All](../assets/All.svg?raw=true)

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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.AsObservable)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx#Observable.AsObservable)

`AsObservableInt` or `AsObservableBool` type asserts an `Observable` to an observable of type `int` or `bool`.
Also `AsObservable` can be called on an `ObservableInt` and `ObservableBool` to convert to an observable of type `interface{}`.


## AuditTime
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.AuditTime)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx#Observable.AuditTime)
[![](../assets/rx.svg?raw=true)](https://rxjs.dev/api/operators/auditTime)

**AuditTime** waits until the source emits and then starts a timer. When the timer
expires, **AuditTime** will emit the last value received from the source during the
time period when the timer was active.

## AsyncSubject

#### TBD

## AutoConnect

#### TBD

## Average

#### TBD

## BehaviorSubject

#### TBD

## Catch

#### TBD

## CatchError

#### TBD

## CombineLatest

#### TBD

## CombineLatestAll

#### TBD

## CombineLatestMap

#### TBD

## CombineLatestMapTo

#### TBD

## CombineLatestWith

#### TBD

## Concat

#### TBD

## ConcatAll

#### TBD

## ConcatMap

#### TBD

## ConcatMapTo

#### TBD

## ConcatWith

![ConcatWith](../assets/ConcatWith.svg?raw=true)

#### TBD

## Connect

#### TBD

## Count

#### TBD

## Create

#### TBD

## CreateFutureRecursive

#### TBD

## CreateRecursive

#### TBD

## DebounceTime

#### TBD

## Defer

#### TBD

## Delay

#### TBD

## Distinct

#### TBD

## Do

#### TBD

## DoOnComplete

#### TBD

## DoOnError

#### TBD

## ElementAt

#### TBD

## Empty

#### TBD

## Filter

#### TBD

## Finally

#### TBD

## First

#### TBD

## From

#### TBD

## FromChan

#### TBD

## IgnoreCompletion

#### TBD

## IgnoreElements

#### TBD

## Interval

#### TBD

## Just

#### TBD

## Last

#### TBD

## Map

#### TBD

## Max

#### TBD

## Merge

![Merge](../assets/Merge.svg?raw=true)

#### TBD

## MergeAll

#### TBD

## MergeDelayError

![MergeDelayError](../assets/MergeDelayError.svg?raw=true)

#### TBD

## MergeDelayErrorWith

![MergeDelayErrorWith](../assets/MergeDelayErrorWith.svg?raw=true)

#### TBD

## MergeMap

#### TBD

## MergeMapTo

#### TBD

## MergeWith

![MergeWith](../assets/MergeWith.svg?raw=true)

#### TBD

## Min

#### TBD

## Never

#### TBD

## ObserveOn

#### TBD

## Of

#### TBD

## Only

#### TBD

## Passthrough

#### TBD

## Println

#### TBD

## Publish

#### TBD

## PublishReplay

#### TBD

## Range

#### TBD

## Reduce

#### TBD

## RefCount

#### TBD

## Repeat

#### TBD

## ReplaySubject

#### TBD

## Retry

#### TBD

## SampleTime

#### TBD

## Scan

#### TBD

## Serialize

#### TBD

## Single

#### TBD

## Skip

#### TBD

## SkipLast

#### TBD

## Start

#### TBD

## StartWith
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.StartWith)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx#Observable.StartWith)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/startwith.html)

**StartWith** returns an observable that, at the moment of subscription, will
synchronously emit all values provided to this operator, then subscribe to
the source and mirror all of its emissions to subscribers.

![StartWith](../assets/StartWith.svg?raw=true)

Code:
```go
rx.From(2, 3).StartWith(1).Println()
```
Output:
```
1
2
3
```
## Subject

#### TBD

## Subscribe

#### TBD

## SubscribeOn
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc#Observable.SubscribeOn)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx#Observable.SubscribeOn)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/subscribeon.html)

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

#### TBD

## SwitchAll

#### TBD

## SwitchMap

#### TBD

## Take

#### TBD

## TakeLast

#### TBD

## TakeUntil

#### TBD

## TakeWhile

#### TBD

## ThrottleTime

#### TBD

## Throw

#### TBD

## Ticker

#### TBD

## TimeInterval

#### TBD

## Timeout

#### TBD

## Timer

#### TBD

## Timestamp

#### TBD

## ToChan

#### TBD

## ToSingle

#### TBD

## ToSlice

#### TBD

## Wait

#### TBD

## WithLatestFrom

![WithLatestFrom](../assets/WithLatestFrom.svg?raw=true)

#### TBD

## WithLatestFromAll

![WithLatestFromAll](../assets/WithLatestFromAll.svg?raw=true)

#### TBD

