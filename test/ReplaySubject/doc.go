/*
	ReplaySubject	http://reactivex.io/documentation/subject.html

NewReplaySubject creates a new ReplaySubject. ReplaySubject ensures that
all observers see the same sequence of emitted items, even if they
subscribe after. When bufferCapacity argument is 0, then MaxReplayCapacity is
used (currently 16383). When windowDuration argument is 0, then entries added
to the buffer will remain fresh forever.
*/
package ReplaySubject

import _ "github.com/reactivego/rx"
