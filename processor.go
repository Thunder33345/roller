package ranker

import "sort"

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
		return []Group{}, MissingGroups{groups: missing}
	}
	return gs, nil
}

func Process() Processor {
	return ProcessorFunc(process)
}

func process(r RawPermissible, pr GroupProvider) (Permissible, error) {
	gs, err := pGroup(r.Groups, pr)
	if err != nil {
		return Permissible{}, err
	}
	sort.Slice(gs, func(i, j int) bool {
		return gs[i].Order > gs[j].Order
	})

	var p Permissible
	for _, g := range gs {
		p = pProcessSet(p, g.Default)
	}
	p = pProcessSet(p, r.Overwrites)
	return p, nil
}

func pProcessSet(p Permissible, set PermissionSet) Permissible {
	if l := set.Level; l != 0 {
		p.Level = l
	}
	if set.EmptySet {
		p.Permission = []string{}
	} else {
		pRemoveNodes(p.Permission, set.Revoke)
	}
	p.Permission = append(p.Permission, set.Grants...)
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
