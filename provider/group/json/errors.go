package json

import (
	"fmt"
)

var _ error = (*GroupNotFoundError)(nil)

type GroupNotFoundError struct {
	id string
}

func NewGroupNotFoundError(id string) GroupNotFoundError {
	return GroupNotFoundError{id: id}
}

func (e GroupNotFoundError) Error() string {
	return fmt.Sprintf("group ID \"%s\" cant be found", e.id)
}

func (e GroupNotFoundError) ID() string {
	return e.id
}

var _ error = (*ReadOnlyError)(nil)

type ReadOnlyError struct{}

func (e ReadOnlyError) Error() string {
	return "provider is set to readonly mode"
}
