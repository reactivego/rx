package DoOnError

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoOnError(t *testing.T) {
	var oerr error
	_, err := ThrowInt(errors.New("error")).DoOnError(func(err error) { oerr = err }).ToSlice()
	assert.Equal(t, err, oerr)
}
