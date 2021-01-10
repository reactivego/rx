# rx

    import "github.com/reactivego/rx"

[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/intro.html)

Package `rx` provides *Reactive Extensions* for Go, an API for asynchronous programming with [observables](#observables) and [operators](#operators).

> Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed.
> For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.
>
> *Edsger W. Dijkstra*, March 1968

## Prerequisites

A working [Go](https://golang.org/dl/) environment on your system.

## Installation

You can use this package [directly](#standard-package) without installing anything.

Optionally, install the [jig](http://github.com/reactivego/jig) tool in order to [regenerate](#regenerating-this-package) this package.
The [jig](http://github.com/reactivego/jig) tool is also needed to support using the `rx/generic` sub-package as a generics library.

```bash
$ go get github.com/reactivego/jig
```
The [jig](http://github.com/reactivego/jig) tool provides the parametric polymorphism capability that Go 1 is missing.
It works by generating code, replacing place-holder types in generic functions and datatypes with specific types.

## Usage

This package can be used either directly as a standard package or as a generics library for Go 1.

### Standard Package
To use it as a standard package, import the root `rx` package to access the API directly.
```go
package main

import "github.com/reactivego/rx"

func main() {
    rx.From(1,2,"hello").Println()
}
```
For more information see the [Go Package Documentation](https://pkg.go.dev/github.com/reactivego/rx?tab=doc). In this document there is a [selection](#operators) of the most used operators and there is also a complete list of operators in a separate [document](OPERATORS.md).

### Generics Library
To use it as a *Generics Library* for *Go 1*, import the `rx/generic` package.
Generic programming supports writing programs that use statically typed *Observables* and *Operators*.
For example:

```go
package main

import _ "github.com/reactivego/rx/generic"

func main() {
	FromString("Hello", "Gophers!").Println()
}
```
> Note that From**String** is statically typed.

Generate statically typed code by running the [jig](http://github.com/reactivego/jig) tool.

```bash
$ jig -v
found 149 templates in package "rx" (github.com/reactivego/rx/generic)
...
```
Details on how to use generics can be found in the [doc](doc) folder.

## Observables
The main focus of `rx` is [Observables](http://reactivex.io/documentation/observable.html).
In accordance with Dijkstra's observation that we are ill-equipped to visualize how dynamic processes evolve over time, the use of Observables indeed narrows the conceptual gap between the static program and the dynamic process. Observables have a dynamic process at their core because they are defined as values which change over time. They are combined in static relations with the help of many different operators. This means that at the program level (spread out in the text space) you no longer have to deal explicitly with dynamic processes, but with more tangible static relations.


An Observable:

- is a stream of events.
- assumes zero to many values over time.
- pushes values
- can take any amount of time to complete (or may never)
- is cancellable
- is lazy (it doesn't do anything until you subscribe).

This package uses `interface{}` for value types, so an observable can emit a
mix of differently typed values. To create an observable that emits three
values of different types you could write the following little program.

```go
package main

import "github.com/reactivego/rx"

func main() {
    rx.From(1,"hi",2.3).Println()
}
```
> Note the program creates a mixed type observable from an int, string and a float.

Observables in `rx` are somewhat similar to Go channels but have much richer
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

Following are the most commonly used operators:

<details><summary>BufferTime</summary>

buffers the source Observable values for a specific time period and emits those as a
slice periodically in time.

![BufferTime](../assets/BufferTime.svg?raw=true)

</details>
<details><summary>Catch</summary>

recovers from an error notification by continuing the sequence without
emitting the error but by switching to the catch ObservableInt to provide
items.

</details>
<details><summary>CatchError</summary>

catches errors on the Observable to be handled by returning a
new Observable or throwing an error. It is passed a selector function 
that takes as arguments err, which is the error, and caught, which is the
source observable, in case you'd like to "retry" that observable by
returning it again. Whatever observable is returned by the selector will be
used to continue the observable chain.

</details>
<details><summary>CombineLatest</summary>

will subscribe to all Observables. It will then wait for all of
them to emit before emitting the first slice. Whenever any of the subscribed
observables emits, a new slice will be emitted containing all the latest value.

</details>
<details><summary>Concat</summary>

emits the emissions from two or more observables without interleaving them.

</details>
<details><summary>ConcatMap</summary>

transforms the items emitted by an Observable by applying a
function to each item and returning an Observable. The stream of
Observable items is then flattened by concattenating the emissions from
the observables without interleaving.

</details>
<details><summary>ConcatWith</summary>

emits the emissions from two or more observables without interleaving them.

![ConcatWith](../assets/ConcatWith.svg?raw=true)

</details>
<details><summary>Create</summary>

provides a way of creating an Observable from scratch by calling observer 
methods programmatically.

The create function provided to Create will be called once
to implement the observable. It is provided with a Next, Error,
Complete and Canceled function that can be called by the code that
implements the Observable.

</details>
<details><summary>DebounceTime</summary>

only emits the last item of a burst from an Observable if a
particular timespan has passed without it emitting another item.

</details>
<details><summary>DistinctUntilChanged</summary>

#### TBD

</details>
<details><summary>Do</summary>

calls a function for each next value passing through the observable.

</details>
<details><summary>Filter</summary>

emits only those items from an observable that pass a predicate test.

</details>
<details><summary>From</summary>

creates an observable from multiple values passed in.

</details>
<details><summary>Just</summary>

creates an observable that emits a particular item.

</details>
<details><summary>Map</summary>

transforms the items emitted by an Observable by applying a
function to each item.

</details>
<details><summary>Merge</summary>

combines multiple Observables into one by merging their emissions.
An error from any of the observables will terminate the merged observables.

![Merge](../assets/Merge.svg?raw=true)

Code:
```go
a := rx.From(0, 2, 4)
b := rx.From(1, 3, 5)
rx.Merge(a, b).Println()
```
Output:
```
0
1
2
3
4
5
```

</details>
<details><summary>MergeMap</summary>

transforms the items emitted by an Observable by applying a
function to each item an returning an Observable. The stream of Observable
items is then merged into a single stream of items using the MergeAll
operator.

This operator was previously named FlatMap. The name FlatMap is deprecated as
MergeMap more accurately describes what the operator does with the observables
returned from the Map project function.

</details>
<details><summary>MergeWith</summary>

combines multiple Observables into one by merging their emissions.
An error from any of the observables will terminate the merged observables.

![MergeWith](../assets/MergeWith.svg?raw=true)

Code:
```go
a := rx.From(0, 2, 4)
b := rx.From(1, 3, 5)
a.MergeWith(b).Println()
```
Output:
```
0
1
2
3
4
5
```
</details>
<details><summary>Of</summary>

emits a variable amount of values in a sequence and then emits a complete
notification.

</details>
<details><summary>Publish</summary>

returns a Multicaster for a Subject to an underlying Observable
and turns the subject into a connnectable observable.

A Subject emits to an observer only those items that are emitted by the
underlying Observable subsequent to the time the observer subscribes.

When the underlying Obervable terminates with an error, then subscribed
observers will receive that error. After all observers have unsubscribed due
to an error, the Multicaster does an internal reset just before the next
observer subscribes. So this Publish operator is re-connectable, unlike the
RxJS 5 behavior that isn't. To simulate the RxJS 5 behavior use
Publish().AutoConnect(1) this will connect on the first subscription but will
never re-connect.

</details>
<details><summary>PublishReplay</summary>

returns a Multicaster for a ReplaySubject to an underlying
Observable and turns the subject into a connectable observable. A
ReplaySubject emits to any observer all of the items that were emitted by
the source observable, regardless of when the observer subscribes. When the
underlying Obervable terminates with an error, then subscribed observers
will receive that error. After all observers have unsubscribed due to an
error, the Multicaster does an internal reset just before the next observer
subscribes.

</details>
<details><summary>Scan</summary>

applies a accumulator function to each item emitted by an Observable and
the previous accumulator result.

The operator accepts a seed argument that is passed to the accumulator for the
first item emitted by the Observable. Scan emits every value, both intermediate
and final.

</details>
<details><summary>StartWith</summary>

returns an observable that, at the moment of subscription, will
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
</details>
<details><summary>SwitchMap</summary>

transforms the items emitted by an Observable by applying a function
to each item an returning an Observable.

In doing so, it behaves much like MergeMap (previously FlatMap), except that
whenever a new Observable is emitted SwitchMap will unsubscribe from the
previous Observable and begin emitting items from the newly emitted one.

</details>
<details><summary>Take</summary>

emits only the first n items emitted by an Observable.

</details>
<details><summary>TakeUntil</summary>

emits items emitted by an Observable until another Observable emits an item.

</details>
<details><summary>WithLatestFrom</summary>

will subscribe to all Observables and wait for all of them to emit before emitting
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
</details>

There is also a complete list of supported [operators](OPERATORS.md).

## Concurency

The `rx` package does not use any 'bare' go routines internally. Concurrency is tightly controlled by the use of a specific [scheduler](https://github.com/reactivego/scheduler).
Currently there is a choice between 2 different schedulers; a *trampoline* schedulers and a *goroutine* scheduler.

By default all subscribing operators except `ToChan` use the *trampoline* scheduler. A *trampoline* puts tasks on a task queue and only starts processing them when the `Wait` method is called on a returned subscription. The subscription itself calls the `Wait` method of the scheduler.

Only the `Connect` and `Subscribe` methods return a subscription. The other subscribing operators `Println`, `ToSingle`, `ToSlice` and `Wait` are blocking by default and only return when the scheduler returns after wait. In the case of `ToSingle` and `ToSlice` both a value and error are returned and in the case of `Println` and `Wait` just an error is returned or nil when the observable completed succesfully.

The only operator that uses the *goroutine* scheduler by default is the `ToChan` operator. `ToChan` returns a channel that it feeds by subscribing on the *goroutine* scheduler.

To change the scheduler on which subscribing needs to occur, use the [SubscribeOn](OPERATORS.md#subscribeon) operator

## Regenerating this Package
This package is generated from generics in the sub-folder `generic` by the [jig](http://github.com/reactivego/jig) tool.
You don't need to regenerate this package in order to use it. However, if you are interested in regenerating it, then read on.

If not already done so, install the [jig](http://github.com/reactivego/jig) tool.

```bash
$ go get github.com/reactivego/jig
```
Change the current working directory to the package directory and run the [jig](http://github.com/reactivego/jig) tool as follows:
```bash
$ jig -v
```
This works by replacing place-holder types of generic functions and datatypes with the `interface{}` type.

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

This implementation takes most of its cues from [RxJS](https://github.com/ReactiveX/rxjs).
It is the ReactiveX incarnation that pushes the envelope in evolving operator semantics.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file for copyright notice and exact wording.
