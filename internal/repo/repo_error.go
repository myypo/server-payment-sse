package repo

import (
	"fmt"
	dom "payment-sse/internal/domain"
	verbErr "payment-sse/internal/error/verbose"
	"slices"
)

type RepoError interface {
	verbErr.VerboseError
	Type() RepoErrorType
}

type repoError struct {
	userErr error
	techErr error
	typ     RepoErrorType
}

func (e *repoError) Error() string {
	return e.userErr.Error()
}

func (e *repoError) Verbose() error {
	return e.techErr
}

func (e *repoError) User() error {
	return e.userErr
}

func (e *repoError) Type() RepoErrorType {
	return e.typ
}

type RepoErrorType int

const (
	NotFound RepoErrorType = iota
	Conflict
	BadRequest
	Internal
)

func newErrNotFound(modName string, techErr error) RepoError {
	return &repoError{
		userErr: errNotFound(modName),
		techErr: techErr,
		typ:     NotFound,
	}
}

func newErrConflict(modName string, techErr error) RepoError {
	return &repoError{
		userErr: errConflict(modName),
		techErr: techErr,
		typ:     Conflict,
	}
}

func newErrBadRequest(modName string, techErr error, issue string) RepoError {
	return &repoError{
		userErr: errBadRequest(modName, issue),
		techErr: techErr,
		typ:     BadRequest,
	}
}

func newErrInternal(modName string, techErr error) RepoError {
	return &repoError{
		userErr: errInternal(modName),
		techErr: techErr,
		typ:     Internal,
	}
}

func errBadRequest(modName string, issue string) error {
	return fmt.Errorf("the %s related request is invalid because %s", modName, issue)
}

func errNotFound(modName string) error {
	return fmt.Errorf("the requested %s does not exist", modName)
}

func errConflict(modName string) error {
	return fmt.Errorf("the provided %s data is using duplicate values", modName)
}

func errInternal(modName string) error {
	return fmt.Errorf("unexpected error occured when performing operation on %s", modName)
}

func NewUnexpectedRepoError(
	techErr error,
	modName string,
) RepoError {
	return newErrInternal(modName, techErr)
}

func NewRepoError(
	errTyp RepoErrorType,
	techErr error,
	modName string,
	issueName string,
	expTypes ...RepoErrorType,
) RepoError {
	if !slices.Contains(expTypes, errTyp) {
		return newErrInternal(modName, techErr)
	}

	switch errTyp {
	case Conflict:
		return newErrConflict(modName, techErr)
	case NotFound:
		return newErrNotFound(modName, techErr)
	case BadRequest:
		return newErrBadRequest(modName, techErr, issueName)
	default:
		return newErrInternal(modName, techErr)
	}
}

func DomErrorFromRepo(rerr RepoError) dom.DomError {
	switch rerr.Type() {
	case NotFound:
		return dom.DomErrorFromVerbose(rerr, dom.NotFound)
	case Conflict:
		return dom.DomErrorFromVerbose(rerr, dom.Conflict)
	case BadRequest:
		return dom.DomErrorFromVerbose(rerr, dom.BadRequest)
	case Internal:
		return dom.DomErrorFromVerbose(rerr, dom.Internal)
	}

	return dom.DomErrorFromVerbose(rerr, dom.Internal)
}
