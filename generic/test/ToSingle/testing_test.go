package ToSingle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSingle_errMultipleValue(t *testing.T) {
	_, err := FromInts(1, 2).ToSingle()
	assert.Error(t, err)
}

func TestToSingle_errNoValue(t *testing.T) {
	_, err := EmptyInt().ToSingle()
	assert.Error(t, err)
}

func TestToSingle_correct(t *testing.T) {
	value, err := FromInts(3).ToSingle()
	assert.NoError(t, err)
	assert.Equal(t, 3, value)
}
