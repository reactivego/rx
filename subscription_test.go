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


func TestSubscriberLoop(t *testing.T) {
	parent := &subscription{}

	child1 := parent.Add(parent.Unsubscribe)
	child2 := parent.Add(parent.Unsubscribe)
	child3 := parent.Add(parent.Unsubscribe)

	assert.False(t, parent.Canceled())
	assert.False(t, child1.Canceled())
	assert.False(t, child2.Canceled())
	assert.False(t, child3.Canceled())

	child2.Unsubscribe()

	assert.True(t, parent.Canceled())
	assert.True(t, child1.Canceled())
	assert.True(t, child2.Canceled())
	assert.True(t, child3.Canceled())
}