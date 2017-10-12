package rx

//jig:template NewChanNext<Foo>
//jig:needs Next<Foo>

// ChanNextFoo is a "chan NextFoo" with two additional capabilities. Firstly, it
// can be properly canceled from the receiver side by calling the Cancel method.
// And secondly, it can deliver an error in-band (as opposed to out-of-band) to
// the receiver because of how NextFoo is defined.
type ChanNextFoo struct {
	// Channel can be used directly e.g. in a range statement to read NextFoo
	// items from the channel.
	Channel chan NextFoo

	cancel chan struct{}
	closed bool
}

// NewChanNextFoo creates a new ChanNextFoo with given buffer capacity and
// returns a pointer to it. The Send and Close methods are supposed to be used
// from the sending side by a single goroutine. The Channel field and the Cancel
// method are supposed to be used from the receiving side and may be used from
// different goroutines. Multiple goroutines reading from a single channel will
// fight for values though, because a channel does not multicast.
func NewChanNextFoo(capacity int) *ChanNextFoo {
	return &ChanNextFoo{
		Channel: make(chan NextFoo, capacity),
		cancel:  make(chan struct{}),
	}
}

// Send is used from the sender side to send the next value to the channel. This
// call only returns after delivering the value. If the channel has been closed,
// the value is ignored and the call returns immediately.
func (c *ChanNextFoo) Send(value foo) bool {
	if c.closed {
		return false
	}
	select {
	case <-c.cancel:
		c.closed = true
		close(c.Channel)
		return false
	case c.Channel <- NextFoo{Next: value}:
		return true
	}
}

// Close is used from the sender side to deliver an error value before closing
// the channel. Pass nil to indicate a normal close. If the channel has already
// been closed then the call will return immediately.
func (c *ChanNextFoo) Close(err error) bool {
	if c.closed {
		return false
	}
	c.closed = true
	if err != nil {
		select {
		case <-c.cancel:
			return false
		case c.Channel <- NextFoo{Err: err}:
		}
	}
	close(c.Channel)
	return true
}

// Cancel can be called exactly once from the receiver side to indicate it no
// longer is intereseted in the data or completion status. This cancelation will
// be signaled to the sender. The sender will be correctly aborted if it is
// already blocked on a Send or Close call to deliver a value or error to the
// receiver.
func (c *ChanNextFoo) Cancel() {
	close(c.cancel)
}
