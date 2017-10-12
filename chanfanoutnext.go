package rx

import (
	"sync"
)

//jig:template NewChanFanOutNext<Foo>
//jig:needs NewChanNext<Foo>

// ChanFanOutNextFoo is a blocking fan-out implementation writing to multiple
// ChanNextFoo channels that implement message buffering. Use a call to
// NewChannel to add a receiving channel. A newly created channel will start
// receiving any messages that are send subsequently into the fan-out channel.
//
// This channel works by fanning out the Send calls to a slice of ChanNextFoo
// wrapped Go channels. When one of the channels is slow and its buffer size
// reaches its maximum capacity, then the Send into that channel will block.
// This provides so called blocking backpressure to the sender by blocking its
// goroutine. Effectively the slowest channel will dictate the throughput for
// the whole fan-out assembly.
type ChanFanOutNextFoo struct {
	sync.Mutex
	capacity int
	channel  []*ChanNextFoo
	closed   bool
	err      error
}

// NewChanFanOutNextFoo creates a new ChanFanOutNextFoo with the given buffer
// capacity to use for the channels in the fanout
func NewChanFanOutNextFoo(bufferCapacity int) *ChanFanOutNextFoo {
	return &ChanFanOutNextFoo{capacity: bufferCapacity}
}

// NewChannel adds a new ChanNextFoo channel to the fan-out and returns its
// Channel field (<- chan NextFoo) along with a cancel function. It is critical
// that the cancel function is called to indicate you want to stop receiving
// data, as simply abandoning the channel will fill up its buffer and then block
// the whole fan-out assembly from further processing messages.
//
// When the channel has been closed by the sender by calling Close, then the
// cancel function does not have to be called (but doing so does not hurt).
// But never call the cancel function more than once, because that will panic on
// closing an internal cancel channel twice.
//
// It is perfectly fine to call NewChannel on a fan-out channel that was already
// closed. The channel returned will replay the error if present and is then
// closed immediately.
func (m *ChanFanOutNextFoo) NewChannel() (<-chan NextFoo, func()) {
	m.Lock()
	defer m.Unlock()
	if m.closed {
		channel := make(chan NextFoo, 1)
		if m.err != nil {
			channel <- NextFoo{Err: m.err}
		}
		close(channel)
		return channel, func() {}
	}
	ch := NewChanNextFoo(m.capacity)
	m.channel = append(m.channel, ch)
	return ch.Channel, func() {
		ch.Cancel()
		m.Lock()
		defer m.Unlock()
		for i, c := range m.channel {
			if c == ch {
				m.channel = append(m.channel[:i], m.channel[i+1:]...)
				return
			}
		}
	}
}

// Send is used to multicast a value to multiple receiving channels.
func (m *ChanFanOutNextFoo) Send(value foo) {
	m.Lock()
	for _, c := range m.channel {
		c.Send(value)
	}
	m.Unlock()
}

// Close is used to close all receiving channels in the fan-out assembly. If
// err is not nil, then the error is send to the channels first before they
// are closed.
func (m *ChanFanOutNextFoo) Close(err error) {
	m.Lock()
	for _, c := range m.channel {
		c.Close(err)
	}
	m.channel = nil
	m.closed = true
	m.err = err
	m.Unlock()
}
