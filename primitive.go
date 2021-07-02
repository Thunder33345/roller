package ranker

//Group represent a collection of permissions and metadata
type Group struct {
	//Name is a display friendly name of the group, it should only be used for display, it has no filtering or requirement
	Name string
	//RefName is a command friendly runtime name, should be unique
	RefName string
	//UID is the unique identifier for this group used for saving and referencing, should never be changed
	UID string
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
	//if 0, the last value will be used instead
	Level int
	//Grants will give(or overwrite) a permission by
	//string will be the key of the permission
	Grants []string
	//Revoke will remove any permissions that is granted by a prior group
	Revoke []string
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
	Groups []string
}

//Permissible is the compiled result from a RawPermissible
type Permissible struct {
	//Level is the final applicable level
	Level int
	//Permission is th final applicable permission
	Permission []string
}

type WrappedPermissible struct {
	permissible Permissible
	judge       Judge
}

func (w WrappedPermissible) HasPermission(node string) bool {
	return w.judge.HasPermission(w.permissible, node)
}

func (w WrappedPermissible) HasPermissionWithLevel(node string, level int) bool {
	return w.judge.HasPermissionWithLevel(w.permissible, node, level)
}

func (w WrappedPermissible) IsHigherLevel(subject Permissible) bool {
	return w.judge.IsHigherLevel(w.permissible, subject)
}
