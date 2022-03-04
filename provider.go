package roller

import "errors"

//GroupProvider is something that's capable of providing a permission group
type GroupProvider interface {
	//Group will take the gid and return the Group
	//returns an error if there's an issue accessing Group
	Group(gid string) (Group, error)
	Flag(gid string, fid string) (FlagEntry, error)
}

func ErrorIsNotExist(err error) bool {
	return errors.Is(err, notFoundError{})
}

func NewNotFoundError(err error) error {
	return notFoundError{err: err}
}

type notFoundError struct {
	err error
}

func (e notFoundError) Error() string {
	return e.err.Error()
}

func (e notFoundError) Unwrap() error {
	return e.err
}
func (e notFoundError) Is(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}
