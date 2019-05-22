package Start

import (
	"errors"
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
	result, err := StartInt(func() (int, error) { return 0, errors.New("error") }).ToSlice()
	expect := []int(nil)
	assert.EqualError(t, err, "error")
	assert.Equal(t, expect, result)
}
