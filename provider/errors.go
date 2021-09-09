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

var _ error = (*DuplicateIDError)(nil)

type DuplicateIDError struct {
	g1 roller.Group
	g2 roller.Group
}

func NewDuplicateIDError(original roller.Group, duplicate roller.Group) DuplicateIDError {
	return DuplicateIDError{g1: original, g2: duplicate}
}

func (e DuplicateIDError) Error() string {
	return fmt.Sprintf("unique ID not unique: UID \"%s\"(%s[#%s]) already exist, "+
		"cant be shared with UID \"%s\"(%s[#%s])", e.g1.UID, e.g1.Name, e.g1.RefName, e.g2.UID, e.g2.Name, e.g2.RefName)
}

func (e DuplicateIDError) Original() roller.Group {
	return e.g1
}

func (e DuplicateIDError) Duplicate() roller.Group {
	return e.g2
}

var _ error = (*ReadOnlyError)(nil)

type ReadOnlyError struct{}

func (e ReadOnlyError) Error() string {
	return "provider is set to readonly mode"
}
