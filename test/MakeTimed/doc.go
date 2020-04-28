/*
	MakeTimed

MakeTimed provides a way of creating an Observable from scratch by
calling observer methods programmatically. A make function conforming to the
MakeTimedFunc signature will be called by MakeTimed provining a Next, Error and
Complete function that can be called by the code that implements the
Observable. The timeout passed in determines the time between calling the
make function. The time.Duration returned by MakeTimedFunc determines when
to reschedule the next iteration.
*/
package MakeTimed

import _ "github.com/reactivego/rx"
