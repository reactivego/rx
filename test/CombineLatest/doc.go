/*
CombineLatest will subscribe to all Observables. It will then wait for all of
them to emit before emitting the first slice. Whenever any of the subscribed
observables emits, a new slice will be emitted containing all the latest value.
*/
package CombineLatest
