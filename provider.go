package roller

//GroupProvider is something that's capable of providing a permission group
type GroupProvider interface {
	//GetGroup will take an uid and return the Group
	//returns bool as false when Group is not found
	GetGroup(uid string) (Group, bool)
}
