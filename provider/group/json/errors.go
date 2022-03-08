package json

import (
	"fmt"
)

var _ error = (*groupNotFoundError)(nil)

type groupNotFoundError struct {
	id string
}

func (e groupNotFoundError) Error() string {
	return fmt.Sprintf("group ID \"%s\" cant be found", e.id)
}

var _ error = (*readOnlyError)(nil)

type readOnlyError struct{}

func (e readOnlyError) Error() string {
	return "provider is set to readonly mode"
}
