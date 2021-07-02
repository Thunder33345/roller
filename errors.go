package ranker

import "fmt"

var _ error = (*MissingGroupsError)(nil) // ensure MissingGroupsError implements error

type MissingGroupsError struct {
	groups []string
}

func (e MissingGroupsError) Error() string {
	return fmt.Sprintf("Missing Group: %v", e.groups)
}

func (e MissingGroupsError) Groups() []string {
	return e.groups
}

var _ error = (*ErrorMissingRawPermissionList)(nil) // ensure ErrorMissingRawPermissionList implements error

type ErrorMissingRawPermissionList struct {
	uid string
}

func (e ErrorMissingRawPermissionList) Error() string {
	return fmt.Sprintf("Missing PermissionList: %v", e.uid)
}

func (e ErrorMissingRawPermissionList) UID() string {
	return e.uid
}
