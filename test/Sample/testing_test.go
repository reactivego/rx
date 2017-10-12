package Sample

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	const _90ms = 90 * time.Millisecond
	const _200ms = 200 * time.Millisecond
	result, err := Interval(_90ms).Sample(_200ms).Take(3).ToSlice()
	expect := []int{1, 3, 5}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
