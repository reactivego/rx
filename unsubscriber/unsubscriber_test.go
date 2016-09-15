package unsubscriber

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestStandard(t *testing.T) {
	var s Standard
	assert.False(t, s.Unsubscribed())
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
}
