/*
FromChan creates an Observable from a Go channel.

The feeding code can send zero or more items and then closing the channel will
be seen as completion.

FromChan without a specific type also accepts an error to be fed into the
channel. When the feeding code does that, it should then also close the
channel to indicate termination with error.

	FromChan	http://reactivex.io/documentation/operators/from.html
*/
package FromChan
