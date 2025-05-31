# rx

    import "github.com/reactivego/rx"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/rx.svg)](https://pkg.go.dev/github.com/reactivego/rx)

Package `rx` provides _**R**eactive E**x**tensions_, a powerful API for asynchronous programming in Go, built around [observables](#observables) and [operators](#operators) to process streams of data seamlessly.

## Prerequisites

You’ll need [*Go 1.23*](https://golang.org/dl/) or later, as the implementation depends on language support for generics and iterators.

## Observables

In `rx`, an [**Observables**](http://reactivex.io/documentation/observable.html) represents a stream of data that can emit items over time, while an **Observer** subscribes to it to receive and react to those emissions. This reactive approach enables asynchronous and concurrent operations without blocking execution. Instead of waiting for values to become available, an observer passively listens and responds whenever the observable emits data, errors, or a completion signal.

This page introduces the **reactive pattern**, explaining what **Observables** and **Observers** are and how subscriptions work. Other sections explore the powerful set of [**Operators**](https://reactivex.io/documentation/operators.html) that allow you to transform, combine, and control data streams efficiently.

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
    x.From[any](1,"hi",2.3).Println().Wait()
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
    rx.From(1,2,3).Println().Wait()
}
```
> Note the program uses inferred type `int` for the observable.

Output
```
1
2
3
```

Observables in `rx` offer several advantages over standard Go channels:

### Hot vs Cold Observables

- **Hot Observables** emit values regardless of subscription status. Like a live broadcast, any values emitted when no subscribers are listening are permanently missed. Examples include system events, mouse movements, or real-time data feeds.

- **Cold Observables** begin emission only when subscribed to, ensuring subscribers receive the complete data sequence from the beginning. Examples include file contents, database queries, or HTTP requests that are executed on-demand.

### Rich Lifecycle Management

Observables offer comprehensive lifecycle handling. They can complete normally, terminate with errors, or continue indefinitely. Subscriptions provide fine-grained control, allowing subscribers to cancel at any point, preventing resource leaks and unwanted processing.

### Time-Varying Data Model

Unlike traditional variables that represent static values, Observables elegantly model how values evolve over time. They represent the entire progression of a value's state changes, not just its current state, making them ideal for reactive programming paradigms.

### Native Concurrency Support

Concurrency is built into the Observable paradigm. Each Observable conceptually operates as an independent process that asynchronously pushes values to subscribers. This approach naturally aligns with concurrent programming models while abstracting away much of the complexity typically associated with managing concurrent operations.

## Operators

Operators form a language for expressing programs with Observables. They transform, filter, and combine one or more Observables into new Observables, allowing for powerful data stream processing. Each operator performs a specific function in the reactive pipeline, enabling you to compose complex asynchronous workflows through method chaining.

### Index

[__All__](https://pkg.go.dev/github.com/reactivego/rx#Observable.All) converts an Observable stream into a Go 1.22+ iterator sequence that provides each emitted value paired with its sequential zero-based index

[__All2__](https://pkg.go.dev/github.com/reactivego/rx#All2)
converts an Observable of Tuple2 pairs into a Go 1.22+ iterator sequence that yields each tuple's components (First, Second) as separate values.

[__Append__](https://pkg.go.dev/github.com/reactivego/rx#Append)
creates a pipe that appends emitted values to a provided slice while forwarding them to the next observer, with a method variant available for chaining.

[__AsObservable__](https://pkg.go.dev/github.com/reactivego/rx#AsObservable)
provides type conversion between observables, allowing you to safely cast an Observable of one type to another, and to convert a typed Observable to an Observable of 'any' type (and vice versa).

[__AsObserver__](https://pkg.go.dev/github.com/reactivego/rx#AsObserver) converts an Observer of type `any` to an Observer of a specific type T.

[__Assign__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Assign) stores each emitted value from an Observable into a provided pointer variable while passing all emissions through to the next observer, enabling value capture during stream processing.

[__AutoConnect__](https://pkg.go.dev/github.com/reactivego/rx#Connectable.AutoConnect) makes a (Connectable) Multicaster behave like an ordinary Observable that automatically connects the mullticaster to its source when the specified number of observers have subscribed to it.

[__AutoUnsubscribe__](https://pkg.go.dev/github.com/reactivego/rx#Observable.AutoUnsubscribe)

[__BufferCount__](https://pkg.go.dev/github.com/reactivego/rx#BufferCount)

[__Catch__](https://pkg.go.dev/github.com/reactivego/rx#Observable.Catch) recovers from an error notification by continuing the sequence without emitting the error but switching to the catch ObservableInt to provide items.

[__CatchError__](https://pkg.go.dev/github.com/reactivego/rx#Observable.CatchError)  catches errors on the Observable to be handled by returning a new Observable or throwing error.

__CombineAll__

__CombineLatest__ combines multiple Observables into one by emitting an array containing the latest values from each source whenever any input Observable emits a value, with variants (__CombineLatest2__, __CombineLatest3__, __CombineLatest4__, __CombineLatest5__) that return strongly-typed tuples for 2-5 input Observables respectively.

__Concat__ combines multiple Observables sequentially by emitting all values from the first Observable before proceeding to the next one, ensuring emissions never overlap.

__ConcatAll__ transforms a higher-order Observable (an Observable that emits other Observables) into a first-order Observable by subscribing to each inner Observable only after the previous one completes.

__ConcatMap__ projects each source value to an Observable, subscribes to it, and emits its values, waiting for each one to complete before processing the next source value.

__ConcatWith__ extends an Observable by appending additional Observables, ensuring that emissions from each Observable only begin after the previous one completes.

__Connectable__ is an Observable with delayed connection to its source, combining both Observable and Connector interfaces. It separates the subscription process into two parts: observers can register via Subscribe, but the Observable won't subscribe to its source until Connect is explicitly called. This enables multiple observers to subscribe before any emissions begin (multicast behavior), allowing a single source Observable to be efficiently shared among multiple consumers. Besides inheriting all methods from Observable and Connector, Connectable provides the convenience methods __AutoConnect__ and __RefCount__ to manage connection behavior.

__Connect__ establishes a connection to the source Observable and returns a Subscription that can be used to cancel the connection when no longer needed.

__Connector__ provides a mechanism for controlling when a Connectable Observable subscribes to its source, allowing you to connect the Observable independently from when observers subscribe to it. This separation enables multiple subscribers to prepare their subscriptions before the source begins emitting items. It has a single method __Connect__.

__Constraints__ type constraints __Signed__, __Unsigned__, __Integer__ and __Float__ copied verbatim from `golang.org/x/exp` so we could drop the dependency on that package.

__Count__ returns an Observable that emits a single value representing the total number of items emitted by the source Observable before it completes.

__Create__ constructs a new Observable from a Creator function, providing a bridge between imperative code and the reactive Observable pattern. The Observable will continue producing values until the Creator signals completion, the Observer unsubscribes, or the Creator returns an error.

__Creator__ is a function type that generates values for an Observable stream. It receives a zero-based index for the current iteration and returns a tuple containing the next value to emit, any error that occurred, and a boolean flag indicating whether the sequence is complete.

__Defer__

__Delay__

__DistinctUntilChanged__ only emits when the current value is different from the last.

__Do__ calls a function for each next value passing through the observable.

__ElementAt__ emit only item n emitted by an Observable.

__Empty__ creates an Observable that emits no items but terminates normally.

__EndWith__

__Equal__

__Err__

__ExhaustAll__

__ExhaustMap__

__Filter__ emits only those items from an observable that pass a predicate test.

__First__ emits only the first item from an Observable.

__Fprint__

__Fprintf__

__Fprintln__

__From__ creates an observable from multiple values passed in.

__Go__ subscribes to the observable and starts execution on a separate goroutine, ignoring all emissions from the observable sequence. This makes it useful when you only care about side effects and not the actual values. Returns a Subscription that can be used to cancel the subscription when no longer needed.

__Ignore[T]__ creates an Observer[T] that simply discards any emissions from an Observable. It is useful when you need to create an Observer but don't care about its values.

__Interval__ creates an ObservableInt that emits a sequence of integers spaced by a particular time
terval.

__Last__ emits only the last item emitted by an Observable.

__Map__ transforms the items emitted by an Observable by applying a function to each item.

__MapE__

__Marshal__

__MaxBufferSizeOption__, __WithMaxBufferSize__

__Merge__ combines multiple Observables into one by merging their emissions.

__MergeAll__ flattens a higher order observable by merging the observables it emits.

__MergeMap__ transforms the items emitted by an Observable by applying a function to each item an
turning an Observable.

__MergeWith__ combines multiple Observables into one by merging their emissions.

__Multicast__

__Must__

__Never__ creates an Observable that emits no items and does't terminate.

__Observable__

__Observer__

__Of__ emits a variable amount of values in a sequence and then emits a complete notification.

__OnComplete__

__OnDone__

__OnError__

__OnNext__

__Passthrough__ just passes through all output from the Observable.

__Pipe__

__Print__

__Printf__

__Println__ subscribes to the Observable and prints every item to os.Stdout.

__Publish__ returns a multicasting Observable[T] for an underlying Observable[T] as a Connectable[T] type.

__Pull__

__Pull2__

__Race__

__RaceWith__

__Recv__

__Reduce__ applies a reducer function to each item emitted by an Observable and the previous reducer
sult.

__ReduceE__

__RefCount__ makes a Connectable behave like an ordinary Observable.

__Repeat__ creates an observable that emits a sequence of items repeatedly.

__Retry__ if a source Observable sends an error notification, resubscribe to it in the hopes that it
ll complete without error.

__RetryTime__

__SampleTime__ emits the most recent item emitted by an Observable within periodic time intervals.

__Scan__ applies a accumulator function to each item emitted by an Observable and the previous
cumulator result.

__ScanE__

__Scheduler__

__Send__

__Share__

__Skip__ suppresses the first n items emitted by an Observable.

__Slice__

__StartWith__ returns an observable that, at the moment of subscription, will synchronously emit all values provided to this operator, then subscribe to the source and mirror all of its emissions to subscribers.

__Subject__ is a combination of an observer and observable.

__Subscribe__ operates upon the emissions and notifications from an Observable.

__SubscribeOn__ specifies the scheduler an Observable should use when it is subscribed to.

__Subscriber__

__Subscription__

__SwitchAll__

__SwitchMap__

__Take__ emits only the first n items emitted by an Observable.

__TakeWhile__ mirrors items emitted by an Observable until a specified condition becomes false.

__Tap__

__Throw__ creates an observable that emits no items and terminates with an error.

__Ticker__ creates an ObservableTime that emits a sequence of timestamps after an initialDelay has passed.

__Timer__ creates an Observable that emits a sequence of integers (starting at zero) after an initialDelay has passed.

__Tuple__

__Values__

__Wait__ subscribes to the Observable and waits for completion or error.

__WithLatestFrom__ will subscribe to all Observables and wait for all of them to emit before emitting the first slice.

__WithLatestFromAll__ flattens a higher order observable.

__Zip__

__ZipAll__
