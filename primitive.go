package roller

//Group represent a collection of permissions and metadata
type Group struct {
	//Name is a display friendly name of the group, it should only be used for display, it has no filtering or requirement
	Name string `json:"name"`
	//RefName is a command friendly runtime stable name, should be unique
	RefName string `json:"ref_name"`
	//ID is the unique identifier for this group used for saving and referencing, must never be changed
	ID string `json:"id"`
	//Weight dictates the overwriting precedent, where the larger overwrites the smaller
	//must be unique, otherwise behaviour is undefined
	Weight int `json:"weight"`
	//Permission is the permission that is used
	Permission Entry `json:"permission,omitempty"`
}

//Entry represent a collection of permissions and flags
type Entry struct {
	//EmptySet will discard all previously granted permissions
	EmptySet bool `json:"empty_set,omitempty"`
	//Level is the default power level of said entry
	//Only the highest group's level is in used
	Level int `json:"level,omitempty"`
	//SetLevel makes overwrites the Level instead of adding or subtracting from last level
	SetLevel bool `json:"set_level,omitempty"`
	//Grant will add permissions to the List
	Grant []string `json:"grant,omitempty"`
	//Revoke will revoke a permissions that is granted to the List by a prior group
	Revoke []string `json:"revoke,omitempty"`
}

//FlagEntry is an Entry but inside a Group.Flags
//its same as Entry byt with extra flag only fields
type FlagEntry struct {
	//Weight dictates the overwriting precedent, must be unique, otherwise behaviour is undefined
	//behaviour is defined by processor
	Weight int `json:"weight"`
	//Preprocess indicates that this should be processed before Group.Permission
	Preprocess bool `json:"preprocess,omitempty"`
	Entry
}

//RawList is the raw save state for List
type RawList struct {
	//Overwrites has the highest precedent
	//Will overwrite all group based permissions
	Overwrites Entry `json:"overwrites,omitempty"`
	//Groups are a list of group UUID to inherit permission from
	Groups []string `json:"groups,omitempty"`
}

//List is the compiled result from a RawList
type List struct {
	//Level is the final applicable level
	Level int `json:"level,omitempty"`
	//Permission is th final applicable permission
	Permission []string `json:"permission,omitempty"`
}
