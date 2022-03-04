package roller

import (
	"errors"
	"fmt"
)

//IsMissingFlagError checks if an error is an missingFlagError,
//missingFlagError are treated specially ignored
//this is exposed for ease of making a custom processor
func IsMissingFlagError(err error) bool {
	mf := &missingFlagError{}
	r := errors.As(err, &mf)
	return r
}

//NewMissingFlagError creates a new missingFlagError
//used to signal to processor that the flag is missing, but it's not a critical error
//see GroupProvider.Flag
func NewMissingFlagError(group, flag string) error {
	return missingFlagError{
		group: group,
		flag:  flag,
	}
}

type missingFlagError struct {
	group, flag string
}

func (e missingFlagError) Error() string {
	return fmt.Sprintf(`missing error: cant find flag "%s" in "%s"`, e.flag, e.group)
}

//providerGroupError is an internal error used by BasicProcessor
type providerGroupError struct {
	group string
	cause error
}

func (e providerGroupError) Error() string {
	return fmt.Sprintf(`cant retrieve group "%s": %v`, e.group, e.cause)
}
func (e providerGroupError) Unwrap() error {
	return e.cause
}
