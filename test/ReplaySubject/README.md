# ReplaySubject

[![](../../svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/rx/test/ReplaySubject?tab=doc)
[![](../../svg/godoc.svg)](https://godoc.org/github.com/reactivego/rx/test/ReplaySubject)
[![](../../svg/rx.svg)](http://reactivex.io/documentation/subject.html)

**ReplaySubject** ensures that all observers see the same sequence of emitted items,
even if they subscribe after.

When bufferCapacity argument is 0, then DefaultReplayCapacity is used (currently 16383).
When windowDuration argument is 0, then entries added to the buffer will remain
fresh forever.
