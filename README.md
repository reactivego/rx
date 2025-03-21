# rx

    import "github.com/reactivego/rx"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/rx.svg)](https://pkg.go.dev/github.com/reactivego/rx)

Package `rx` delivers _**R**eactive E**x**tensions_, a powerful API for asynchronous programming in Go, built around [observables](#observables) and [operators](#operators) to handle data streams seamlessly.
## Prerequisites

You’ll need [*Go 1.23*](https://golang.org/dl/) or later, as it includes support for generics and iterators.

## Observables
In `rx`, an [**Observables**](http://reactivex.io/documentation/observable.html) represents a stream of data that can emit items over time, while an **observer** subscribes to it to receive and react to those emissions. This reactive approach enables asynchronous and concurrent operations without blocking execution. Instead of waiting for values to become available, an observer passively listens and responds whenever the Observable emits data, errors, or a completion signal.

This page introduces the **reactive pattern**, explaining what **Observables** and **observers** are and how subscriptions work. Other sections explore the powerful set of **Observable operators** that allow you to transform, combine, and control data streams efficiently.

An Observable:

- is a stream of events.
- assumes zero to many values over time.
- pushes values
- can take any amount of time to complete (or may never)
- is cancellable
- is lazy (it doesn't do anything until you subscribe).

Example
```go
package main

import "github.com/reactivego/x"

func main() {
    x.From[any](1,"hi",2.3).Println()
}
```
> Note the program creates a mixed type `any` observable from an int, string and a float64.

Output
```
1
hi
2.3
```
Example
```go
package main

import "github.com/reactivego/rx"

func main() {
    rx.From(1,2,3).Println()
}
```
> Note the program uses inferred type `int` for the observable.

Output
```
1
2
3
```

Observables in `x` are somewhat similar to Go channels but have much richer
semantics:

Observables can be hot or cold. A hot observable will try to emit values even
when nobody is subscribed. Values emitted during that period will be lost.
The position of a mouse pointer or the current time are examples of hot observables.

A cold observable will only start emitting values after somebody subscribes.
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


### Index
- __All__ determines whether all items emitted by an Observable meet some criteria.
- __AsObservable__ when called on an Observable source will type assert the 'any' items of the source to 'bar' items.
- __AsyncSubject__ emits the last value (and only the last value) emitted by the Observable part, and only after that Observable part completes.
- __AuditTime__ waits until the source emits and then starts a timer.
- __AutoConnect__ makes a Multicaster behave like an ordinary Observable that automatically connects the multicaster to its source when the specified number of observers have subscribed to it.
- __Average__ calculates the average of numbers emitted by an Observable and emits this average.
- __BehaviorSubject__ returns a new BehaviorSubject.
- __Buffer__ buffers the source Observable values until closingNotifier emits.
- __BufferTime__ buffers the source Observable values for a specific time period and emits those as a slice periodically in time.
- __Catch__ recovers from an error notification by continuing the sequence without emitting the error but by switching to the catch ObservableInt to provide items.
- __CatchError__ catches errors on the Observable to be handled by returning a new Observable or throwing an error.
- __CombineLatest__ will subscribe to all Observables.
- __CombineLatestAll__ flattens a higher order observable.
- __CombineLatestMap__ maps every entry emitted by the Observable into an Observable, and then subscribe to it, until the source observable completes.
- __CombineLatestMapTo__ maps every entry emitted by the Observable into a single Observable, and then subscribe to it, until the source observable completes.
- __CombineLatestWith__ will subscribe to its Observable and all other Observables passed in.
- __Concat__ emits the emissions from two or more observables without interleaving them.
- __ConcatAll__ flattens a higher order observable by concattenating the observables it emits.
- __ConcatMap__ transforms the items emitted by an Observable by applying a function to each item and returning an Observable.
- __ConcatMapTo__ maps every entry emitted by the Observable into a single Observable.
- __ConcatWith__ emits the emissions from two or more observables without interleaving them.
- __Connect__ instructs a connectable Observable to begin emitting items to its subscribers.
- __Count__ counts the number of items emitted by the source ObservableInt and emits only this value.
- __Create__ provides a way of creating an Observable from scratch by calling observer methods programmatically.
- __CreateFutureRecursive__ provides a way of creating an Observable from scratch by calling observer methods programmatically.
- __CreateRecursive__ provides a way of creating an Observable from scratch by calling observer methods programmatically.
- __DebounceTime__ only emits the last item of a burst from an Observable if a particular timespan has passed without it emitting another item.
- __Defer__ does not create the Observable until the observer subscribes.
- __Delay__ shifts the emission from an Observable forward in time by a particular amount of time.
- __Distinct__ suppress duplicate items emitted by an Observable.
- __DistinctUntilChanged__ only emits when the current value is different from the last.
- __Do__ calls a function for each next value passing through the observable.
- __DoOnComplete__ calls a function when the stream completes.
- __DoOnError__ calls a function for any error on the stream.
- __ElementAt__ emit only item n emitted by an Observable.
- __Empty__ creates an Observable that emits no items but terminates normally.
- __Filter__ emits only those items from an observable that pass a predicate test.
- __Finally__ applies a function for any error or completion on the stream.
- __First__ emits only the first item, or the first item that meets a condition, from an Observable.
- __From__ creates an observable from multiple values passed in.
- __FromChan__ creates an Observable from a Go channel.
- __IgnoreCompletion__ only emits items and never completes, neither with Error nor with Complete.
- __IgnoreElements__ does not emit any items from an Observable but mirrors its termination notification.
- __Interval__ creates an ObservableInt that emits a sequence of integers spaced by a particular time interval.
- __Just__ creates an observable that emits a particular item.
- __Last__ emits only the last item emitted by an Observable.
- __Map__ transforms the items emitted by an Observable by applying a function to each item.
- __MapTo__ transforms the items emitted by an Observable.
- __Max__ determines, and emits, the maximum-valued item emitted by an Observable.
- __Merge__ combines multiple Observables into one by merging their emissions.
- __MergeAll__ flattens a higher order observable by merging the observables it emits.
- __MergeDelayError__ combines multiple Observables into one by merging their emissions.
- __MergeDelayErrorWith__ combines multiple Observables into one by merging their emissions.
- __MergeMap__ transforms the items emitted by an Observable by applying a function to each item an returning an Observable.
- __MergeMapTo__ maps every entry emitted by the Observable into a single Observable.
- __MergeWith__ combines multiple Observables into one by merging their emissions.
- __Min__ determines, and emits, the minimum-valued item emitted by an Observable.
- __Never__ creates an Observable that emits no items and does't terminate.
- __ObserveOn__ specifies a schedule function to use for delivering values to the observer.
- __ObserverObservable__ actually is an observer that is made observable.
- __Of__ emits a variable amount of values in a sequence and then emits a complete notification.
- __Only__ filters the value stream of an observable and lets only the values of a specific type pass.
- __Passthrough__ just passes through all output from the Observable.
- __Println__ subscribes to the Observable and prints every item to os.Stdout while it waits for completion or error.
- __Publish__ returns a Multicaster for a Subject to an underlying Observable and turns the subject into a connnectable observable.
- __PublishBehavior__ returns a Multicaster that shares a single subscription to the underlying Observable returning an initial value or the last value emitted by the underlying Observable.
- __PublishLast__ returns a Multicaster that shares a single subscription to the underlying Observable containing only the last value emitted before it completes.
- __PublishReplay__ returns a Multicaster for a ReplaySubject to an underlying Observable and turns the subject into a connectable observable.
- __Range__ creates an Observable that emits a range of sequential int values.
- __Reduce__ applies a reducer function to each item emitted by an Observable and the previous reducer result.
- __RefCount__ makes a Connectable behave like an ordinary Observable.
- __Repeat__ creates an observable that emits a sequence of items repeatedly.
- __ReplaySubject__ ensures that all observers see the same sequence of emitted items, even if they subscribe after.
- __Retry__ if a source Observable sends an error notification, resubscribe to it in the hopes that it will complete without error.
- __SampleTime__ emits the most recent item emitted by an Observable within periodic time intervals.
- __Scan__ applies a accumulator function to each item emitted by an Observable and the previous accumulator result.
- __Serialize__ forces an observable to make serialized calls and to be well-behaved.
- __Single__ enforces that the observable sends exactly one data item and then completes.
- __Skip__ suppresses the first n items emitted by an Observable.
- __SkipLast__ suppresses the last n items emitted by an Observable.
- __Start__ creates an Observable that emits the return value of a function.
- __StartWith__ returns an observable that, at the moment of subscription, will synchronously emit all values provided to this operator, then subscribe to the source and mirror all of its emissions to subscribers.
- __Subject__ is a combination of an observer and observable.
- __Subscribe__ operates upon the emissions and notifications from an Observable.
- __SubscribeOn__ specifies the scheduler an Observable should use when it is subscribed to.
- __Sum__ calculates the sum of numbers emitted by an Observable and emits this sum.
- __SwitchAll__ converts an Observable that emits Observables into a single Observable that emits the items emitted by the most-recently-emitted of those Observables.
- __SwitchMap__ transforms the items emitted by an Observable by applying a function to each item an returning an Observable.
- __Take__ emits only the first n items emitted by an Observable.
- __TakeLast__ emits only the last n items emitted by an Observable.
- __TakeUntil__ emits items emitted by an Observable until another Observable emits an item.
- __TakeWhile__ mirrors items emitted by an Observable until a specified condition becomes false.
- __ThrottleTime__ emits when the source emits and then starts a timer during which all emissions from the source are ignored.
- __Throw__ creates an observable that emits no items and terminates with an error.
- __Ticker__ creates an ObservableTime that emits a sequence of timestamps after an initialDelay has passed.
- __TimeInterval__ intercepts the items from the source Observable and emits in their place a struct that indicates the amount of time that elapsed between pairs of emissions.
- __Timeout__ mirrors the source Observable, but issue an error notification if a particular period of time elapses without any emitted items.
- __Timer__ creates an Observable that emits a sequence of integers (starting at zero) after an initialDelay has passed.
- __Timestamp__ attaches a timestamp to each item emitted by an observable indicating when it was emitted.
- __ToChan__ returns a channel that emits 'any' values.
- __ToSingle__ blocks until the Observable emits exactly one value or an error.
- __ToSlice__ collects all values from the Observable into an slice.
- __Wait__ subscribes to the Observable and waits for completion or error.
- __WithLatestFrom__ will subscribe to all Observables and wait for all of them to emit before emitting the first slice.
- __WithLatestFromAll__ flattens a higher order observable.
