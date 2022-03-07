package roller

//Group represent a collection of permissions and metadata
type Group struct {
	//Name is a display friendly name of the group, it should only be used for display, it has no filtering or requirement
	Name string `json:"name"`
	//RefName is a command friendly runtime stable name, should be unique
	RefName string `json:"ref_name"`
	//Weight dictates the order of overwriting precedent, by default the larger gets applied after the smaller ones
	//must be unique, otherwise behaviour is undefined
	//this does not define the rank/hierarchy of a group, use Entry.Level instead
	Weight int `json:"weight"`
	//Permission will be given for said group holder
	Permission Entry `json:"permission,omitempty"`
}

//Entry represent a collection of permissions and flags
type Entry struct {
	//EmptySet discard all previously granted permissions
	EmptySet bool `json:"empty_set,omitempty"`
	//Level will be added into List.Level, this can be used to compare hierarchy
	//negatives will subtract from level instead
	Level int `json:"level,omitempty"`
	//SetLevel sets and overwrite the Level of List.Level instead of adding or adding
	SetLevel bool `json:"set_level,omitempty"`
	//Grant adds permissions to the List
	//grants will be processed after revokes
	Grant []string `json:"grant,omitempty"`
	//Revoke will revoke a permissions that are granted to the List by a prior group
	Revoke []string `json:"revoke,omitempty"`
}

//FlagEntry is an Entry but as part of a Group's Flags
//its same as Entry byt with extra flag only fields
type FlagEntry struct {
	//Weight dictates the order of overwriting precedent, by default the larger gets applied after the smaller ones
	//must be unique, otherwise behaviour is undefined
	Weight int `json:"weight"`
	//Preprocess indicates that this should be processed before Group.Permission
	Preprocess bool `json:"preprocess,omitempty"`
	Entry
}

//RawList is the preprocessed save state for List
type RawList struct {
	//Overwrites has the highest precedent
	//Will be applied after all group based permissions
	Overwrites Entry `json:"overwrites,omitempty"`
	//Groups are a list of group UUID to inherit permission from
	Groups []string `json:"groups,omitempty"`
}

//List is the compiled result from a RawList
type List struct {
	//Level is the final applicable level
	Level int `json:"level,omitempty"`
	//Permission is the list of permissions that are applicable
	//only grants that aren't revoked will be stored in the list
	Permission []string `json:"permission,omitempty"`
}
