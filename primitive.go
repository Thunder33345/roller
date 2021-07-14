package ranker

//Group represent a collection of permissions and metadata
type Group struct {
	//Name is a display friendly name of the group, it should only be used for display, it has no filtering or requirement
	Name string
	//RefName is a command friendly runtime stable name, should be unique
	RefName string
	//UID is the unique identifier for this group used for saving and referencing, must never be changed
	UID string
	//Order dictates the overwrite precedent, where largest gets overwritten by smallest, must be unique
	Order int
	//Permission is the permission that is used
	Permission Entry
}

//Entry represent a collection of permissions and flags
type Entry struct {
	//EmptySet will discard all previously granted permissions
	EmptySet bool
	//Level is the default power level of said entry
	//Only the last group's level is in used, and context level overwrites group
	//if 0, the last value will be used instead
	Level int
	//Grant will add a permission
	Grant []string
	//Revoke will subtract a permissions that is granted by a prior group
	Revoke []string
}

//RawList is the raw save state for List
type RawList struct {
	//Overwrites has the highest precedent
	//Will overwrite all group based permissions
	Overwrites Entry
	//Groups are a list of p_group to inherit permission from
	Groups []string
}

//List is the compiled result from a RawList
type List struct {
	//Level is the final applicable level
	Level int
	//Permission is th final applicable permission
	Permission []string
}