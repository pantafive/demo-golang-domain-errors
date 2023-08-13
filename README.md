# Handle Specific Errors in Go

## Preface

Dave Cheney explains an idiomatic way to work with errors
in [GopherCon 2016: Dave Cheney - Dont Just Check Errors Handle Them Gracefully](https://www.youtube.com/watch?v=lsBF58Q-DnY)
([text extraction of the presentation](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully))
where he introduces his [errors](https://github.com/pkg/errors) library.

The library provides a way to annotate errors with the context without making an error type public.

> If your code implements an interface whose contract requires a specific error type, all implementors of that interface
> need to depend on the package that defines the error type.


Significant changes in Go 1.13 and then Go 1.20 made this library obsolete. Nate Finch rethought Dave Cheney's approach
in [Error Flags](https://npf.io/2021/04/errorflags/). This simple test demonstrates the approach:

```go
package fault_test

import (
	"errors"
	"fmt"
)

var ErrSpecific = errors.New("specific error")

func ExampleSentinelError() {
	err := errors.New("third party error")

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

```

The current approach is sufficient for most cases. However, if we want to ensure that all specific errors are handled
correctly, we need to use a more advanced method.

## "Checked" Errors in Go

[Exhaustive](https://github.com/nishanths/exhaustive) linter checks the exhaustiveness of enum switch statements.
To use it with errors we need to implement `Flagged` interface that embeds `error` and has `Flag()` method that returns
a
flag – enum value that will be used in the switch statement in the following way:

```go
package fault_test

import (
	"errors"
	"fmt"

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

```

Exhaustive linter will report if there are missing cases in the switch statement.
You can run `golangci-lint run ./src` to see it in action:

```
❯ golangci-lint run ./fault
fault/flagged_test.go:29:2: missing cases in switch of type fault.Flag: fault.Charlie (exhaustive)
```

---

The implementation of flagged error is pretty simple:

```go
package fault

type Flag string

const (
	Alfa    Flag = "Alfa"
	Bravo   Flag = "Bravo"
	Charlie Flag = "Charlie"
)

type Flagged interface {
	error
	Flag() Flag
}

// New creates a new flagged error from the existing error and provided flag.
func New(err error, flag Flag) Flagged {
	return fault{error: err, flag: flag}
}

type fault struct {
	error
	flag Flag
}

func (e fault) Unwrap() error {
	return e.error
}

func (e fault) Flag() Flag {
	return e.flag
}

```

That's it. Example usage can be found in [fault/flagged_test.go](fault/flagged_test.go).

---

Other articles on the topic:

- [Effective Error Handling in Golang](https://earthly.dev/blog/golang-errors/) by Earthly
- [Error handling in Go HTTP applications](https://www.joeshaw.org/error-handling-in-go-http-applications/) by Joe Shaw
- [Error handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html) by Upspin
- [Failure is your Domain](https://middlemost.com/failure-is-your-domain/) by Ben Johnson
- [Handling Go errors](https://about.sourcegraph.com/blog/go/gophercon-2019-handling-go-errors) by Marwan Sulaiman
- [My Custom HTTP Error in Golang](https://clavinjune.dev/en/blogs/my-custom-http-error-in-golang/) by Clavin June
