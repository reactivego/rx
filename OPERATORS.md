# Operators

Currently 104 operators have been implemented:
    
| A … C                   | D … L                         | M … P                  | R … S              | T … W                   |
|:------------------------|:------------------------------|:-----------------------|:-------------------|:------------------------|
| [All]                   | [DebounceTime] :star:         | [Map] :star:           | [Range]            | [Take] :star:           |
| [AsObservable]          | [Defer]                       | [MapTo]                | [Reduce]           | [TakeLast]              |
| [AsyncSubject]          | [Delay]                       | [Max]                  | [RefCount]         | [TakeUntil] :star:      |
| [AuditTime]             | [Distinct]                    | [Merge] :star:         | [Repeat]           | [TakeWhile]             |
| [AutoConnect]           | [DistinctUntilChanged] :star: | [MergeAll]             | [ReplaySubject]    | [ThrottleTime]          |
| [Average]               | [Do] :star:                   | [MergeDelayError]      | [Retry]            | [Throw]                 |
| [BehaviorSubject]       | [DoOnComplete]                | [MergeDelayErrorWith]  | [SampleTime]       | [Ticker]                |
| [Buffer]                | [DoOnError]                   | [MergeMap] :star:      | [Scan] :star:      | [TimeInterval]          |
| [BufferTime] :star:     | [ElementAt]                   | [MergeMapTo]           | [Serialize]        | [Timeout]               |
| [Catch] :star:          | [Empty]                       | [MergeWith] :star:     | [Single]           | [Timer]                 |
| [CatchError] :star:     | [Filter] :star:               | [Min]                  | [Skip]             | [Timestamp]             |
| [CombineLatest] :star:  | [Finally]                     | [Never]                | [SkipLast]         | [ToChan]                |
| [CombineLatestAll]      | [First]                       | [ObserveOn]            | [Start]            | [ToSingle]              |
| [CombineLatestMap]      | [From] :star:                 | [Of] :star:            | [StartWith] :star: | [ToSlice]               |
| [CombineLatestMapTo]    | [FromChan]                    | [Only]                 | [Subject]          | [Wait]                  |
| [CombineLatestWith]     | [IgnoreCompletion]            | [Passthrough]          | [Subscribe]        | [WithLatestFrom] :star: |
| [Concat] :star:         | [IgnoreElements]              | [Println]              | [SubscribeOn]      | [WithLatestFromAll]     |
| [ConcatAll]             | [Interval]                    | [Publish] :star:       | [Sum]              |
| [ConcatMap] :star:      | [Just] :star:                 | [PublishBehavior]      | [SwitchAll]        |
| [ConcatMapTo]           | [Last]                        | [PublishLast]          | [SwitchMap] :star: |
| [ConcatWith] :star:     |                               | [PublishReplay] :star: |
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
[Buffer]: #buffer
[BufferTime]: #buffertime
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
[DistinctUntilChanged]: #distinctuntilchanged
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
[MapTo]: #mapto
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
[PublishBehavior]: #publishbehavior
[PublishLast]: #publishlast
[PublishReplay]: #publishreplay
[Range]: #range
[Reduce]: #reduce
[RefCount]: #refcount
[Repeat]: #repeat
[ReplaySubject]: #replaysubject
[Retry]: #retry
[SampleTime]: #sampletime
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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.All)
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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.AsObservable)

`AsObservableInt` or `AsObservableBool` type asserts an `Observable` to an observable of type `int` or `bool`.
Also `AsObservable` can be called on an `ObservableInt` and `ObservableBool` to convert to an observable of type `interface{}`.


## AuditTime
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.AuditTime)
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

## Buffer
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.Buffer)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/buffer.html)

**Buffer** buffers the source Observable values until closingNotifier emits.

![Buffer](../assets/Buffer.svg?raw=true)

## Example
```go
import _ "github.com/reactivego/rx/generic"
```
Code:
```go
source := rx.Timer(0*ms, 100*ms).Take(4).ConcatMap(func(i interface{}) rx.Observable {
    switch i.(int) {
    case 0:
        return rx.From("a", "b")
    case 1:
        return rx.From("c", "d", "e")
    case 3:
        return rx.From("f", "g")
    }
    return rx.Empty()
})
closingNotifier := rx.Interval(100 * ms)
source.Buffer(closingNotifier).Println()
```
Output:
```
[a b]
[c d e]
[]
[f g]
```
## BufferTime
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.BufferTime)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/buffer.html)

