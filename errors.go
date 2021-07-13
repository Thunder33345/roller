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