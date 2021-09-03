package provider

import (
	"github.com/Thunder33345/roller"
)

var _ roller.GroupProvider = (GroupStoreProvider)(nil)

//GroupStoreProvider is something that is capable of store and provide groups
type GroupStoreProvider interface {
	SetGroup(gid string, group roller.Group) error
	GetGroup(gid string) (roller.Group, error)
}

//ListStoreProvider is something that can store and provide lists
type ListStoreProvider interface {
	SetList(lid string, list roller.List) error
	GetList(lid string) (roller.List, error)
}

//FileProvider is something that lives on the file system
type FileProvider interface {
	//Close should be called before closing for cleaning up
	Close() error
	//Save will cause the provider to save all information to disk
	Save() error
	//Reload will cause the provider to reload all data from disk
	Reload() error
}
