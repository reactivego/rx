# rx

    import "github.com/reactivego/rx"

[![](../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx?tab=doc)
[![](../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx)
[![](../assets/rx.svg?raw=true)](http://reactivex.io/intro.html)

Package `rx` provides *Reactive Extensions* for Go, an API for asynchronous programming with observable streams.

> Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed.
> For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.
>
> *Edsger W. Dijkstra*, March 1968

### Installation
Use the go tool to get the package:

```bash
$ go get github.com/reactivego/rx
```

Optionally install the [jig](http://github.com/reactivego/jig) tool, in order to use the package as a generics library.

```bash
$ go get github.com/reactivego/jig
```
The [jig](http://github.com/reactivego/jig) tool provides the parametric polymorphism capability that Go 1 is missing.
It works by generating code, replacing place-holder types in generic functions and datatypes with specific types.

## Usage

This package can be used either directly as a standard package or as a generics library for Go 1.

### Standard Package
To use as a standard package, import the root `rx` package to access the API directly.
```go
package main

import "github.com/reactivego/rx"

func main() {
    rx.From(1,2,"hello").Println()
}
```
Links to documentation are at the top of this page.

### Generics Library
To use as a *Generics Library* for *Go 1*, import the `rx/generic` package.
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

This package uses `interface{}` for entry types, so an observable can emit a
mix of differently typed entries. To create an observable that emits three
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
when nobody is subscribed. As long as there are no subscribers the values of
a hot observable are lost. The position of a mouse pointer or the current time
are examples of hot observables. 

A cold observable will only start emitting values when somebody subscribes.
The contents of a file or a database are examples of cold observables.

An observable can complete normally or with an error, it uses subscriptions
that can be canceled from the subscriber side. Where a normal variable is
just a place where you read and write values from, an observable captures how
the value of this variable changes over time.

Concurrency flows naturally from the fact that an observable is an ever
changing stream of values. Every Observable conceptually has at its core a
concurrently running process that pushes out values.

## Operators 
[Operators](OPERATORS.md) form a language in which programs featuring Observables can be expressed.
They work on one or more Observables to transform, filter and combine them into new Observables.

Following are the most commonly used [operators](OPERATORS.md):

<details><summary>BufferTime</summary>

#### TBD

</details>
<details><summary>Catch</summary>

#### TBD

</details>
<details><summary>CatchError</summary>

#### TBD

</details>
<details><summary>CombineLatest</summary>

#### TBD

</details>
<details><summary>Concat</summary>

#### TBD

</details>
<details><summary>ConcatMap</summary>

#### TBD

</details>
<details><summary>ConcatWith</summary>

#### TBD

</details>
<details><summary>Create</summary>

#### TBD

</details>
<details><summary>DebounceTime</summary>

#### TBD

</details>
<details><summary>DistinctUntilChanged</summary>

#### TBD

</details>
<details><summary>Do</summary>

#### TBD

</details>
<details><summary>Filter</summary>

#### TBD

</details>
<details><summary>From</summary>

#### TBD

</details>
<details><summary>Just</summary>

#### TBD

</details>
<details><summary>Map</summary>

#### TBD

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

#### TBD

</details>
<details><summary>MergeWith</summary>

combines multiple Observables into one by merging their emissions.
An error from any of the observables will terminate the merged observables.

![Merge](../assets/MergeWith.svg?raw=true)

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

#### TBD

</details>
<details><summary>Publish</summary>

#### TBD

</details>
<details><summary>PublishReplay</summary>

#### TBD

</details>
<details><summary>Scan</summary>

#### TBD

</details>
<details><summary>StartWith</summary>

#### TBD

</details>
<details><summary>SwitchMap</summary>

#### TBD

</details>
<details><summary>Take</summary>

#### TBD

</details>
<details><summary>TakeUntil</summary>

#### TBD

</details>
<details><summary>WithLatestFrom</summary>

#### TBD

</details>

## Concurency

The `rx` package does not use any 'bare' go routines internally. Concurrency is tightly controlled by the use of a specific [scheduler](https://github.com/reactivego/scheduler).
Currently there is a choice between 2 different schedulers; a *trampoline* schedulers and a *goroutine* scheduler.

By default all subscribing operators except `ToChan` use the *trampoline* scheduler. A *trampoline* puts tasks on a task queue and only starts processing them when the `Wait` method is called on a returned subscription. The subscription itself calls the `Wait` method of the scheduler.

Only the `Connect` and `Subscribe` methods return a subscription. The other subscribing operators `Println`, `ToSingle`, `ToSlice` and `Wait` are blocking by default and only return when the scheduler returns after wait. In the case of `ToSingle` and `ToSlice` both a value and error are returned and in the case of `Println` and `Wait` just an error is returned or nil when the observable completed succesfully.

The only operator that uses the *goroutine* scheduler by default is the `ToChan` operator. `ToChan` returns a channel that it feeds by subscribing on the *goroutine* scheduler.

To change the scheduler on which subscribing needs to occur, use the [SubscribeOn](OPERATORS.md#subscribeon) operator

## Regenerating this Package
This package is generated from generics in the sub-folder `generic` by the [jig](http://github.com/reactivego/jig) tool.
You don't need to regenerate this package in order to use it. However, if you are
interested in regenerating it, then read on.

To regenerate, change the current working directory to the package directory
and run the [jig](http://github.com/reactivego/jig) tool as follows:

```bash
$ go get -d github.com/reactivego/jig
$ go run github.com/reactivego/jig -v
```

This works by replacing place-holder types of generic functions and datatypes with the `interface{}` type.

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

This implementation takes most of its cues from [RxJS](https://github.com/ReactiveX/rxjs).
It is the ReactiveX incarnation that pushes the envelope in evolving operator semantics.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file for copyright notice and exact wording.
