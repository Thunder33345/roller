package provider

import (
	"github.com/Thunder33345/roller"
)

var _ roller.GroupProvider = (GroupStorer)(nil)

//GroupStorer is something that is capable of store and provide groups
type GroupStorer interface {
	AddGroup(group roller.Group) error
	Group(id string) (roller.Group, error)
	RemoveGroup(id string) error
}

//Walker is an iterable provider
type Walker interface {
	//WalkGroup will iterate through all the groups with provided callback
	//if the function returns true, it will halt the process
	WalkGroup(func(roller.Group) (halt bool)) error
}

//Saver is a provider that needs manual saving
type Saver interface {
	//Save will flush and save internal state
	Save() error
}

//Closer is a provider that needs cleaning up
//Behaviour is undefined if used after closing
type Closer interface {
	//Close will make the close and free all relevant data
	Close() error
}

//Reloader is a provider that's capable of reinitialize its internal state
type Reloader interface {
	//Reload will reload the provider
	Reload() error
}
