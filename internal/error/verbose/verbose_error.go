package verbErr

import "fmt"

type VerboseError interface {
	error
	Verbose() error
	User() error
}

type verboseError struct {
	techErr error
	userErr error
}

func DefaultVerboseError(techErr error) VerboseError {
	return &verboseError{
		techErr: techErr,
		userErr: fmt.Errorf("unexpected server error"),
	}
}

func (e *verboseError) Error() string {
	return e.userErr.Error()
}

func (e *verboseError) User() error {
	return e.userErr
}

func (e *verboseError) Verbose() error {
	return e.techErr
}
