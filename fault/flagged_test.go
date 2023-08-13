package fault_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pantafive/demo-golang-domain-errors/fault"
)

func ExampleFlag() {
	err := fault.New(errors.New("some error"), fault.Alfa) //nolint:goerr113

	if err == nil { //nolint:revive
		// do nothing
	}

	// error.As will "unwrap" the error and assign it to flaggedError
	var flaggedError fault.Flagged

	_ = errors.As(err, &flaggedError)

	// In the example we intentionally ignore flagged.Charlie to demonstrate that
	// exhaustive linter will alert us about it:
	// missing cases in a switch of type flagged.Flag: flagged.Charlie (exhaustive).
	// In this way, we achieve Checked Exceptions effect.
	switch flaggedError.Flag() {
	case fault.Alfa:
		fmt.Print("Error with Alfa flag handled")
	case fault.Bravo:
		// handle Bravo flag
	default:
		// handle generic error
	}

	// Output: Error with Alfa flag handled
}

var (
	errRoot  = errors.New("root error")
	errChild = errors.New("child error")
)

func TestFlagged_Is(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)
	assert.True(t, errors.Is(err, errRoot))

	wrappedErr := fmt.Errorf("wrapped: %w", err)
	assert.True(t, errors.Is(wrappedErr, errRoot))
}

func TestFlagged_As(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	var want fault.Flagged

	assert.True(t, errors.As(err, &want))
	assert.Equal(t, want.Flag(), fault.Alfa)

	wrappedErr := fmt.Errorf("wrapped: %w", err)

	assert.True(t, errors.As(wrappedErr, &want))
	assert.Equal(t, want.Flag(), fault.Alfa)
}

func TestFlagged_Wrap(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	wrappedErr := fmt.Errorf("wrapped: %w", err)

	var want fault.Flagged

	assert.True(t, errors.Is(wrappedErr, errRoot))
	assert.True(t, errors.As(wrappedErr, &want))
}

func TestFlagged_Join(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	joinedErr := errors.Join(err, errChild)

	var want fault.Flagged

	assert.True(t, errors.Is(joinedErr, errRoot))
	assert.True(t, errors.Is(joinedErr, errChild))
	assert.True(t, errors.As(joinedErr, &want))
}

func TestFlagged_FlagOverFlag(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)
	err = fault.New(err, fault.Bravo)

	var want fault.Flagged
	assert.True(t, errors.As(err, &want))
	assert.Equal(t, want.Flag(), fault.Bravo)
}
