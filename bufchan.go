package rx

import (
	"errors"
	"sync"
	"time"
)

//jig:template NewBufChan

// BufChan is a fixed capacity buffer non-blocking channel where entries
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
type BufChan struct {
	sync.RWMutex
	recv     *sync.Cond
	channel  []bufChanMessage
	read     int64
	write    int64
	size     int64
	duration time.Duration
}

// NewBufChan returns a non-blocking BufChan with given buffer capacity and time
// window. The window specifies how long items send to the channel will remain
// fresh. After sent values become stale they are no longer returned when the
// channel buffer is iterated with an endpoint.
//
// A bufferCapacity of 0 will result in a channel that cannot send data, but
// that can signal that it has been closed. A windowDuration of 0 will make the
// sent values remain fresh forever.
func NewBufChan(bufferCapacity int, windowDuration time.Duration) *BufChan {
	b := &BufChan{
		channel:  make([]bufChanMessage, bufferCapacity+1),
		size:     int64(bufferCapacity + 1),
		duration: windowDuration,
	}
	b.recv = sync.NewCond(b.RLocker())
	return b
}

// Send will append the value at the end of the channel buffer. If the channel
// buffer is full, the first entry in the channel buffer is overwritten and the
// channel buffer start moved to the second entry. If the channel was closed
// previously by calling Close, this call to Send is ignored.
func (b *BufChan) Send(value interface{}) {
	b.Lock()
	defer b.Unlock()
	if b.channel[b.write%b.size].value != nil {
		return
	}
	var staleAfter time.Time
	if b.duration != 0 {
		staleAfter = time.Now().Add(b.duration)
	}
	b.channel[b.write%b.size] = bufChanMessage{value, staleAfter}
	b.write++
	if b.write-b.read == b.size {
		b.read++
		b.channel[b.write%b.size] = bufChanMessage{}
	}
	b.recv.Broadcast()
}

// Close will mark the channel as closed. Pass either an error value to indicate
// an error, or nil to indicate normal completion. Once the channel has been
// closed, all calls to Send will return immediately without modifying the
// channel buffer.
func (b *BufChan) Close(err error) {
	b.Lock()
	defer b.Unlock()
	if err != nil {
		b.channel[b.write%b.size] = bufChanMessage{value: err}
	} else {
		b.channel[b.write%b.size] = bufChanMessage{value: "completion"}
	}
	b.recv.Broadcast()
}

// NewEndpoint will return a receive endpoint that can be used to receive from
// the channel. A new endpoint may also be created for a closed channel, even
// if it was closed with an error. The buffered content of a closed channel can
// be received normally.
func (b *BufChan) NewEndpoint() *BufEndpoint {
	b.RLock()
	defer b.RUnlock()
	return &BufEndpoint{b, b.read, false}
}

// BufEndpoint is a receive endpoint used for receiving from a channel buffer. A
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
type BufEndpoint struct {
	*BufChan
	cursor   int64
	overflow bool
}

// ErrMissingBackpressure is delivered to an endpoint when the channel overflows
// because the endpoint can't keep-up with the data rate at which the sender
// sends values. Other sibling endpoints that are fast enough won't get this
// error and continue to operate normally.
var ErrMissingBackpressure = errors.New("missing backpressure")

// Recv will return the next value in the channel and true. Or when at the end
// of the channel buffer nil and false. The channel buffer can still be
// iterated after it is finalized by calling Close. If the endpoint could not
// keep up with the sender, then it returns nil and false. Closed will in that
// case report ErrMissingBackpressure.
func (ep *BufEndpoint) Recv() (interface{}, bool) {
	ep.RLock()
	defer ep.RUnlock()
	now := time.Now()
	if ep.cursor < ep.read {
		ep.overflow = true
	}
	if ep.overflow {
		return nil, false
	}
	for ep.cursor != ep.write {
		entry := ep.channel[ep.cursor%ep.size]
		ep.cursor++
		if entry.stale.IsZero() || entry.stale.After(now) {
			return entry.value, true
		}
	}
	return nil, false
}

// Closed returns true once the channel buffer has been finalized by calling
// Close on the channel. The returned error either has a value (the error) or
// is nil (to indicate completion). Note, that this call is independent of how
// far the endpoint has currently iterated the channel buffer.
//
// The error ErrMissingBackpressure will be delivered if the endpoint does not
// receive data fast enough to keep up with the sender.
func (ep *BufEndpoint) Closed() (error, bool) {
	ep.RLock()
	defer ep.RUnlock()
	if ep.overflow {
		return ErrMissingBackpressure, true
	}
	entry := ep.channel[ep.write%ep.size]
	if entry.value != nil {
		if err, ok := entry.value.(error); ok {
			return err, true
		} else {
			return nil, true
		}
	}
	return nil, false
}

// Wait will wait for a Send or Close call on the channel by the sender.
func (ep *BufEndpoint) Wait() {
	ep.RLock()
	defer ep.RUnlock()
	if ep.overflow {
		return
	}
	if ep.cursor != ep.write {
		return
	}
	if ep.channel[ep.write%ep.size].value != nil {
		return
	}
	ep.recv.Wait()
}

// Range will call the function for every received value and will call the
// function one final time when the channel is closed.
func (ep *BufEndpoint) Range(f func(interface{}, error, bool) bool) {
	for more := true; more; {
		if next, ok := ep.Recv(); ok {
			more = f(next, nil, false)
		} else {
			if err, ok := ep.Closed(); ok {
				f(nil, err, true)
				return
			} else {
				ep.Wait()
			}
		}
	}
}

// bufChanMessage instances store the values in a BufChan. It also
// contains a timestamp indicating when it will become stale. You will never
// need to deal with bufChanMessage instances directly.
type bufChanMessage struct {
	value interface{}
	stale time.Time
}
