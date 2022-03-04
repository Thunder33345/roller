package roller

import (
	"errors"
	"fmt"
)

//IsErrorNotExist checks if a given error is NotFoundError
func IsErrorNotExist(err error) bool {
	return errors.Is(err, notFoundError{})
}

//NewNotFoundError creates a new NotFoundError with given string as error text
func NewNotFoundError(str string) error {
	return notFoundError{err: str}
}

type notFoundError struct {
	err string
}

func (e notFoundError) Error() string {
	return e.err
}

func (e notFoundError) Is(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}

var _ error = (*MissingGroupError)(nil) // ensure MissingGroupError implements error

//MissingGroupError Is an error raised by Process when group provider fails to load a certain groups
type MissingGroupError struct {
	group string
	error error
}

func NewMissingGroupsError(gid string, err error) MissingGroupError {
	return MissingGroupError{
		group: gid,
		error: err,
	}
}

func (e MissingGroupError) Error() string {
	return fmt.Sprintf("failed to access group \"%v\": %v", e.Group(), e.error)
}

func (e MissingGroupError) Unwrap() error {
	return e.error
}

func (e MissingGroupError) Group() string {
	return e.group
}
