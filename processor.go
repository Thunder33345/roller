package ranker

import "sort"

type Processor interface {
	Process(r RawPermissible, pr GroupProvider) (Permissible, error)
}

type ProcessorFunc func(r RawPermissible, pr GroupProvider) (Permissible, error)

func (p ProcessorFunc) Process(r RawPermissible, pr GroupProvider) (Permissible, error) {
	return p(r, pr)
}

type CachedProcessor interface {
	ClearCache()
	Process(uid string) (Permissible, error)
	DirectProcess(uid string) (Permissible, error)
}

func Process(r RawPermissible, pr GroupProvider) (Permissible, error) {
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

type ProcessorWithCache struct {
	cache       map[string]Permissible
	process     Processor
	provider    Provider
	lastChanged int64
}

func (p *ProcessorWithCache) Process(uid string) (Permissible, error) {
	if c, ok := p.GetCache(uid); ok {
		return c, nil
	}
	return p.DirectProcess(uid)
}

func (p *ProcessorWithCache) DirectProcess(uid string) (Permissible, error) {
	r, e := p.provider.GetRawPermissible(uid)
	if e != nil {
		return Permissible{}, e
	}
	return p.process.Process(r, p.provider)
}

func (p *ProcessorWithCache) GetCache(uid string) (Permissible, bool) {
	if p.provider.LastChanged() > p.lastChanged {
		p.ClearCache()
		p.lastChanged = p.provider.LastChanged()
		return Permissible{}, false
	}
	if c, ok := p.cache[uid]; ok {
		return c, true
	}
	return Permissible{}, false
}

func (p *ProcessorWithCache) StoreCache(uid string, permissible Permissible) {
	p.cache[uid] = permissible
}

func (p *ProcessorWithCache) ClearCache() {
	p.cache = map[string]Permissible{}
}
