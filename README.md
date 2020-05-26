# rx

    import "github.com/reactivego/rx"

[![](https://godoc.org/github.com/reactivego/rx?status.png)](http://godoc.org/github.com/reactivego/rx)

Package `rx` provides [Reactive Extensions](http://reactivex.io/) for Go.
An API for asynchronous programming with observable streams.

## Table of Contents

<!-- MarkdownTOC -->

- [Installation](#installation)
- [Usage](#usage)
- [Observables](#observables)
- [Operators](#operators)
- [Regenerating this Package](#regenerating-this-package)
- [Obligatory Dijkstra Quote](#obligatory-dijkstra-quote)
- [Acknowledgements](#acknowledgements)
- [License](#license)

<!-- /MarkdownTOC -->

### Installation

To use the `rx` package, you just import it in your code:

```go
import "github.com/reactivego/rx"
```
Access the package api directly through the `rx` prefix.

To use the package purely as a generic programming library, first install the [*jig*](https://github.com/reactivego/jig) tool, then import as follows:
```go
import _ "github.com/reactivego/rx"
```
The use of `_` stops the *go* tool from complaining about you not using any code from the package while at the same time allowing the [*jig*](https://github.com/reactivego/jig) tool to actually find the generics in the library.

## Usage
The library can be used in two ways

1. Direct API from `github.com/reactivego/rx`
2. Generics from `github.com/reactivego/rx/generic`

### Direct API
To use the API for asynchronous programming with observables 

You can use this package directly as follows:

	import "github.com/reactivego/rx"

Then you just use the code directly from the library:

	rx.From(1,2,"hello").Println()	

### Generics
Alternatively you can use the library as a generics library and use
a tool to generate statically typed  observables and operators:

	import _ "github.com/reactivego/rx/generic"

Then you use the code as follows:
	
	FromInt(1,2).Println()

You'll need to generate the observables and operators by running the [jig tool](https://github.com/reactivego/jig).

#### Overview
- You write code that references templates from the `rx` library.
- You run the `jig` command in the directory where your code is located.
- **Now `jig` analyzes your code and determines what additional code is needed to make it build**.
- *Jig* takes templates from the `rx` library and specializes them on specific types.
- Specializations are generated into the file `rx.go` alongside your own code.
- If all went well, your code will now build.

#### Preparing the Working Directory
Let's create a new folder for our simple program and start editing the file `main.go`.

```bash
$ mkdir -p ~/go/hellorx
$ cd ~/go/hello/rx
$ subl main.go
```
> We use [*Sublime Text*](https://www.sublimetext.com/) for code editing, hence the use of the `subl` command.

#### Writing Code
Now that you have your `main.go` file open in your editor of choice, type the following code:

```go
package main

import _ "github.com/reactivego/rx"

func main() {
	FromString("You!", "Gophers!", "World!").
		MapString(func(x string) string {
			return "Hello, " + x
		}).
		Println()
}
```

If you go to the command-line now and run `jig` it would do the following:

1. Import `github.com/reactivego/rx/generic` to access the generics in that library.
2. Generate `FromString` by specializing the template `From<Foo>` from the library on type `string`. Generate dependencies of `FromString` like the type `ObservableString` that is returned by `FromString`.
3. Generate the `MapString` method of `ObservableString` by specializing the template `Observable<Foo> Map<Bar>` for `Foo` as type `string` and `Bar` (the type to map to) also as type `string`.
4. Map function from `string` to `string` just concatenates two strings and returns the result.
5. Print every string returned by the Map function.
6. The output you can expect when you run the program.

#### Generating Code

Now actually go to the command-line and run `jig -v`. Use the verbose flag `-v` because otherwise `jig` will be silent.

```bash
$ jig -v
found 126 templates in package "rx" (github.com/reactivego/rx/generic)
found 16 templates in package "multicast" (github.com/reactivego/multicast/generic)
generating "FromString"
  Scheduler
  Subscriber
  StringObserveFunc
  zeroString
  ObservableString
  FromSliceString
  FromString
generating "ObservableString MapString"
  ObservableString MapString
generating "ObservableString Println"
  Schedulers
  ObservableString Println
writing file "rx.go"
```

#### Running the Program

Now we can try to run the code and see what it does.

```bash
$ go run *.go
Hello, You!
Hello, Gophers!
Hello, World!
```

Success! `jig` generated the code into the file `rx.go` and we were able to run the program.
Turns out the generated file [`rx.go`](../example/rx/rx.go) contains less than 125 lines of code.

If you add additional code to the program that uses different generics of the `rx` library, then you should run `jig` again to generate specializations of those generics.

## Observables

The main focus of `rx` is on [Observables](http://reactivex.io/documentation/observable.html).

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

The code above creates an observable from numbers and strings and then prints them.

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

Concurrency follows naturally from the fact that an observable is an ever
changing stream of values.

## Operators 

The combination of Observables and a set of expressive Operators is the real strength of Reactive Extensions. 
Operators work on one or more Observables. They are the language you use to describe the way in which observables should be combined to form new Observables. Operators specify how Observables representing streams of values are e.g. merged, transformed, concatenated, split, multicasted, replayed, delayed and debounced.

This implementation takes most of its cues from [RxJS 6](https://github.com/ReactiveX/rxjs) and [RxJava 2](https://github.com/ReactiveX/RxJava). Both libaries have been pushing the envelope in evolving operator semantics.

A list of implemented [operators](http://reactivex.io/documentation/operators.html) can be found in the [test directory](test).

## Regenerating this Package

This package is generated from the sub-folder generic by the [jig](http://github.com/reactivego/jig) tool.
You don't need to regenerate the package in order to use it. However, if you are
interested in regenerating it, then read on.

The jig tool provides the parametric polymorphism capability that Go 1 is missing.
It works by replacing place-holder types of generic functions and datatypes
with interface{} (it can also generate statically typed code though).

To regenerate, change the current working directory to the package directory
and run the jig tool as follows:

```bash
$ go get -d github.com/reactivego/jig
$ go run github.com/reactivego/jig -v
```

## Obligatory Dijkstra Quote

Our intellectual powers are rather geared to master static relations and our powers to visualize processes evolving in time are relatively poorly developed. For that reason we should do our utmost to shorten the conceptual gap between the static program and the dynamic process, to make the correspondence between the program (spread out in text space) and the process (spread out in time) as trivial as possible.

*Edsger W. Dijkstra*, March 1968

## Acknowledgements
This library started life as the [Reactive eXtensions for Go](https://github.com/alecthomas/gorx) library by *Alec Thomas*. Although the library has been through the metaphorical meat grinder a few times, its DNA is still clearly present in this library and I owe Alec a debt of grattitude for the work he has made so generously available.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file for copyright notice and exact wording.
