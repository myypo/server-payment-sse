package dom

import (
	"errors"
	verbErr "payment-sse/internal/error/verbose"
)

var (
	errBadRequest = errors.New("Bad request input. Please, provide valid data")

	errUnauthorized = errors.New("Access denied. Please, authorize")

	errConflict = errors.New("Please, try providing alternative input")

	errNotFound = errors.New(
		"Sorry, but it seems like what you have requested does not exist. Please, try providing different input",
	)

	errInternal = errors.New(
		"We are sorry, but something unexpected happened. Please, try again later",
	)

	errForbidden = errors.New("Sorry, but the action you have tried to take is forbidden")
)

type DomErrorType int

const (
	NotFound DomErrorType = iota
	Conflict
	Internal
	BadRequest
	Forbidden
	Gone
)

type DomError interface {
	verbErr.VerboseError
	Type() DomErrorType
}

type domError struct {
	userErr error
	techErr error
	typ     DomErrorType
}

func (e *domError) Verbose() error {
	return e.techErr
}

func (e *domError) User() error {
	return e.userErr
}

func (e *domError) Error() string {
	return e.userErr.Error()
}

func (e *domError) Type() DomErrorType {
	return e.typ
}

func NewDomError(userErr, techErr error, typ DomErrorType) DomError {
	return &domError{userErr, techErr, typ}
}

func NewBadRequest(err error) DomError {
	return &domError{userErr: err, techErr: err, typ: BadRequest}
}

func DomErrorFromVerbose(verr verbErr.VerboseError, typ DomErrorType) DomError {
	userErr := verr.User()
	techErr := verr.Verbose()

	return &domError{userErr, techErr, typ}
}
