package subscriber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	var s subscription
	assert.False(t, s.Closed())
	s.Unsubscribe()
	assert.True(t, s.Closed())
}