**BufferTime** buffers the source Observable values for a specific time period and emits those as a
slice periodically in time.

![BufferTime](../assets/BufferTime.svg?raw=true)

Code:
```go
const ms = time.Millisecond
source := rx.Timer(0*ms, 100*ms).Take(4).ConcatMap(func(i interface{}) rx.Observable {
    switch i.(int) {
    case 0:
        return rx.From("a", "b")
    case 1:
        return rx.From("c", "d", "e")
    case 3:
        return rx.From("f", "g")
    }
    return rx.Empty()
})
source.BufferTime(100 * ms).Println()
```
Output:
```
[a b]
[c d e]
[]
[f g]
```
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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.Distinct)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/distinct.html)

**Distinct** suppress duplicate items emitted by an Observable.

The operator filters an Observable by only allowing items through that have not already been emitted.
In some implementations there are variants that allow you to adjust the criteria by which two items are
considered “distinct.” In some, there is a variant of the operator that only compares an item against its
immediate predecessor for distinctness, thereby filtering only consecutive duplicate items from the sequence.

![Distinct](../assets/Distinct.svg?raw=true)

Code:
```go
rx.From(1, 2, 2, 1, 3).Distinct().Println()
```
Output:
```
1
2
3
```
## DistinctUntilChanged
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.DistinctUntilChanged)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/documentation/operators/distinct.html)

**DistinctUntilChanged** only emits when the current value is different from the last.

The operator only compares emitted items from the source Observable against their immediate
predecessors in order to determine whether or not they are distinct.

![DistinctUntilChanged](../assets/DistinctUntilChanged.svg?raw=true)

Code:
```go
rx.From(1, 2, 2, 1, 3).DistinctUntilChanged().Println()
```
Output:
```
1
2
1
3
```
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

## MapTo

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

## PublishBehavior

#### TBD

## PublishLast

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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.StartWith)
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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.SubscribeOn)
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
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#Observable.WithLatestFrom)
[![](../assets/rx.svg?raw=true)](https://rxjs-dev.firebaseapp.com/api/operators/withLatestFrom)

**WithLatestFrom** will subscribe to all Observables and wait for all of them to emit before emitting
the first slice. The source observable determines the rate at which the values are emitted. The idea
is that observables that are faster than the source, don't determine the rate at which the resulting
observable emits. The observables that are combined with the source will be allowed to continue
emitting but only will have their last emitted value emitted whenever the source emits.

Note that any values emitted by the source before all other observables have emitted will
effectively be lost. The first emit will occur the first time the source emits after all other
observables have emitted.

![WithLatestFrom](../assets/WithLatestFrom.svg?raw=true)

Code:
```go
a := rx.From(1,2,3,4,5)
b := rx.From("A","B","C","D","E")
a.WithLatestFrom(b).Println()
```
Output:
```
[2 A]
[3 B]
[4 C]
[5 D]
```
## WithLatestFromAll
[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx#ObservableObservable.WithLatestFromAll)
[![](../assets/rx.svg?raw=true)](https://rxjs-dev.firebaseapp.com/api/operators/withLatestFrom)

**WithLatestFromAll** flattens a higher order observable (e.g. ObservableObservable) by subscribing
to all emitted observables (ie. Observable entries) until the source completes. It will then wait
for all of the subscribed Observables to emit before emitting the first slice. The first observable
that was emitted by the source will be used as the trigger observable. Whenever the trigger
observable emits, a new slice will be emitted containing all the latest values.

Note that any values emitted by the source before all other observables have emitted will
effectively be lost. The first emit will occur the first time the source emits after all other
observables have emitted.

![WithLatestFromAll](../assets/WithLatestFromAll.svg?raw=true)

Code:
```go
a := rx.From(1, 2, 3, 4, 5)
b := rx.From("A", "B", "C", "D", "E")
c := rx.FromObservable(a, b)
c.WithLatestFromAll().Println()
```
Output:
```
[2 A]
[3 B]
[4 C]
[5 D]
```
