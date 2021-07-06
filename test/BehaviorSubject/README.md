# BehaviorSubject

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/BehaviorSubject#section-documentation)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/subject.html)

When an observer subscribes to a `BehaviorSubject`, it begins by emitting the item most
recently emitted by the Observable part of the subject (or a seed/default
value if none has yet been emitted) and then continues to emit any other
items emitted later by the Observable part.
