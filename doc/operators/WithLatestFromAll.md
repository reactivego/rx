# WithLatestFromAll

[![](../../../assets/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/WithLatestFromAll?tab=doc)
[![](../../../assets/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test/WithLatestFromAll)
[![](../../../assets/rx.svg)](https://rxjs-dev.firebaseapp.com/api/operators/withLatestFrom)

**WithLatestFromAll** flattens a higher order observable (e.g. ObservableObservable)
by subscribing to all emitted observables (ie. Observable entries) until the
source completes. It will then wait for all of the subscribed Observables to
emit before emitting the first slice. The first observable that was emitted by
the source will be used as the trigger observable. Whenever the trigger
observable emits, a new slice will be emitted containing all the latest values.

Note that any values emitted by the trigger before all other observables have emitted will
effectively be lost. The first emit will occur the first time the trigger emits after all other
observables have emitted.