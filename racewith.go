package rx

import (
	"sync"
)

func (observable Observable[T]) RaceWith(others ...Observable[T]) Observable[T] {
	if len(others) == 0 {
		return observable
	}
	return func(observe Observer[T], scheduler Scheduler, subscriber Subscriber) {
		var race struct {
			sync.Mutex
			subscribers []Subscriber
		}
		race.subscribers = make([]Subscriber, 1+len(others))
		for i := range race.subscribers {
			race.subscribers[i] = subscriber.Add()
		}
		subscribe := func(subscriber Subscriber) (Observer[T], Scheduler, Subscriber) {
			observer := func(next T, err error, done bool) {
				race.Lock()
				defer race.Unlock()
				if subscriber.Subscribed() {
					for i := range race.subscribers {
						if race.subscribers[i] != subscriber {
							race.subscribers[i].Unsubscribe()
						}
					}
					race.subscribers = nil
					observe(next, err, done)
					if done {
						subscriber.Unsubscribe()
					}
				}
			}
			return observer, scheduler, subscriber
		}
		race.Lock()
		defer race.Unlock()
		observable(subscribe(race.subscribers[0]))
		for i, other := range others {
			other(subscribe(race.subscribers[i+1]))
		}
	}
}
