package roller

//GroupProvider is something that's capable of providing a permission group
type GroupProvider interface {
	//Group will take the gid and return the Group
	//returns an error if there's an issue accessing Group
	Group(gid string) (Group, error)
	//Flag will take the group's id and flag's id to retrieve FlagEntry
	//If FlagEntry cannot be found, NewNotFoundError should be returned
	//IsErrorNotExist is used by the processor to suppress missing flags
	Flag(gid string, fid string) (FlagEntry, error)
}
