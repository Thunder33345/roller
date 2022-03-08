package group

import (
	"github.com/Thunder33345/roller"
)

var _ roller.GroupProvider = (Provider)(nil)

//Provider store and provide groups
type Provider interface {
	roller.GroupProvider
	SetGroup(groupID string, group roller.Group) error
	RemoveGroup(groupID string) error
	SetFlag(groupID string, flagID string, flag roller.FlagEntry) error
	RemoveFlag(groupID string, flagID string) error
}

type Saver interface {
	Provider
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
