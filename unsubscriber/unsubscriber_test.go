package unsubscriber

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestUnsubscriber(t *testing.T) {
	var s unsubscriber
	assert.False(t, s.Unsubscribed())
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
}
