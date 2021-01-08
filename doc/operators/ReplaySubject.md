# ReplaySubject

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/ReplaySubject?tab=doc)
[![](../../../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx/test/ReplaySubject)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/subject.html)

**ReplaySubject** ensures that all observers see the same sequence of emitted items,
even if they subscribe after.

When bufferCapacity argument is 0, then DefaultReplayCapacity is used (currently 16383).
When windowDuration argument is 0, then entries added to the buffer will remain
fresh forever.
