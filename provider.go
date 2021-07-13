package ranker

type GroupProvider interface {
	GetGroup(uid string) (Group, bool)
}

type GroupProviderFunc func(uid string) (Group, bool)

func (g GroupProviderFunc) GetGroup(uid string) (Group, bool) {
	return g(uid)
}
