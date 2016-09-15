# TODO

## Miscellaneous
Use observable.SubscribeOn(sched).Subscribe(operator).Wait() in Concat implementation.
Implement Subject<T> to be used alongside Create<T> where you call OnNext/OnError/OnComplete directly on the instance of the Subject.
Rewrite Fork test into Publish test and remove the Fork method.
Build test cases using early unsubscribe to test for proper cleanup of internal goroutines.
Switch to alternative interpretation of ObserverFunc using 'done' where 'done' indicates 'error or completed'.
Use OnNext,OnError,OnComplete as names for <T>Observer interface.

## Scheduler
Pass scheduler and unsubscriber to the filter implementation.
Test complex situations with the Immediate scheduler i.e. Publish
Implement timed tasks use them in Timeout, Interval implementations.
Create test cases combining scheduler Go and Do with SubscribeOn and ObserveOn to look for nice conncurrency patterns.
Build a leak detector for goroutines, by never calling go directly, always via wrapper call that also intercepts panics.
Integrate asynchronous CurrentThread Scheduler based on SerialQueue implementation.
Fix Debounce and Sample that now use goroutines directly to use the scheduler.
Fix situations where Unsubscribe may cause goroutines to leak in filters.go

