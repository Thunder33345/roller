package ranker

//GroupProvider is something that's capable of providing a permission group
type GroupProvider interface {
	//GetGroup will take an uid and return the Group
	//returns bool as false when Group is not found
	GetGroup(uid string) (Group, bool)
}

//GroupProviderFunc is a function type that also makes it compliant to the GroupProvider interface
type GroupProviderFunc func(uid string) (Group, bool)

var _ GroupProvider = (*GroupProviderFunc)(nil)

func (g GroupProviderFunc) GetGroup(uid string) (Group, bool) {
	return g(uid)
}
