package roller

//GroupProvider is something that provides group permission data
type GroupProvider interface {
	//Group will take the gid and return the Group
	//returns an error if there's an issue accessing Group
	Group(gid string) (Group, error)
	//Flag will take the group's id and flag's id to retrieve FlagEntry
	//error indicates the presence of an error,
	//bool indicates if an FlagEntry is found
	//On success bool should be true and error should be nil
	Flag(gid string, fid string) (FlagEntry, bool, error)
}
