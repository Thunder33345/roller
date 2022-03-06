package roller

//GroupProvider provides groups and flags data
type GroupProvider interface {
	//Group take the unique group ID to retrieve the Group
	//returns an error if Group does not exist, or there's an error accessing the data
	Group(gid string) (Group, error)
	//Flag take the group's id and flag's id to retrieve the FlagEntry
	//error indicates the presence of an error
	//bool indicates if an FlagEntry is found
	//On success bool should be true and error should be nil
	Flag(gid string, fid string) (FlagEntry, bool, error)
}
