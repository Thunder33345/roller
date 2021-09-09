package provider

import (
	"github.com/Thunder33345/roller"
	"io"
)

var _ roller.GroupProvider = (GroupStorer)(nil)

//GroupStorer is something that is capable of store and provide groups
type GroupStorer interface { //todo set a better name
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
type Reloader interface { //todo set a better name
	//Reload will reload the provider
	Reload() error
}

//truncateSeeker is an io.ReadWriter that can be seeked and truncated
//internally used for compat on JSON.Save where Truncate(0) and Seek(0,0) will be called before writing
type truncateSeeker interface {
	Truncate(size int64) error
	io.Seeker
}
