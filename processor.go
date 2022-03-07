package roller

import (
	"sort"
)

//Processor is something that processes RawList into List with help of Provider retrieving Group
//Processor defines the behaviour of how the final RawList is generated
type Processor interface {
	//Process generates a List out of RawList, returns error if there's any problem
	Process(r RawList) (List, error)
	//ProcessFlags generates a List out of RawList with flags included, returns error if there's any problem
	ProcessFlags(r RawList, flags ...string) (List, error)
	//MergeEntry merges processed List with a list of permission Entry to generate a new List
	//this should not alter the original list
	MergeEntry(l List, es ...Entry) List
}

var _ Processor = (*BasicProcessor)(nil)

//BasicProcessor is the simple default Processor
//todo short circuit for Entry.EmptySet in Group.Permission
type BasicProcessor struct {
	//Provider retrieve Group data specified in RawList
	Provider GroupProvider
	//WeightAscending controls the order of applying the roles
	//false: the larger will be applied after the smaller(larger has more precedence)
	//true: the smaller will be applied after the larger(smaller has more precedence)
	WeightAscending bool
}

//compare is a helper function to sort misc stuff
func (p BasicProcessor) compare(i, j int) bool {
	if p.WeightAscending {
		return i > j
	}
	return i < j
}

//Process generates a List out of RawList,
//returns error if there's any issues with the provider
func (p BasicProcessor) Process(r RawList) (List, error) {
	gs, err := p.getGroups(r.Groups)
	if err != nil {
		return List{}, err
	}
	sort.Slice(gs, func(i, j int) bool {
		return p.compare(gs[i].Weight, gs[j].Weight)
	})

	var l List
	for _, g := range gs {
		l = p.processSet(l, g.Permission)
	}
	l = p.processSet(l, r.Overwrites)
	return l, nil
}

//ProcessFlags generates a List out of RawList with flags included,
//returns error if there's any issues with the provider
func (p BasicProcessor) ProcessFlags(r RawList, flags ...string) (List, error) {
	gs, err := p.getGroups(r.Groups)
	if err != nil {
		return List{}, err
	}
	sort.Slice(gs, func(i, j int) bool {
		return p.compare(gs[i].Weight, gs[j].Weight)
	})

	var l List
	for _, g := range gs {
		pre, post, err2 := p.getFlags(g.id, flags)
		if err2 != nil {
			return List{}, err2
		}
		for _, v := range pre {
			l = p.processSet(l, v.Entry)
		}
		l = p.processSet(l, g.Permission)
		for _, v := range post {
			l = p.processSet(l, v.Entry)
		}
	}
	l = p.processSet(l, r.Overwrites)
	return l, nil
}

//MergeEntry merges  List with a list of Entry to generate a new List
func (p BasicProcessor) MergeEntry(l List, es ...Entry) List {
	for _, e := range es {
		l = p.processSet(l, e)
	}
	return l
}

//getGroups returns a list of Group from a list of GroupID
//returns error if there's any issues with the provider
func (p BasicProcessor) getGroups(r []string) ([]keyedGroup, error) {
	var gs []keyedGroup
	for _, gid := range r {
		v, err := p.Provider.Group(gid)
		if err != nil {
			return nil, providerGroupError{
				group: gid,
				cause: err,
			}
		}
		gs = append(gs, keyedGroup{id: gid, Group: v})
	}
	return gs, nil
}

//processSet is a function that merges a List and an Entry to generate a new List
func (p BasicProcessor) processSet(l List, set Entry) List {
	if set.SetLevel {
		l.Level = set.Level
	} else {
		l.Level += set.Level
	}

	if set.EmptySet {
		l.Permission = []string{}
	} else if len(set.Revoke) > 0 {
		l.Permission = p.removeNodes(l.Permission, set.Revoke)
	}
	l.Permission = append(l.Permission, set.Grant...)
	return l
}

//getFlags takes a GroupID, and a list of flags
//and return a ordered, separated list of FlagEntry for pre-process and post-process
func (p BasicProcessor) getFlags(gid string, selected []string) (pre []FlagEntry, post []FlagEntry, xErr error) {
	fl := make([]FlagEntry, 0, len(selected))
	for _, sel := range selected {
		if f, found, err := p.Provider.Flag(gid, sel); err == nil {
			if found {
				fl = append(fl, f)
			}
		} else {
			return nil, nil, err
		}
	}
	sort.Slice(fl, func(i, j int) bool {
		return p.compare(fl[i].Weight, fl[j].Weight)
	})
	for _, f := range fl {
		if f.Preprocess {
			pre = append(pre, f)
		} else {
			post = append(post, f)
		}
	}
	return pre, post, nil
}

//removeNodes removes needles from a specified stack,
//inputs will not be altered and should be nonzero length slices
func (p BasicProcessor) removeNodes(stack []string, needle []string) []string {
	check := func(v string) bool {
		for _, r := range needle {
			if v == r {
				return false
			}
		}
		return true
	}
	ret := make([]string, 0, len(stack))
	for _, v := range stack {
		if check(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

type keyedGroup struct {
	id string
	Group
}
