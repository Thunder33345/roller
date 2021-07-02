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

var _ error = (*MissingPermissible)(nil) // ensure MissingPermissible implements error

type MissingPermissible struct {
	uid string
}

func (e MissingPermissible) Error() string {
	return fmt.Sprintf("Missing Permissible: %v", e.uid)
}

func (e MissingPermissible) UID() string {
	return e.uid
}
