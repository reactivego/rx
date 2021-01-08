# CatchError

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/CatchError?tab=doc)
[![](../../../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx/test/CatchError)
[![](../../../assets/rx.svg?raw=true)](https://rxjs-dev.firebaseapp.com/api/operators/catchError)

**CatchError** catches errors on the Observable to be handled by returning a
new Observable or throwing an error. It is passed a selector function 
that takes as arguments err, which is the error, and caught, which is the
source observable, in case you'd like to "retry" that observable by
returning it again. Whatever observable is returned by the selector will be
used to continue the observable chain.