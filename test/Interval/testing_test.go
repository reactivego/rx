package Interval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type tiStruct struct {
	v int
	e bool // True if elapsed time was >= 10ms
}

func TestInterval(t *testing.T) {
	start := time.Now()
	seen := []tiStruct{}
	Interval(time.Millisecond * 10).
		Take(5).
		SubscribeNext(func(n int) {
			seen = append(seen, tiStruct{n, time.Now().Sub(start) >= 10*time.Millisecond})
		})
	assert.Equal(t, []tiStruct{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
	}, seen)
}
