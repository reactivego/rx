package rx

import (
	"errors"
	"math"
	"sync"
	"time"
)

//jig:template ErrBufferOverflow

// ErrMissingBackpressure is delivered to an endpoint when the channel overflows
// because the endpoint can't keep-up with the data rate at which the sender
// sends values. Other sibling endpoints that are fast enough won't get this
// error and continue to operate normally.
var ErrBufferOverflow = errors.New("buffer overflow")

//jig:template NewReplayChan<Foo>
//jig:needs  ErrBufferOverflow

// ReplayChanFoo is a fixed capacity buffer non-blocking channel where entries
// are appended in a circular fashion. Use Send to append entries to the channel
// buffer, use NewEndpoint to create an endpoint to receive from the channel.
// If the channel buffer is full (contains bufferCapacity items) then the next
// call to Send will overwrite the first entry in the channel buffer.
//
// The actual channel buffer capacity is 1 entry larger than the bufferCapacity
// passed to NewBufChan. A full channel has the write position for the
// next entry immediately adjoining the read postion of the first entry in the
// channel buffer. This means that even in a full channel buffer, there is a
// single emtpy slot at the write position. This single empty slot is used to
// store a "tombstone" when the channel is closed and thus will not be updated
// further. To close the channel, call the Close method with nil or error.
type ReplayChanFoo struct {
	sync.RWMutex
	send            *sync.Cond
	recv            *sync.Cond
	channel         []replayMessageFoo
	tombstone       interface{}
	read            int64
	write           int64
	size            int64
	duration        time.Duration
	endpoints       []*ReplayEndpointFoo
	overflowhandled func(*ReplayChanFoo) bool
}

// BackpressureBlockFoo strategy for handling overflow will block the calling
// goroutine on a condition until an endpoint is connected and starts draining
// the buffer or Close is called to cancel the active Send.
func BackpressureBlockFoo(b *ReplayChanFoo) bool {
	firstfresh := func() int64 {
		now := time.Now()
		for i := b.read; i < b.write; i++ {
			entry := b.channel[i]
			stale := !entry.stale.IsZero() && entry.stale.Before(now)
			if !stale {
				return i
			}
		}
		return b.write
	}
	for b.tombstone == nil {
		fresh := firstfresh()
		if fresh > b.read {
			for _, ep := range b.endpoints {
				if fresh > ep.cursor {
					ep.cursor = fresh
				}
			}
			b.read = fresh
			return true
		}
		if len(b.endpoints) != 0 {
			cursormin := int64(math.MaxInt64)
			for _, ep := range b.endpoints {
				if ep.cursor < cursormin {
					cursormin = ep.cursor
				}
			}
			if b.read < cursormin {
				b.read++
				return true
			}
		}
		b.send.Wait()
	}
	return false
}

// BackpressureBufferFoo strategy for handling overflow will manage a fixed size
// buffer, tracking the endpoint cursors and advancing the read index as little
// as possible. Individual endpoints may overflow and be terminated because they
// can't keep up with the data, while the buffer as a whole continues to operate.
//
// This strategy generates overflow errors in endpoints when they can't keep up
// and for the whole channel when all endpoints overflowed.
func BackpressureErrorFoo(b *ReplayChanFoo) bool {
	firstfresh := func() int64 {
		now := time.Now()
		for i := b.read; i < b.write; i++ {
			entry := b.channel[i]
			stale := !entry.stale.IsZero() && entry.stale.Before(now)
			if !stale {
				return i
			}
		}
		return b.write
	}
	fresh := firstfresh()
	if fresh > b.read {
		for _, ep := range b.endpoints {
			if fresh > ep.cursor {
				ep.cursor = fresh
			}
		}
		b.read = fresh
		return true
	}
	if len(b.endpoints) == 0 {
		b.tombstone = ErrBufferOverflow
	} else {
		cursormax := int64(0)
		for _, ep := range b.endpoints {
			if ep.cursor > cursormax {
				cursormax = ep.cursor
			}
			if b.read > ep.cursor {
				ep.overflow = ErrBufferOverflow
			}
		}
		if b.read > cursormax {
			b.tombstone = ErrBufferOverflow
		} else {
			b.read++
			return true
		}
	}
	return false
}

// BackpressureLatestFoo strategy for handling overflow will will keep the
// latest bufferCapacity number of items in the buffer, dropping the oldest
// ones. This never generates an error, even when holes appear in the data
// received by the endpoints because they can't keep up with the source.
func BackpressureLatestFoo(b *ReplayChanFoo) bool {
	b.read++
	return true
}

// NewBufChan returns a non-blocking ReplayChanFoo with given buffer capacity and time
// window. The window specifies how long items send to the channel will remain
// fresh. After sent values become stale they are no longer returned when the
// channel buffer is iterated with an endpoint.
//
// A bufferCapacity of 0 will result in a channel that cannot send data, but
// that can signal that it has been closed. A windowDuration of 0 will make the
// sent values remain fresh forever.
func NewReplayChanFoo(bufferCapacity int, windowDuration time.Duration) *ReplayChanFoo {
	b := &ReplayChanFoo{
		channel:         make([]replayMessageFoo, bufferCapacity+1),
		size:            int64(bufferCapacity + 1),
		duration:        windowDuration,
		overflowhandled: BackpressureBlockFoo,
	}
	b.recv = sync.NewCond(b.RLocker())
	b.send = sync.NewCond(b)
	return b
}

