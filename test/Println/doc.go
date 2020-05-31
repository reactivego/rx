/*
Println subscribes to the Observable and prints every item to os.Stdout
while it waits for completion or error.

Returns either the error or nil when the Observable completed normally.
Println is performed on the Trampoline scheduler.
*/
package Println
