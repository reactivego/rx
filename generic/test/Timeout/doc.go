/*
Operator Timeout documentation and tests.

	Timeout	http://reactivex.io/documentation/operators/timeout.html

Timeout mirrors the source Observable, but issue an error notification if a
particular period of time elapses without any emitted items.

The Timeout operator allows you to abort an Observable with an onError
termination if that Observable fails to emit any items during a specified
span of time.
*/
package Timeout

import _ "github.com/reactivego/rx"
