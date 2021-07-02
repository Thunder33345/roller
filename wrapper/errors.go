package wrapper

import "fmt"

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
