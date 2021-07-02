package ranker

import "fmt"

var _ error = (*MissingGroups)(nil) // ensure MissingGroups implements error

type MissingGroups struct {
	groups []UID
}

func (e MissingGroups) Error() string {
	return fmt.Sprintf("Missing Group: %v", e.groups)
}

func (e MissingGroups) Groups() []UID {
	return e.groups
}
