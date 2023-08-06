package fault_test

import (
	"errors"
	"fmt"
)

var ErrSpecific = errors.New("specific error")

func ExampleSentinelError() { //nolint:govet
	err := errors.New("third party error") //nolint:goerr113

	// add context to the error, mark it as ErrSpecific
	err = errors.Join(ErrSpecific, err)

	err = fmt.Errorf("wrap the error: %w", err)

	switch {
	case err == nil:
		// do nothing
	case errors.Is(err, ErrSpecific):
		fmt.Print("Specific error occurred")
	default:
		// generic error handling
	}

	// Output: Specific error occurred
}
