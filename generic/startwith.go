package rx

//jig:template Observable<Foo> StartWith
//jig:needs From<Foo>, Observable<Foo> ConcatWith

// StartWith returns an observable that, at the moment of subscription, will
// synchronously emit all values provided to this operator, then subscribe to
// the source and mirror all of its emissions to subscribers.
func (o ObservableFoo) StartWith(values ...foo) ObservableFoo {
	return FromFoo(values...).ConcatWith(o)
}
