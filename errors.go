package roller

import "fmt"

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
