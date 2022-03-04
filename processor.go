package roller

import (
	"sort"
)

//Processor is something that processes RawList List and Group
//Processor defines the rule set of how something get processed
type Processor interface {
	//Process generates a List out of RawList, returns error if there's any problem
	Process(r RawList) (List, error)
	//ProcessFlags generates a List out of RawList with flags included, returns error if there's any problem
	ProcessFlags(r RawList, flags ...string) (List, error)
	//MergeEntry merges List with a list of Entry to generate a new RawList
	MergeEntry(l List, es ...Entry) List
}

var _ Processor = (*BasicProcessor)(nil)

type BasicProcessor struct {
	Provider GroupProvider
	//WeightAscending controls whether smaller or bigger number holds precedent
	//by default the larger will overwrite the smaller
	WeightAscending bool
}

func (p BasicProcessor) compare(i, j int) bool {
	if p.WeightAscending {
		return i > j
	}
	return i < j
}

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
		pre, post, err2 := p.getFlags(g.ID, flags)
		if err2 != nil {
			if IsMissingFlagError(err2) {
				continue
			} else {
				return List{}, err2
			}
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

func (p BasicProcessor) MergeEntry(l List, es ...Entry) List {
	for _, e := range es {
		l = p.processSet(l, e)
	}
	return l
}

func (p BasicProcessor) getGroups(r []string) ([]Group, error) {
	var gs []Group
	for _, gid := range r {
		v, err := p.Provider.Group(gid)
		if err == nil {
			gs = append(gs, v)
		} else {
			return []Group{}, providerGroupError{
				group: gid,
				cause: err,
			}
		}
	}
	return gs, nil
}

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

//getFlags tries to get all selected flags from the map then return the sorted slice into preprocess and postprocess
func (p BasicProcessor) getFlags(gid string, selected []string) (pre []FlagEntry, post []FlagEntry, xErr error) {
	fl := make([]FlagEntry, 0, len(selected))
	for _, sel := range selected {
		if f, err := p.Provider.Flag(gid, sel); err == nil {
			fl = append(fl, f)
		} else {
			if IsMissingFlagError(err) {
				continue
			} else {
				return nil, nil, err
			}
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
//inputs will not be altered and are assumed to be nonzero length slices
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
