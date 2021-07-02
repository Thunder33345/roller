package ranker

//Level is the numerical representation if hierarchy
//A higher number(100) power level will overrule a lower number(1)
//In context of ContextPermissionSet: A 0 will inherit group's default
//In context of Group: A 0 will inherit a previous group's value
type Level int

//Nodes is a slice of permission nodes
type Nodes []string


//UID is a unique identifier given for a certain group
type UID string

//Group represent a collection of permissions and metadata
type Group struct {
	//Name is a display friendly name of the group, it should only be used for display, it has no filtering or requirement
	Name string
	//RefName is a command friendly runtime name, should be unique
	RefName string
	//UID is the unique identifier for this group used for saving and referencing, should never be changed
	UID UID
	//Order dictates the overwrite precedent, where largest gets overwritten by smallest, must be unique
	Order int
	//Default is the default permission that is used
	Default PermissionSet
	//Context is permissions that are only applicable in certain context
	//Where string will be the identifier of the context
	//Unless context have ContextPermissionSet.IgnoreDefault set to true, group permission will be inherited
	//Context map[string]ContextPermissionSet
	//ContextFallback will be used if Context is not found
	//If this is included, it will always have the highest precedence
	//ContextFallback ContextPermissionSet
}

//PermissionSet represent a collection of permissions and flags
type PermissionSet struct {
	//EmptySet will discard all previously granted permissions
	EmptySet bool
	//Level is the default power level of said entry
	//Only the last group's level is in used, and context level overwrites group
	Level Level
	//Grants will give(or overwrite) a permission by
	//string will be the key of the permission
	Grants Nodes
	//Revoke will remove any permissions that is granted by a prior group
	Revoke Nodes
}

//ContextPermissionSet represent a collection of permission with extra flags only valid in context
//type ContextPermissionSet struct {
//	//PermissionSet embeds the default set
//	PermissionSet
//	//IgnoreDefault will ignore the default group permission for this context
//	IgnoreDefault bool
//	//Order is used when merging multiple context, must be unique
//	//When overwriting precedent is based on largest to smallest
//	Order int
//}

//RawPermissible is the raw save state for Permissible
type RawPermissible struct {
	//Overwrites has the highest precedent
	//Will overwrite all group based permissions
	Overwrites PermissionSet
	//Groups are a list of p_group to inherit permission from
	Groups []UID
}

//Permissible is the compiled result from a RawPermissible
type Permissible struct {
	//Level is the final applicable level
	Level Level
	//Permission is th final applicable permission
	Permission Nodes
}
