/*
CreateFutureRecursive provides a way of creating an Observable from
scratch by calling observer methods programmatically.

The create function provided to CreateFutureRecursive will be called
repeatedly to implement the observable. It is provided with a Next, Error
and Complete function that can be called by the code that implements the
Observable.

The timeout passed in determines the time before calling the create
function. The time.Duration returned by the create function determines how
long CreateFutureRecursive has to wait before calling the create function
again.

	CreateFutureRecursive
*/
package CreateFutureRecursive

