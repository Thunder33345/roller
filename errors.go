package perms_manager

import "fmt"

var _ error = (*MissingGroupsError)(nil) // ensure MissingGroupsError implements error

//MissingGroupsError Is an error raised by Process when group provider fails to load a certain groups
type MissingGroupsError struct {
	groups []string
}

func NewMissingGroupsError(groups []string) MissingGroupsError {
	return MissingGroupsError{
		groups: groups,
	}
}

func (e MissingGroupsError) Error() string {
	return fmt.Sprintf("Missing Group: %v", e.groups)
}

func (e MissingGroupsError) Groups() []string {
	return e.groups
}
