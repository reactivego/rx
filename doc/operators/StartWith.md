# StartWith

[![](../../../assets/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/StartWith?tab=doc)
[![](../../../assets/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test/StartWith)
[![](../../../assets/rx.svg)](http://reactivex.io/documentation/operators/startwith.html)

**StartWith** returns an observable that, at the moment of subscription, will
synchronously emit all values provided to this operator, then subscribe to
the source and mirror all of its emissions to subscribers.
