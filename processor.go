package ranker

import "sort"

type Processor interface {
	Process(r RawPermissionList, pr GroupProvider) (PermissionList, error)
}

type ProcessorFunc func(r RawPermissionList, pr GroupProvider) (PermissionList, error)

func (p ProcessorFunc) Process(r RawPermissionList, pr GroupProvider) (PermissionList, error) {
	return p(r, pr)
}

type WrappedProcessor interface {
	Process(uid string) (PermissionList, error)
}

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
		p = pProcessSet(p, g.Default)
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

type CachedProcessor struct {
	cache       map[string]PermissionList
	processor   Processor
	provider    DataProvider
	lastChanged int64
}

func (p *CachedProcessor) Process(uid string) (PermissionList, error) {
	if c, ok := p.GetCache(uid); ok {
		return c, nil
	}
	return p.DirectProcess(uid)
}

func (p *CachedProcessor) DirectProcess(uid string) (PermissionList, error) {
	r, e := p.provider.GetRawPermission(uid)
	if e != nil {
		return PermissionList{}, e
	}
	return p.processor.Process(r, p.provider)
}

func (p *CachedProcessor) GetCache(uid string) (PermissionList, bool) {
	if p.provider.LastChanged() > p.lastChanged {
		p.ClearCache()
		p.lastChanged = p.provider.LastChanged()
		return PermissionList{}, false
	}
	if c, ok := p.cache[uid]; ok {
		return c, true
	}
	return PermissionList{}, false
}

func (p *CachedProcessor) StoreCache(uid string, pl PermissionList) {
	p.cache[uid] = pl
}

func (p *CachedProcessor) ClearCache() {
	p.cache = map[string]PermissionList{}
}
