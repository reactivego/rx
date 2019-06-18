package Start

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	source := StartInt(func() (int, error) { return 42, nil })
	result, err := source.ToSingle()
	expect := 42
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
	result, err = source.ToSingle()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

func TestStartWithError(t *testing.T) {
	result, err := StartInt(func() (int, error) { return 0, RxError("start-with-error") }).ToSlice()
	expect := []int(nil)
	assert.EqualError(t, err, "start-with-error")
	assert.Equal(t, expect, result)
}
