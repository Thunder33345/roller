package roller

//GroupProvider is something that's capable of providing a permission group
type GroupProvider interface {
	//GetGroup will take the gid and return the Group
	//returns an error if there's an issue accessing Group
	GetGroup(gid string) (Group, error)
}
