/*
Operator Single documentation and tests.

	Single

Single enforces that the observable sends exactly one data item and then
completes. If the observable sends no data before completing or sends more than
1 item before completing  this reported as an error to the observer.
*/
package Single

import _ "github.com/reactivego/rx"

//jig:file {{.package}}.go
