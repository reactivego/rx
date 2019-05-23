package DoOnComplete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoOnComplete(t *testing.T) {
	complete := false
	result, err := EmptyInt().DoOnComplete(func() { complete = true }).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int(nil), result)
	assert.True(t, complete)
}

func TestDoOnCompleteSubscribeNext(t *testing.T) {
	wait := make(chan struct{})
	result := []int{}
	_ = FromInts(1, 2, 3, 4, 5).
		DoOnComplete(func() { close(wait) }).
		SubscribeNext(func(v int) { result = append(result, v) }, SubscribeOn(NewGoroutine()))
	<-wait
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}
