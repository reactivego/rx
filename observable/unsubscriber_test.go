package observable

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestUnsubscriberInt(t *testing.T) {
	var s UnsubscriberInt
	assert.False(t, s.Unsubscribed())
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
}

func TestUnsubscriberChannel(t *testing.T) {
	done := make(chan bool)
	unsubscribed := false
	var s Unsubscriber = NewUnsubscriberChannel()
	events, ok := s.(UnsubscriberEvents)
	assert.True(t, ok)
	events.OnUnsubscribe(func() { unsubscribed = true; done <- true })
	assert.False(t, s.Unsubscribed())
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
	<-done
	assert.True(t, unsubscribed)
}

func TestUnsubscriberFunc(t *testing.T) {
	callPerformed := false
	makeUnsub := func() Unsubscriber {
		unsub := UnsubscriberFunc(func() {
			callPerformed = true
		})
		return &unsub
	}
	unsub := makeUnsub()

	unsub.Unsubscribe()
	assert.True(t, callPerformed)
	assert.True(t, unsub.Unsubscribed())
	callPerformed = false
	assert.NotPanics(t, func() {
		unsub.Unsubscribe()
	}, "Calling unsub.Unsubscribe() again should not panic")
	assert.True(t, unsub.Unsubscribed())
	// Should not have changed callPerformed again
	assert.False(t, callPerformed)
}

// func TestLinkedUnsubscriber(t *testing.T) {
// 	linked := NewLinkedUnsubscriber()
// 	var sub UnsubscriberInt
// 	assert.False(t, linked.Unsubscribed())
// 	assert.False(t, sub.Unsubscribed())
// 	linked.Link(sub)
// 	assert.Panics(t, func() { linked.Link(sub) })
// 	linked.Unsubscribe()
// 	assert.True(t, sub.Unsubscribed())
// 	assert.True(t, linked.Unsubscribed())
// }

// func TestLinkedUnsubscriberUnsubscribesTargetOnLink(t *testing.T) {
// 	linked := NewLinkedUnsubscriber()
// 	sub := NewGenericUnsubscriber()
// 	linked.Unsubscribe()
// 	assert.True(t, linked.Unsubscribed())
// 	assert.False(t, sub.Unsubscribed())
// 	linked.Link(sub)
// 	assert.True(t, linked.Unsubscribed())
// 	assert.True(t, sub.Unsubscribed())
// }
