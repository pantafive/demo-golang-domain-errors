package fault_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pantafive/demo-golang-domain-errors/fault"
)

var (
	errRoot  = errors.New("root error")
	errChild = errors.New("child error")
)

func TestFlagged_ExampleUsage(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	if err == nil {
		require.True(t, false) // no errors - unexpected
	}

	flaggedError := fault.Blank()

	_ = errors.As(err, &flaggedError)

	// In the example we intentionally ignore flagged.Charlie to demonstrate that
	// exhaustive linter will alert us about it:
	// missing cases in a switch of type flagged.Flag: flagged.Charlie (exhaustive).
	// In this way, we achieve Checked Exceptions effect.
	switch flaggedError.Flag() {
	case fault.Alfa:
		require.True(t, true) // handle alfa - expected
	case fault.Bravo:
		require.True(t, false) // handle bravo - unexpected
	default:
		require.True(t, false) // handle unknown - unexpected
	}
}

func TestFlagged_Is(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)
	assert.True(t, errors.Is(err, errRoot))

	wrappedErr := fmt.Errorf("wrapped: %w", err)
	assert.True(t, errors.Is(wrappedErr, errRoot))
}

func TestFlagged_As(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	want := fault.Blank()

	assert.True(t, errors.As(err, &want))
	assert.Equal(t, want.Flag(), fault.Alfa)

	wrappedErr := fmt.Errorf("wrapped: %w", err)

	assert.True(t, errors.As(wrappedErr, &want))
	assert.Equal(t, want.Flag(), fault.Alfa)
}

func TestFlagged_Wrap(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	wrappedErr := fmt.Errorf("wrapped: %w", err)

	want := fault.Blank()

	assert.True(t, errors.Is(wrappedErr, errRoot))
	assert.True(t, errors.As(wrappedErr, &want))
}

func TestFlagged_Join(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)

	joinedErr := errors.Join(err, errChild)

	want := fault.Blank()

	assert.True(t, errors.Is(joinedErr, errRoot))
	assert.True(t, errors.Is(joinedErr, errChild))
	assert.True(t, errors.As(joinedErr, &want))
}

func TestFlagged_FlagOverFlag(t *testing.T) {
	err := fault.New(errRoot, fault.Alfa)
	err = fault.New(err, fault.Bravo)

	want := fault.Blank()
	assert.True(t, errors.As(err, &want))
	assert.Equal(t, want.Flag(), fault.Bravo)
}
