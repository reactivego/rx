/*
CombineLatestAll flattens a higher order observable (e.g. ObservableObservable)
by subscribing to all emitted observables (ie. Observable entries) until the
source completes. It will then wait for all of the subscribed Observables to
emit before emitting the first slice. Whenever any of the subscribed
observables emits, a new slice will be emitted containing all the latest value.
*/
package CombineLatestAll

type Vector struct{ x, y int }
