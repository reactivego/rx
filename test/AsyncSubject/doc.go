// AsyncSubject emits the last value (and only the last value) emitted by the
// Observable part, and only after that Observable part completes. (If the
// Observable part does not emit any values, the AsyncSubject also completes
// without emitting any values.)
//
// It will also emit this same final value to any subsequent observers.
// However, if the Observable part terminates with an error, the AsyncSubject
// will not emit any items, but will simply pass along the error notification
// from the Observable part.
package AsyncSubject
