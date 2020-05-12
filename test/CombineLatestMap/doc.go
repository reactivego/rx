/*
CombinesLatestMap maps every entry emitted by the Observable into an
Observable, and then subscribe to it, until the source observable
completes. It will then wait for all of the Observables to emit before
emitting the first slice. Whenever any of the subscribed observables emits,
a new slice will be emitted containing all the latest value.

	CombineLatestMap
*/
package CombineLatestMap
