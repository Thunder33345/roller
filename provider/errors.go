package provider

import (
	"errors"
	"fmt"
)

var ReadOnlyError = errors.New("provider is set to read only")

var _ error = (*groupNotFoundError)(nil)

func NewGroupNotFoundError(id string) error {
	return &groupNotFoundError{id: id}
}

type groupNotFoundError struct {
	id string
}

func (e groupNotFoundError) Error() string {
	return fmt.Sprintf("group ID \"%s\" cant be found", e.id)
}
