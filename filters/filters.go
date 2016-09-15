package filters

import (
	"errors"
	"sync"
	"time"
)

type ObserverFunc func(interface{}, error, bool)

func (f ObserverFunc) Next(next interface{}) {
	f(next, nil, false)
}

func (f ObserverFunc) Error(err error) {
	f(nil, err, err == nil)
}

func (f ObserverFunc) Complete() {
	f(nil, nil, true)
}

type Filter func(ObserverFunc) ObserverFunc

func Distinct() Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		seen := map[interface{}]struct{}{}
		operator := func(next interface{}, err error, completed bool) {
			if err == nil && !completed {
				if _, ok := seen[next]; ok {
					return
				}
				seen[next] = struct{}{}
			}
			observer(next, err, completed)
		}
		return operator
	}
	return filter
}

func ElementAt(n int) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		i := 0
		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed || i == n {
				observer(next, err, completed)
			}
			i++
		}
		return operator
	}
	return filter
}

func Where(f func(next interface{}) bool) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed || f(next) {
				observer(next, err, completed)
			}
		}
		return operator
	}
	return filter
}

func First() Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		start := true
		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed || start {
				observer(next, err, completed)
			}
			start = false
		}
		return operator
	}
	return filter
}

func Last() Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		have := false
		var last interface{}
		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed {
				if have {
					observer(last, nil, false)
				}
				observer(nil, err, completed)
			}
			last = next
			have = true
		}
		return operator
	}
	return filter
}

func Skip(n int) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		i := 0
		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed || i >= n {
				observer(next, err, completed)
			}
			i++
		}
		return operator
	}
	return filter
}

func SkipLast(n int) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		read := 0
		write := 0
		n++
		buffer := make([]interface{}, n)
		operator := func(next interface{}, err error, completed bool) {
			switch {
			case err != nil:
				observer.Error(err)
			case completed:
				observer.Complete()
			default:
				buffer[write] = next
				write = (write + 1) % n
				if write == read {
					observer.Next(buffer[read])
					read = (read + 1) % n
				}
			}
		}
		return operator
	}
	return filter
}

func Take(n int) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		taken := 0
		operator := func(next interface{}, err error, completed bool) {
			if taken < n {
				observer(next, err, completed)
				if err == nil && !completed {
					taken++
					if taken >= n {
						observer.Complete()
					}
				}
			}
		}
		return operator
	}
	return filter
}

func TakeLast(n int) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		read := 0
		write := 0
		n++
		buffer := make([]interface{}, n)
		operator := func(next interface{}, err error, completed bool) {
			switch {
			case err != nil:
				for read != write {
					observer.Next(buffer[read])
					read = (read + 1) % n
				}
				observer.Error(err)
			case completed:
				for read != write {
					observer.Next(buffer[read])
					read = (read + 1) % n
				}
				observer.Complete()
			default:
				buffer[write] = next
				write = (write + 1) % n
				if write == read {
					read = (read + 1) % n
				}
			}
		}
		return operator
	}
	return filter
}

func IgnoreElements() Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed {
				observer(next, err, completed)
			}
		}
		return operator
	}
	return filter
}

func IgnoreCompletion() Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		operator := func(next interface{}, err error, completed bool) {
			if !completed {
				observer(next, err, completed)
			}
		}
		return operator
	}
	return filter
}

func One() Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		count := 0
		var value interface{}
		operator := func(next interface{}, err error, completed bool) {
			if count < 2 {
				switch {
				case err != nil:
					observer.Error(err)
				case completed:
					if count == 1 {
						observer.Next(value)
						observer.Complete()
					} else {
						observer.Error(errors.New("expected one value, got none"))
					}
				default:
					count++
					if count == 1 {
						value = next
					} else {
						observer.Error(errors.New("expected one value, got multiple"))
					}
				}
			}
		}
		return operator
	}
	return filter
}

type timedEntry struct {
	v interface{}
	t time.Time
}

func Replay(size int, duration time.Duration) Filter {
	read := 0
	write := 0
	if size == 0 {
		size = 16384
	}
	if duration == 0 {
		duration = time.Hour * 24 * 7 * 52
	}
	size++
	buffer := make([]timedEntry, size)

	filter := func(observer ObserverFunc) ObserverFunc {
		operator := func(next interface{}, err error, completed bool) {
			now := time.Now()
			switch {
			case err != nil:
				cursor := read
				for cursor != write {
					if buffer[cursor].t.After(now) {
						observer.Next(buffer[cursor].v)
					}
					cursor = (cursor + 1) % size
				}
				observer.Error(err)
			case completed:
				cursor := read
				for cursor != write {
					if buffer[cursor].t.After(now) {
						observer.Next(buffer[cursor].v)
					}
					cursor = (cursor + 1) % size
				}
				observer.Complete()
			default:
				buffer[write] = timedEntry{next, time.Now().Add(duration)}
				write = (write + 1) % size
				if write == read {
					if buffer[read].t.After(now) {
						observer.Next(buffer[read].v)
					}
					read = (read + 1) % size
				}
			}
		}
		return operator
	}
	return filter
}

func Sample(window time.Duration) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		var cancel = make(chan struct{})
		var last struct {
			sync.Mutex
			Value interface{}
			Fresh bool
		}

		// TODO: make sure sampler gets killed on Unsubcribe
		sampler := func() {
			for {
				select {
				case <-time.After(window):
					last.Lock()
					if last.Fresh {
						observer.Next(last.Value)
						last.Fresh = false
					}
					last.Unlock()
				case <-cancel:
					return
				}
			}
		}
		go sampler()

		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed {
				close(cancel)
				if err != nil {
					observer.Error(err)
				} else {
					observer.Complete()
				}
			} else {
				last.Lock()
				defer last.Unlock()
				last.Value = next
				last.Fresh = true
			}
		}
		return operator
	}
	return filter
}

func Debounce(duration time.Duration) Filter {
	filter := func(observer ObserverFunc) ObserverFunc {
		errch := make(chan error)
		valuech := make(chan interface{})

		// TODO: make sure debouncer gets killed on Unsubscribe
		debouncer := func() {
			var nextValue interface{}
			var timeout <-chan time.Time
			for {
				select {
				case nextValue = <-valuech:
					timeout = time.After(duration)
				case <-timeout:
					observer.Next(nextValue)
					timeout = nil
				case err := <-errch:
					if timeout != nil {
						observer.Next(nextValue)
					}
					if err != nil {
						observer.Error(err)
					} else {
						observer.Complete()
					}
					return
				}
			}
		}
		go debouncer()

		operator := func(next interface{}, err error, completed bool) {
			if err != nil || completed {
				errch <- err
			} else {
				valuech <- next
			}
		}
		return operator
	}
	return filter
}
