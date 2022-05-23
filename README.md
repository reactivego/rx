# rx

    import "github.com/reactivego/x"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/x.svg)](https://pkg.go.dev/github.com/reactivego/x)

Package `reactivego/x` provides *Reactive Extensions* for Go, an API for asynchronous programming with [observables](#observables) and [operators](#operators).

## Prerequisites

A version of [Go](https://golang.org/dl/) that supports generics.

## Observables
The main focus of `reativego/x` is [Observables](http://reactivex.io/documentation/observable.html). Observables are dynamic and active at their core, because they are defined as values that change over time. Observables are combined into compound observables with the help of many different operators. This means that at the program level you don't have to deal explicitly with dynamic nature of the observable, but with more tangible static relations between observables as declared by using operators.

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
Output
```
1
hi
2.3
```
> Note the program creates a mixed type `any` observable from an int, string and a float64.


Example
```go
package main

import "github.com/reactivego/x"

func main() {
    x.From(1,2,3).Println()
}
```
Output
```
1
2
3
```
> Note the program uses inferred type int for the observable.

Observables in `reactivego/x` are somewhat similar to Go channels but have much richer
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