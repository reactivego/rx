package rx

// Share returns a new Observable that multicasts (shares) the original
// Observable. As long as there is at least one Subscriber this Observable
// will be subscribed and emitting data. When all subscribers have unsubscribed
// it will unsubscribe from the source Observable. Because the Observable is
// multicasting it makes the stream hot.
//
// This method is useful when you have an Observable that is expensive to
// create or has side-effects, but you want to share the results of that
// Observable with multiple subscribers. By using `Share`, you can avoid
// creating multiple instances of the Observable and ensure that all
// subscribers receive the same data.
func (observable Observable[T]) Share() Observable[T] {
	return observable.Publish().RefCount()
}
