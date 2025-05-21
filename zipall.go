package rx

import (
	"errors"
	"sync"
)

var ErrZipBufferOverflow = errors.Join(Err, errors.New("zip buffer overflow"))

func ZipAll[T any](observable Observable[Observable[T]], options ...MaxBufferSizeOption) Observable[[]T] {
	var maxBufferSize = 0
	for _, option := range options {
		option(&maxBufferSize)
	}
	return func(observe Observer[[]T], scheduler Scheduler, subscriber Subscriber) {
		var sources []Observable[T]
		var buffers struct {
			sync.Mutex
			values    [][]T  // Buffer for each source
			completed []bool // Track which sources have completed
			done      bool
		}
		makeObserver := func(sourceIndex int) Observer[T] {
			return func(next T, err error, done bool) {
				buffers.Lock()
				defer buffers.Unlock()
				if !buffers.done {
					switch {
					case !done:
						// Check if adding this item would exceed the buffer size limit
						if maxBufferSize > 0 && len(buffers.values[sourceIndex]) >= maxBufferSize {
							buffers.done = true
							var zero []T
							observe(zero, ErrZipBufferOverflow, true)
							return
						}
						buffers.values[sourceIndex] = append(buffers.values[sourceIndex], next)
						result := make([]T, len(buffers.values))
						for _, values := range buffers.values {
							if len(values) == 0 {
								return
							}
						}
						for i := range buffers.values {
							result[i] = buffers.values[i][0]
							buffers.values[i] = buffers.values[i][1:] // Remove the used value
						}
						observe(result, nil, false)
					case err != nil:
						buffers.done = true
						var zero []T
						observe(zero, err, true)
					default:
						buffers.completed[sourceIndex] = true
						for i, completed := range buffers.completed {
							if completed && len(buffers.values[i]) == 0 {
								buffers.done = true
								var zero []T
								observe(zero, nil, true)
								return
							}
						}
					}
				}
			}
		}
		observer := func(next Observable[T], err error, done bool) {
			switch {
			case !done:
				sources = append(sources, next)
			case err != nil:
				var zero []T
				observe(zero, err, true)
			default:
				if len(sources) == 0 {
					var zero []T
					observe(zero, nil, true)
					return
				}
				runner := scheduler.Schedule(func() {
					if subscriber.Subscribed() {
						numSources := len(sources)
						buffers.values = make([][]T, numSources)
						buffers.completed = make([]bool, numSources)
						for sourceIndex, source := range sources {
							if !subscriber.Subscribed() {
								return
							}
							source.AutoUnsubscribe()(makeObserver(sourceIndex), scheduler, subscriber)
						}
					}
				})
				subscriber.OnUnsubscribe(runner.Cancel)
			}
		}
		observable.AutoUnsubscribe()(observer, scheduler, subscriber)
	}
}
