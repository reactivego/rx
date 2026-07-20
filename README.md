# rx

    import "github.com/reactivego/rx"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/rx.svg)](https://pkg.go.dev/github.com/reactivego/rx)

Package `rx` provides _**R**eactive E**x**tensions_ for Go: an API for asynchronous programming built around [Observables](#observable-and-observer) — typed streams of data — and [Operators](#operators) that transform, filter, and combine them into pipelines.

## Prerequisites

You'll need [*Go 1.23*](https://golang.org/dl/) or later; the implementation depends on generics and on iterator (range-over-func) support.

## Quick Start

Create an Observable from some values, print each one, and wait for completion:

```go
package main

import "github.com/reactivego/rx"

func main() {
    rx.From(1, 2, 3).Println().Wait()
}
```
```
1
2
3
```

The element type is inferred; use an explicit type parameter for mixed-type streams:

```go
rx.From[any](1, "hi", 2.3).Println().Wait()
```
```
1
hi
2.3
```

## Core Concepts

### Observable and Observer

An [`Observable[T]`](https://pkg.go.dev/github.com/reactivego/rx#Observable) is a stream of values that are pushed to an [`Observer[T]`](https://pkg.go.dev/github.com/reactivego/rx#Observer) over time. Both are plain function types, so any function with the right signature works:

```go
type Observable[T any] func(Observer[T], Scheduler, Subscriber)

type Observer[T any] func(next T, err error, done bool)
```

An Observable:

- is a stream of events.
- assumes zero to many values over time.
- pushes values to its observers.
- can take any amount of time to complete (or may never).
- is cancellable.
- is lazy: it does nothing until you subscribe.

The observer receives each value with `done == false`. The stream terminates with exactly one final call where `done == true`: with `err == nil` for normal completion, or a non-nil `err` for failure.

```go
sub := rx.From("a", "b").Subscribe(rx.GoroutineContext(), func(next string, err error, done bool) {
    switch {
    case !done:
        fmt.Println(next)
    case err != nil:
        fmt.Println("error:", err)
    default:
        fmt.Println("complete")
    }
})
sub.Wait()
```
```
a
b
complete
```

### Subscriptions

[`Subscribe`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Subscribe) returns a [`Subscription`](https://pkg.go.dev/github.com/reactivego/rx#Subscription) for observing and controlling the stream's lifecycle:

- `Subscribed()` reports whether the subscription is still active.
- `Unsubscribe()` cancels it.
- `Done()` returns a channel that closes on termination, for use in `select`.
- `Wait()` blocks until termination and returns the outcome.
- `Err()` returns `nil` after normal completion, the stream's error after failure, [`ErrSubscriptionCanceled`](https://pkg.go.dev/github.com/reactivego/rx#ErrSubscriptionCanceled) after `Unsubscribe`, or [`ErrSubscriptionActive`](https://pkg.go.dev/github.com/reactivego/rx#ErrSubscriptionActive) while still running.

Every error produced by this package wraps the base error [`Err`](https://pkg.go.dev/github.com/reactivego/rx#Err), so `errors.Is(err, rx.Err)` identifies package errors.

### Context and Schedulers

`Subscribe` takes a `context.Context` that serves two purposes:

1. **Cancellation.** Canceling the context unsubscribes the subscription; conversely the context is available to source operators through the subscriber, so sources can stop producing when it is canceled.
2. **Scheduler selection.** A [`Scheduler`](https://pkg.go.dev/github.com/reactivego/rx#Scheduler) controls where and when the stream runs, and is carried by the context:

```go
// Concurrent: every task runs on its own goroutine.
ctx := rx.GoroutineContext()

// Serial: tasks run sequentially on the calling goroutine (the default).
ctx := rx.SchedulerContextWith(context.Background(), rx.NewScheduler())
```

Cancellation composes with scheduler selection: context values flow down, so a cancelable context derived from a scheduler context still carries the scheduler. Canceling it unsubscribes every subscription made with it:

```go
ctx, cancel := context.WithTimeout(rx.GoroutineContext(), time.Second)
defer cancel()

err := rx.Interval[int](100 * time.Millisecond).Println().Wait(ctx)
// after one second: err is rx.ErrSubscriptionCanceled
```

[`GoroutineContext`](https://pkg.go.dev/github.com/reactivego/rx#GoroutineContext) also accepts an optional parent, so `rx.GoroutineContext(ctx)` attaches the scheduler to an existing context instead.

Convenience entry points choose sensible defaults when no scheduler is attached: [`Wait`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Wait), [`First`](https://pkg.go.dev/github.com/reactivego/rx#Observable.First), [`Last`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Last), and [`Slice`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Slice) use a new serial scheduler, while [`Go`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Go) uses the concurrent Goroutine scheduler.

### Iterators and Channels

Observables interoperate with Go's native streaming constructs in both directions:

- [`Values`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Values) and [`All`](https://pkg.go.dev/github.com/reactivego/rx#Observable.All) expose an Observable as an iterator usable in a `range` loop; `All` pairs each value with its index.
- [`Pull`](https://pkg.go.dev/github.com/reactivego/rx#Pull) and [`Pull2`](https://pkg.go.dev/github.com/reactivego/rx#Pull2) adapt an `iter.Seq`/`iter.Seq2` into an Observable.
- [`Recv`](https://pkg.go.dev/github.com/reactivego/rx#Recv) turns a receive channel into an Observable; [`Send`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Send) forwards emissions into a channel.

```go
for value, err := range rx.From(1, 2, 3).Values(rx.GoroutineContext()) {
    if err != nil {
        break // terminal error; value is the zero value
    }
    fmt.Println(value)
}
```

### Hot vs Cold Observables

- **Hot Observables** emit values regardless of subscription status. Like a live broadcast, values emitted while nobody is subscribed are missed. Examples include system events, mouse movements, or real-time data feeds.

- **Cold Observables** begin emission only when subscribed to, so each subscriber receives the complete sequence from the beginning. Examples include file contents, database queries, or HTTP requests executed on demand.

### Multicasting

By default a cold Observable restarts for every subscriber. Multicasting shares one upstream subscription among many subscribers:

- [`Publish`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Publish) returns a [`Connectable[T]`](https://pkg.go.dev/github.com/reactivego/rx#Connectable): subscribers register first, and the source is subscribed only when [`Connect`](https://pkg.go.dev/github.com/reactivego/rx#Connector.Connect) is called.
- [`Behavior`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Behavior) is like `Publish` but replays the latest value (starting from a seed) to each new subscriber, with conflating, non-backpressuring semantics suited to distributing current state.
- [`RefCount`](https://pkg.go.dev/github.com/reactivego/rx#Connectable.RefCount) turns a Connectable back into an ordinary Observable that connects when the first subscriber arrives and disconnects when the last one leaves.
- [`AutoConnect`](https://pkg.go.dev/github.com/reactivego/rx#Connectable.AutoConnect) connects once a given number of subscribers have subscribed, and reconnects after an error or cancellation when the threshold is crossed again.
- [`Share`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Share) is shorthand for `Publish().RefCount()`.
- [`Multicast`](https://pkg.go.dev/github.com/reactivego/rx#Multicast) and [`Subject`](https://pkg.go.dev/github.com/reactivego/rx#Subject) provide a raw Observer/Observable pair for pushing values imperatively, with configurable buffering and replay.

```go
state := source.Behavior(initial).RefCount()
// every subscriber immediately observes the current value, then tracks updates
```

## Operators

Operators form a language for expressing programs with Observables. They transform, filter, and combine one or more Observables into new Observables, enabling complex asynchronous workflows through method chaining. Most operators exist both as a method on `Observable[T]` and as a standalone function returning a [`Pipe`](https://pkg.go.dev/github.com/reactivego/rx#Pipe) segment for use with [`Pipe`](https://pkg.go.dev/github.com/reactivego/rx#Observable.Pipe); operators that change the element type (like `Map` between arbitrary types) are standalone functions.

### Index

[__All__](https://pkg.go.dev/github.com/reactivego/rx#Observable.All) converts the Observable into an iterator that yields each value paired with its zero-based index.

[__Append__](https://pkg.go.dev/github.com/reactivego/rx#Append) appends each emitted value to a provided slice while forwarding all emissions.

[__AsObservable__](https://pkg.go.dev/github.com/reactivego/rx#AsObservable) converts between a typed `Observable[T]` and an `Observable[any]`.

[__AsObserver__](https://pkg.go.dev/github.com/reactivego/rx#AsObserver) converts between a typed `Observer[T]` and an `Observer[any]`, emitting `ErrTypecastFailed` when a conversion fails.

[__Assign__](https://pkg.go.dev/github.com/reactivego/rx#Assign) stores each emitted value into a provided pointer variable while forwarding all emissions.

[__AutoConnect__](https://pkg.go.dev/github.com/reactivego/rx#Connectable.AutoConnect) connects a Connectable to its source once the specified number of subscribers have subscribed.

[__AutoUnsubscribe__](https://pkg.go.dev/github.com/reactivego/rx#AutoUnsubscribe) unsubscribes from the source automatically as soon as the stream terminates.

[__Behavior__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Behavior) returns a replay-latest, conflating multicast Connectable seeded with an initial value, suited to distributing current state.

[__BufferCount__](https://pkg.go.dev/github.com/reactivego/rx#BufferCount) collects values into slices of a given size, starting a new buffer every `startBufferEvery` values.

[__Catch__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Catch) recovers from an error by continuing with a replacement Observable instead of emitting the error.

[__CatchError__](https://pkg.go.dev/github.com/reactivego/rx#Observable.CatchError) recovers from an error via a selector function that maps the error to a replacement Observable.

[__CombineAll__](https://pkg.go.dev/github.com/reactivego/rx#CombineAll) collects the Observables emitted by a higher-order Observable and combines their latest values into slices.

[__CombineLatest__](https://pkg.go.dev/github.com/reactivego/rx#CombineLatest) emits a slice of the latest values of all inputs whenever any input emits; __CombineLatest2__–__CombineLatest5__ return strongly typed tuples.

[__Concat__](https://pkg.go.dev/github.com/reactivego/rx#Concat) emits all values from each Observable in turn, never overlapping emissions.

[__ConcatAll__](https://pkg.go.dev/github.com/reactivego/rx#ConcatAll) flattens a higher-order Observable by subscribing to each inner Observable only after the previous one completes.

[__ConcatMap__](https://pkg.go.dev/github.com/reactivego/rx#ConcatMap) projects each value to an inner Observable and concatenates their emissions in order.

[__ConcatWith__](https://pkg.go.dev/github.com/reactivego/rx#Observable.ConcatWith) appends other Observables, each starting only after the previous completes.

[__Connect__](https://pkg.go.dev/github.com/reactivego/rx#Connector.Connect) subscribes a Connectable to its source and returns a Subscription; the scheduler can be provided via context.

[__Connectable__](https://pkg.go.dev/github.com/reactivego/rx#Connectable) is an Observable with delayed connection to its source: observers subscribe first, emission starts on Connect.

[__Connector__](https://pkg.go.dev/github.com/reactivego/rx#Connector) is the Connect half of a Connectable.

[__Count__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Count) emits a single value: the number of items the source emitted before completing.

[__Create__](https://pkg.go.dev/github.com/reactivego/rx#Create) constructs an Observable from a Creator function called with an incrementing index, bridging imperative code and the reactive pattern.

[__Creator__](https://pkg.go.dev/github.com/reactivego/rx#Creator) is the function type used by Create: it returns the next value, an error, and a done flag.

[__Defer__](https://pkg.go.dev/github.com/reactivego/rx#Defer) calls a factory to create a fresh Observable for every subscription.

[__Delay__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Delay) time-shifts each emission by a duration.

[__DistinctUntilChanged__](https://pkg.go.dev/github.com/reactivego/rx#Observable.DistinctUntilChanged) suppresses values that equal the previously emitted value.

[__Do__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Do) calls a function for each value passing through the Observable.

[__ElementAt__](https://pkg.go.dev/github.com/reactivego/rx#Observable.ElementAt) emits only the nth value.

[__Empty__](https://pkg.go.dev/github.com/reactivego/rx#Empty) emits no values and terminates normally.

[__EndWith__](https://pkg.go.dev/github.com/reactivego/rx#Observable.EndWith) emits additional values after the source completes.

[__Equal__](https://pkg.go.dev/github.com/reactivego/rx#Equal) returns an equality function for comparable types, for use with DistinctUntilChanged.

[__Err__](https://pkg.go.dev/github.com/reactivego/rx#Err) is the base error joined into every error this package produces; test with `errors.Is`.

[__ExhaustAll__](https://pkg.go.dev/github.com/reactivego/rx#ExhaustAll) flattens a higher-order Observable, ignoring inner Observables that arrive while one is still active.

[__ExhaustMap__](https://pkg.go.dev/github.com/reactivego/rx#ExhaustMap) projects each value to an inner Observable, ignoring source values while the current inner Observable is active.

[__Filter__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Filter) emits only the values that pass a predicate.

[__First__](https://pkg.go.dev/github.com/reactivego/rx#Observable.First) blocks and returns the first emitted value.

[__Fprint__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Fprint), [__Fprintf__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Fprintf), [__Fprintln__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Fprintln) write each value to an `io.Writer` while forwarding it.

[__From__](https://pkg.go.dev/github.com/reactivego/rx#From) creates an Observable from the values passed in.

[__Go__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Go) subscribes on the Goroutine scheduler, discards emissions, and returns a Subscription — useful for side effects.

[__GoroutineContext__](https://pkg.go.dev/github.com/reactivego/rx#GoroutineContext) returns a context carrying the concurrent Goroutine scheduler, optionally wrapping a parent context.

[__Ignore__](https://pkg.go.dev/github.com/reactivego/rx#Ignore) creates an Observer that discards all emissions.

[__Interval__](https://pkg.go.dev/github.com/reactivego/rx#Interval) emits an increasing sequence of numbers spaced by a time interval.

[__Last__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Last) blocks and returns the last emitted value.

[__Map__](https://pkg.go.dev/github.com/reactivego/rx#Map), [__MapE__](https://pkg.go.dev/github.com/reactivego/rx#MapE) transform each value by applying a function; the E variant's function may also return an error.

[__Marshal__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Marshal) marshals each value to `[]byte` using a provided function.

[__Merge__](https://pkg.go.dev/github.com/reactivego/rx#Merge) combines multiple Observables into one by interleaving their emissions.

[__MergeAll__](https://pkg.go.dev/github.com/reactivego/rx#MergeAll) flattens a higher-order Observable by merging all inner Observables concurrently.

[__MergeMap__](https://pkg.go.dev/github.com/reactivego/rx#MergeMap) projects each value to an inner Observable and merges their emissions.

[__MergeWith__](https://pkg.go.dev/github.com/reactivego/rx#Observable.MergeWith) merges other Observables into the source.

[__Multicast__](https://pkg.go.dev/github.com/reactivego/rx#Multicast) returns an Observer/Observable pair that multicasts pushed values to every subscriber, with blocking or dropping backpressure depending on buffer-size sign.

[__Must__](https://pkg.go.dev/github.com/reactivego/rx#Must) returns its value argument, panicking if the accompanying error is non-nil.

[__Never__](https://pkg.go.dev/github.com/reactivego/rx#Never) emits no values and never terminates.

[__Of__](https://pkg.go.dev/github.com/reactivego/rx#Of) emits a single value and completes.

[__OnComplete__](https://pkg.go.dev/github.com/reactivego/rx#Observable.OnComplete), [__OnDone__](https://pkg.go.dev/github.com/reactivego/rx#Observable.OnDone), [__OnError__](https://pkg.go.dev/github.com/reactivego/rx#Observable.OnError), [__OnNext__](https://pkg.go.dev/github.com/reactivego/rx#Observable.OnNext) register callbacks for the corresponding stream events while forwarding all emissions.

[__Passthrough__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Passthrough) forwards all emissions unchanged.

[__Pipe__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Pipe) applies a chain of standalone operator segments to the Observable.

[__Print__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Print), [__Printf__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Printf), [__Println__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Println) print each value to stdout while forwarding it.

[__Publish__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Publish) returns a multicasting Connectable for the source Observable.

[__Pull__](https://pkg.go.dev/github.com/reactivego/rx#Pull), [__Pull2__](https://pkg.go.dev/github.com/reactivego/rx#Pull2) adapt a Go iterator (`iter.Seq`, `iter.Seq2`) into an Observable.

[__Race__](https://pkg.go.dev/github.com/reactivego/rx#Race) mirrors the first Observable to emit or terminate, unsubscribing the others.

[__RaceWith__](https://pkg.go.dev/github.com/reactivego/rx#Observable.RaceWith) races the source against other Observables.

[__Recv__](https://pkg.go.dev/github.com/reactivego/rx#Recv) creates an Observable from a receive channel.

[__Reduce__](https://pkg.go.dev/github.com/reactivego/rx#Reduce), [__ReduceE__](https://pkg.go.dev/github.com/reactivego/rx#ReduceE) accumulate all values into a single final result using a seed and accumulator function.

[__RefCount__](https://pkg.go.dev/github.com/reactivego/rx#Connectable.RefCount) converts a Connectable into an Observable that connects on the first subscription and disconnects when the last subscriber leaves.

[__Repeat__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Repeat) resubscribes to the source each time it completes, indefinitely or a given number of times.

[__Retry__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Retry), [__RetryTime__](https://pkg.go.dev/github.com/reactivego/rx#Observable.RetryTime) resubscribe to the source when it errors, optionally limited and with a backoff schedule.

[__SampleTime__](https://pkg.go.dev/github.com/reactivego/rx#Observable.SampleTime) emits the most recent value within periodic time intervals.

[__Scan__](https://pkg.go.dev/github.com/reactivego/rx#Scan), [__ScanE__](https://pkg.go.dev/github.com/reactivego/rx#ScanE) apply an accumulator to each value, emitting every intermediate result.

[__Send__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Send) sends each value into a provided channel while forwarding it.

[__Share__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Share) multicasts the source with reference counting: `Publish().RefCount()`.

[__Skip__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Skip) suppresses the first n values.

[__Slice__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Slice) blocks and collects all values into a slice.

[__StartWith__](https://pkg.go.dev/github.com/reactivego/rx#Observable.StartWith) synchronously emits the given values on subscription before mirroring the source.

[__Subject__](https://pkg.go.dev/github.com/reactivego/rx#Subject) returns an Observer/Observable pair with a replay buffer bounded by age and capacity.

[__Subscribe__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Subscribe) attaches an Observer with a context and returns a Subscription; the context carries the scheduler and cancels the subscription.

[__SubscribeOn__](https://pkg.go.dev/github.com/reactivego/rx#Observable.SubscribeOn) specifies the concurrent scheduler the Observable should use when subscribed to.

[__Take__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Take) emits only the first n values, then completes.

[__TakeWhile__](https://pkg.go.dev/github.com/reactivego/rx#Observable.TakeWhile) mirrors values until a predicate becomes false.

[__Tap__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Tap) observes all emissions with a secondary Observer without altering the stream.

[__Throw__](https://pkg.go.dev/github.com/reactivego/rx#Throw) emits no values and terminates with the given error.

[__Ticker__](https://pkg.go.dev/github.com/reactivego/rx#Ticker) emits timestamps after an initial delay, following a schedule of intervals.

[__Timer__](https://pkg.go.dev/github.com/reactivego/rx#Timer) emits an increasing sequence of numbers after an initial delay, following a schedule of intervals.

[__Tuple2__](https://pkg.go.dev/github.com/reactivego/rx#Tuple2)–[__Tuple5__](https://pkg.go.dev/github.com/reactivego/rx#Tuple5) are the strongly typed tuples produced by the numbered combining operators.

[__Values__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Values) converts the Observable into an iterator that yields each value, ending with the terminal error if the stream fails.

[__Wait__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Wait) subscribes and blocks until completion, returning nil, the stream's error, or a cancellation error.

[__WithLatestFrom__](https://pkg.go.dev/github.com/reactivego/rx#WithLatestFrom) emits a combination of each source value with the latest values from the other Observables; __WithLatestFrom2__–__WithLatestFrom5__ return strongly typed tuples.

[__WithLatestFromAll__](https://pkg.go.dev/github.com/reactivego/rx#WithLatestFromAll) applies WithLatestFrom to the Observables emitted by a higher-order Observable.

[__Zip__](https://pkg.go.dev/github.com/reactivego/rx#Zip) combines multiple Observables in lockstep, pairing the nth values of every input; __Zip2__–__Zip5__ return strongly typed tuples and accept a buffer-size option.

[__ZipAll__](https://pkg.go.dev/github.com/reactivego/rx#ZipAll) zips the Observables emitted by a higher-order Observable.
