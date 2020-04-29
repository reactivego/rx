/*
Start creates an Observable that emits the return value of a function.

It is designed to be used with a function that returns a (value, error) tuple.
If the error is non-nil the returned Observable will be an Observable that
emits and error, otherwise it will be a single-value Observable of the value.

	Start	http://reactivex.io/documentation/operators/start.html
*/
package Start
