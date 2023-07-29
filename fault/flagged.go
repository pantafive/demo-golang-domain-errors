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

// Blank creates empty Flagged error.
// It is used to pass as a pointer to errors.As.
func Blank() Flagged {
	return fault{}
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
