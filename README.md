# rx

    import "github.com/reactivego/rx"

[![](https://godoc.org/github.com/reactivego/rx?status.png)](http://godoc.org/github.com/reactivego/rx)

Package `rx` offers Reactive Extensions ([ReactiveX](http://reactivex.io/)) for [Go](https://golang.org/), an API for asynchronous programming with observable streams (Observables).

## What rx is

Observables are the main focus of rx and they are (sort of) sets of events.
- They assume zero to many values over time.
- They push values.
- They can take any amount of time to complete (or may never).
- They are cancellable.
- They are lazy; they don't do anything until you subscribe.

This implementation of rx uses `interface{}` for value types, so you can
mix different types of values in function and method calls. To create an
observable you might write the following:

```go
From(1,"hi",2.3).Println()
```

The call above creates an observable from numbers and strings and then prints
them.

For an overview of all implemented operators, see [Operators](generic/README.md#operators)

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file in this repository for copyright notice and exact wording.
