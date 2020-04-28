/*
	Make

Make provides a way of creating an Observable from scratch by
calling observer methods programmatically. A make function conforming to the
MakeFunc signature will be called by Make providing a Next, Error and
Complete function that can be called by the code that implements the
Observable.
*/
package Make

import _ "github.com/reactivego/rx"
