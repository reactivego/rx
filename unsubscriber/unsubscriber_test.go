package unsubscriber

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestInt32(t *testing.T) {
	var s Int32
	assert.False(t, s.Unsubscribed())
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
}

type E struct{}

func (e E) Escape() {

}

type U struct{}

func (u U) Unsubscribe() {

}

func (u U) Unsubscribed() bool {
	return false
}

func TestUnsubscriberEscape(t *testing.T) {

	type EscapeArtist interface {
		Escape()
	}

	type Combined interface {
		EscapeArtist
		Unsubscriber
	}

	var e E
	var u U

	var c Combined = struct {
		EscapeArtist
		Unsubscriber
	}{e, u}

	c.Escape()
	c.Unsubscribe()
	c.Unsubscribed()

	if _, ok := c.(Unsubscriber); !ok {
		t.Fail()
	}

	if _, ok := c.(EscapeArtist); !ok {
		t.Fail()
	}

	// Prevent EscapeArtist from being available to users
	var d Unsubscriber = &struct{ Unsubscriber }{c.(Unsubscriber)}

	if _, ok := d.(EscapeArtist); ok {
		t.Fail()
	}

}
