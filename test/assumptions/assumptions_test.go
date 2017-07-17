package assumptions

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestNilRange(t *testing.T) {
	var items []*struct{}
	assert.NotPanics(t, func() {
		for _, v := range items {
			_ = v
		}
	}, "Calling range on nil slice should not panic")
}
