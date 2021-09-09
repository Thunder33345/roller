package roller

//GroupProvider is something that's capable of providing a permission group
type GroupProvider interface {
	//Group will take the gid and return the Group
	//returns an error if there's an issue accessing Group
	Group(gid string) (Group, error)
}
