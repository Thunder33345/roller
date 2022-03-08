package provider

import (
	"github.com/Thunder33345/roller"
	"io"
)

var _ roller.GroupProvider = (Provider)(nil)

//Provider store and provide groups
type Provider interface {
	roller.GroupProvider
	SetGroup(groupID string, group roller.Group) error
	RemoveGroup(groupID string) error
	SetFlag(groupID string, flagID string, flag roller.FlagEntry) error
	RemoveFlag(groupID string, flagID string) error
	Save() error
}

type Lister interface {
	Provider
	Groups() (map[string]roller.Group, error)
	Flags(group roller.Group) (map[string]roller.FlagEntry, error)
}

//Walker is an iterable provider
type Walker interface {
	Provider
	//WalkGroups will iterate through all the groups with provided callback
	//func should take in group as the current group
	//and last to indicate if this is the last group
	//if last is true, the callback won't receive any more new groups, and the iteration will end
	//if the function returns halt as true, it will stop further iteration
	WalkGroups(func(group roller.Group, last bool) (halt bool)) error
	WalkFlags(roller.Group, func(Flag roller.FlagEntry, last bool) (halt bool)) error
}

//truncateSeeker is an io.ReadWriter that can be seeked and truncated
//internally used for compat on JSON.Save where Truncate(0) and Seek(0,0) will be called before writing
type truncateSeeker interface {
	Truncate(size int64) error
	io.Seeker
	io.ReadWriter
}

//reseter is an io.ReadWriter that can be reset before being written to, used for JSON.Save
type reseter interface {
	Reset()
	io.ReadWriter
}
