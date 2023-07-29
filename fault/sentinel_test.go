package fault_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrSpecific = errors.New("specific error")

func TestSentinelError_ExampleUsage(t *testing.T) {
	err := errors.New("third party error") //nolint:goerr113

	// add context to the error, mark it as ErrSpecific
	err = errors.Join(ErrSpecific, err)

	err = fmt.Errorf("wrap the error: %w", err)

	specificErrorOccurred := false

	switch {
	case err == nil:
		// do nothing
	case errors.Is(err, ErrSpecific):
		// handle specific error
		specificErrorOccurred = true
	default:
		// generic error handling
	}

	assert.True(t, specificErrorOccurred)
}
