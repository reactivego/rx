package Just

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJust(t *testing.T) {
	result, err := JustInt(1).ToSingle()
	expect := 1
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
