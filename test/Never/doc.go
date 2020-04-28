/*
	Never	http://reactivex.io/documentation/operators/empty-never-throw.html

Never creates an Observable that emits no items and does't terminate.
It does not schedule a task and as such you won't be able to use a subscriber's
Wait method for it to terminate if you use the operator by itself. Combine
Never with e.g. Timeout to make it usefull.
*/
package Never

import _ "github.com/reactivego/rx"
