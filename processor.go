package ranker

import "sort"

//Processor is something that can compile a raw permission list with a given group provider
type Processor interface {
	Process(r RawPermissionList, pr GroupProvider) (PermissionList, error)
}

type ProcessorFunc func(r RawPermissionList, pr GroupProvider) (PermissionList, error)

var _ Processor = (*ProcessorFunc)(nil)

func (p ProcessorFunc) Process(r RawPermissionList, pr GroupProvider) (PermissionList, error) {
	return p(r, pr)
}

//Process is the default built in permission list processing method
func Process(r RawPermissionList, pr GroupProvider) (PermissionList, error) {
	gs, err := pGroup(r.Groups, pr)
	if err != nil {
		return PermissionList{}, err
	}
	sort.Slice(gs, func(i, j int) bool {
		return gs[i].Order > gs[j].Order
	})

	var p PermissionList
	for _, g := range gs {
		p = pProcessSet(p, g.Permission)
	}
	p = pProcessSet(p, r.Overwrites)
	return p, nil
}

func pGroup(r []string, pr GroupProvider) ([]Group, error) {
	var gs []Group
	var missing []string
	for _, uid := range r {
		if v, ok := pr.GetGroup(uid); ok {
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

func pProcessSet(p PermissionList, set PermissionEntry) PermissionList {
	if l := set.Level; l != 0 {
		p.Level = l
	}
	if set.EmptySet {
		p.Permission = []string{}
	} else {
		pRemoveNodes(p.Permission, set.Revoke)
	}
	p.Permission = append(p.Permission, set.Grant...)
	return p
}

func pRemoveNodes(stack []string, needle []string) {
	for i, s := range stack {
		for _, r := range needle {
			if s == r {
				stack = stack[i-1 : i+1]
			}
		}
	}
}