// Send will append the value at the end of the channel buffer. If the channel
// buffer is full, the first entry in the channel buffer is overwritten and the
// channel buffer start moved to the second entry. If the channel was closed
// previously by calling Close, this call to Send is ignored.
func (b *ReplayChanFoo) Send(value foo) bool {
	if b.tombstone != nil {
		return false
	}

	b.Lock()

	var staleAfter time.Time
	if b.duration != 0 {
		staleAfter = time.Now().Add(b.duration)
	}

	if b.write < b.read+b.size-1 {
		b.channel[b.write%b.size] = replayMessageFoo{value, staleAfter}
		b.write++
	} else if b.overflowhandled(b) {
		b.channel[b.write%b.size] = replayMessageFoo{value, staleAfter}
		b.write++
	}

	b.Unlock()
	b.recv.Broadcast()
	return b.tombstone == nil
}

// Close will mark the channel as closed. Pass either an error value to indicate
// an error, or nil to indicate normal completion. Once the channel has been
// closed, all calls to Send will return immediately without modifying the
// channel buffer. A Send blocked on backpressure blocking will be canceled when
// Close is called.
func (b *ReplayChanFoo) Close(err error) {
	b.Lock()
	if err != nil {
		b.tombstone = err
	} else {
		b.tombstone = "closed"
	}
	b.endpoints = nil
	b.Unlock()
	b.send.Signal()
	b.recv.Broadcast()
}

// NewEndpoint will return a receive endpoint that can be used to receive from
// the channel. A new endpoint may also be created for a closed channel, even
// if it was closed with an error. The buffered content of a closed channel can
// be received normally.
func (b *ReplayChanFoo) NewEndpoint() *ReplayEndpointFoo {
	b.Lock()
	ep := &ReplayEndpointFoo{b, b.read, nil}
	if b.tombstone == nil {
		b.endpoints = append(b.endpoints, ep)
	}
	b.Unlock()
	return ep
}

func (b *ReplayChanFoo) RemoveEndpoint(ep *ReplayEndpointFoo) {
	b.Lock()
	for i, e := range b.endpoints {
		if e == ep {
			b.endpoints = append(b.endpoints[:i], b.endpoints[i+1:]...)
			b.Unlock()
			return
		}
	}
	b.Unlock()
}

// ReplayEndpointFoo is a receive endpoint used for receiving from a channel buffer. A
// newly created endpoint will start reading at the start of the channel buffer
// at the moment NewEndpoint was called. Reading from the endpoint using Recv
// calls will continue until the end of the channel buffer is reached.
//
// The channel buffer may grow while it is being iterated, the endpoint will
// reflect that. The channel may grow too fast for an endpoint to be able to
// keep up. This causes the end of the channel buffer to overflow the endpoint
// current position. At that point the endpoint will have effectively dropped
// all data. If that happens, the endpoint will fail and emit an
// ErrMissingBackpressure error. Sibling endpoints are not affected by this, nor
// is the channel itself.
type ReplayEndpointFoo struct {
	*ReplayChanFoo
	cursor   int64
	overflow error
}

// Recv will return the next value in the channel and true. Or when at the end
// of the channel buffer nil and false. The channel buffer can still be
// iterated after it is finalized by calling Close. If the endpoint could not
// keep up with the sender, then it returns nil and false. Closed will in that
// case report ErrMissingBackpressure.
func (ep *ReplayEndpointFoo) Recv() (foo, bool) {
	ep.RLock()
	var zeroFoo foo
	if ep.overflow != nil {
		ep.RUnlock()
		return zeroFoo, false
	}
	now := time.Now()
	if ep.cursor < ep.read {
		ep.cursor = ep.read // lost items in hole between ep.cursor and ep.read
	}
	for ep.cursor != ep.write {
		entry := ep.channel[ep.cursor%ep.size]
		ep.cursor++
		if entry.stale.IsZero() || entry.stale.After(now) {
			ep.RUnlock()
			ep.send.Signal()
			return entry.value, true
		}
	}
	ep.RUnlock()
	return zeroFoo, false
}

// Closed returns true once the channel buffer has been finalized by calling
// Close on the channel. The returned error either has a value (the error) or
// is nil (to indicate completion). Note, that this call is independent of how
// far the endpoint has currently iterated the channel buffer.
//
// The error ErrMissingBackpressure will be delivered if the endpoint does not
// receive data fast enough to keep up with the sender.
func (ep *ReplayEndpointFoo) Closed() (error, bool) {
	ep.RLock()
	if ep.overflow != nil {
		ep.RUnlock()
		return ep.overflow, true
	}
	if ep.tombstone == nil {
		ep.RUnlock()
		return nil, false
	}
	if err, ok := ep.tombstone.(error); ok {
		ep.RUnlock()
		return err, true
	}
	ep.RUnlock()
	return nil, true
}

// Wait will wait for a Send or Close call on the channel by the sender.
func (ep *ReplayEndpointFoo) Wait() {
	ep.RLock()
	if ep.overflow != nil {
		ep.RUnlock()
		return
	}
	if ep.cursor != ep.write {
		ep.RUnlock()
		return
	}
	if ep.tombstone != nil {
		ep.RUnlock()
		return
	}
	ep.recv.Wait()
	ep.RUnlock()
}

// Range will call the function for every received value and will call the
// function one final time when the channel is closed.
func (ep *ReplayEndpointFoo) Range(f func(foo, error, bool) bool) {
	var zeroFoo foo
	for more := true; more; {
		if next, ok := ep.Recv(); ok {
			more = f(next, nil, false)
		} else {
			if err, ok := ep.Closed(); ok {
				f(zeroFoo, err, true)
				return
			} else {
				ep.Wait()
			}
		}
	}
}

// replayMessageFoo instances store the values in a ReplayChanFoo. It also
// contains a timestamp indicating when it will become stale. You will never
// need to deal with replayMessageFoo instances directly.
type replayMessageFoo struct {
	value foo
	stale time.Time
}
