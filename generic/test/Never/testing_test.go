package Never

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNever(t *testing.T) {
	never := true
	subscription := NeverInt().SubscribeNext(func(next int) {
		never = false
	}, SubscribeOn(NewGoroutineScheduler()))

	go func() {
		time.Sleep(100 * time.Millisecond)
		subscription.Unsubscribe()
	}()

	assert.False(t, subscription.Closed())
	subscription.Wait()
	assert.True(t, subscription.Closed())
	assert.True(t, never)
}
