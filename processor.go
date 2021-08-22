package perms_manager

import "sort"

//Processor is something that processes RawList List and Group
//Processor defines the rule set of how something get processed
type Processor interface {
	//Process generates a List out of RawList, returns error if there's any problem
	Process(r RawList) (List, error)
	//MergeEntry merges List with a list of Entry to generate a new RawList
	MergeEntry(l List, es ...Entry) List
}

type BasicProcessor struct {
	Provider GroupProvider
}

func NewProcessor(provider GroupProvider) BasicProcessor {
	return BasicProcessor{Provider: provider}
}
func (p BasicProcessor) Process(r RawList) (List, error) {
	gs, err := p.getGroups(r.Groups)
	if err != nil {
		return List{}, err
	}
	sort.Slice(gs, func(i, j int) bool {
		return gs[i].Order > gs[j].Order
	})

	var l List
	for _, g := range gs {
		l = p.processSet(l, g.Permission)
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
	var missing []string
	for _, uid := range r {
		if v, ok := p.Provider.GetGroup(uid); ok {
			gs = append(gs, v)
		} else {
			missing = append(missing, uid)
		}
	}
	if len(missing) > 0 {
		return []Group{}, MissingGroupsError{groups: missing}
	}
	return gs, nil
}

func (p BasicProcessor) processSet(l List, set Entry) List {
	if lv := set.Level; !set.IgnoreLevel {
		if set.AddLevel {
			l.Level += lv
		} else {
			l.Level = lv
		}
	}
	if set.EmptySet {
		l.Permission = []string{}
	} else {
		p.removeNodes(l.Permission, set.Revoke)
	}
	l.Permission = append(l.Permission, set.Grant...)
	return l
}

func (p BasicProcessor) removeNodes(stack []string, needle []string) {
	for i, s := range stack {
		for _, r := range needle {
			if s == r {
				stack = stack[i-1 : i+1]
			}
		}
	}
}
