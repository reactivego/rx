package test

import (
	"math"
	"sync/atomic"
)

const MaxCapacity uint64 = 4096

type Channel struct {
	read     uint64
	write    uint64
	cap      uint64
	mod      uint64
	data     [MaxCapacity]interface{}
	endpoint []*Endpoint
}

type Endpoint struct {
	*Channel
	cursor uint64
}

func NewChannel(capacity int) *Channel {
	// Round capacity up to power of 2
	cap := uint64(1) << uint(math.Ceil(math.Log2(float64(capacity))))
	return &Channel{cap: cap, mod: cap - 1}
}

func (c *Channel) Send(value interface{}) {
	if c.write < c.read+c.cap {
		c.data[c.write&c.mod] = value
	}
}

func (c *Channel) Close(err error) {
	if c.write < c.read+c.cap {
	}
}

func (c *Channel) Full() bool {
	return c.write == c.read+c.cap
}

func (c *Channel) NewEndpoint() *Endpoint {
	endpoint := &Endpoint{c, c.read}
	c.endpoint = append(c.cursor, endpoint)
	return endpoints
}

func (e *Endpoint) Range(f func(interface{}, error, bool) bool) {
	for ; e.cursor < e.write; e.cursor += 1 {
		if !f(e.data[e.cursor], nil, false) {
			return
		}
	}
	f(nil, nil, true)
}
