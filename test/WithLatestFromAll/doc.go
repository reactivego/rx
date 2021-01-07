/*
WithLatestFromAll flattens a higher order observable (e.g. ObservableObservable) by subscribing
to all emitted observables (ie. Observable entries) until the source completes. It will then wait
for all of the subscribed Observables to emit before emitting the first slice. The first observable
that was emitted by the source will be used as the trigger observable. Whenever the trigger
observable emits, a new slice will be emitted containing all the latest values.

Note that any values emitted by the source before all other observables have emitted will
effectively be lost. The first emit will occur the first time the source emits after all other
observables have emitted.
*/
package WithLatestFromAll
