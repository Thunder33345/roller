package provider

import (
	"fmt"
	"github.com/Thunder33345/roller"
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

var _ error = (*DuplicateGroupIDError)(nil)

type DuplicateGroupIDError struct {
	g1 roller.Group
	g2 roller.Group
}

func NewDuplicateIDError(original roller.Group, duplicate roller.Group) DuplicateGroupIDError {
	return DuplicateGroupIDError{g1: original, g2: duplicate}
}

func (e DuplicateGroupIDError) Error() string {
	return fmt.Sprintf("group ID not unique: ID \"%s\"(%s[#%s]) already exist, "+
		"cant be shared with ID \"%s\"(%s[#%s])", e.g1.ID, e.g1.Name, e.g1.RefName, e.g2.ID, e.g2.Name, e.g2.RefName)
}

func (e DuplicateGroupIDError) Original() roller.Group {
	return e.g1
}

func (e DuplicateGroupIDError) Duplicate() roller.Group {
	return e.g2
}

var _ error = (*ReadOnlyError)(nil)

type ReadOnlyError struct{}

func (e ReadOnlyError) Error() string {
	return "provider is set to readonly mode"
}
