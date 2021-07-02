package ranker

import "fmt"

var _ error = (*MissingGroups)(nil) // ensure MissingGroups implements error

type MissingGroups struct {
	groups []string
}

func (e MissingGroups) Error() string {
	return fmt.Sprintf("Missing Group: %v", e.groups)
}

func (e MissingGroups) Groups() []string {
	return e.groups
}
