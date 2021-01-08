# Subject

[![](../../../assets/godev.svg?raw=true)](https://pkg.go.dev/github.com/reactivego/rx/test/Subject?tab=doc)
[![](../../../assets/godoc.svg?raw=true)](https://godoc.org/github.com/reactivego/rx/test/Subject)
[![](../../../assets/rx.svg?raw=true)](http://reactivex.io/documentation/subject.html)

**Subject** is a combination of an observer and observable. Subjects are
special because they are the only reactive constructs that support
multicasting. The items sent to it through its observer side are
multicasted to multiple clients subscribed to its observable side.

A Subject embeds both an Observer and an Observable. This exposes the methods
and fields of both types on Subject. Use the Observable methods to subscribe
to it. Use the ObserveFunc Next, Error and Complete methods to feed data to
it.

After a subject has been terminated by calling either Error or Complete, it
goes into terminated state. All subsequent calls to its observer side will be
silently ignored. All subsequent subscriptions to the observable side will be
handled according to the specific behavior of the subject. There are different
types of subjects, see the NewReplaySubject function for more info about a 
subject that is capable of replaying a sequence of emissions. The standard
subject will not buffer any emissions but will just multicast them.
