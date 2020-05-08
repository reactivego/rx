/*
CreateRecursive provides a way of creating an Observable from
scratch by calling observer methods programmatically.

The create function provided to CreateRecursive will be called
repeatedly to implement the observable. It is provided with a Next, Error
and Complete function that can be called by the code that implements the
Observable.

	CreateRecursive
*/
package CreateRecursive
