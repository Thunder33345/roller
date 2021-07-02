package wrapper

import "ranker"

type CachedProcessor interface {
	ClearCache()
	Process(uid string) (ranker.Permissible, error)
	DirectProcess(uid string) (ranker.Permissible, error)
}

type ProcessorWithCache struct {
	cache       map[string]ranker.Permissible
	process     ranker.ProcessorFunc
	provider    Provider
	lastChanged int64
}

func (p *ProcessorWithCache) Process(uid string) (ranker.Permissible, error) {
	if c, ok := p.GetCache(uid); ok {
		return c, nil
	}
	return p.DirectProcess(uid)
}

func (p *ProcessorWithCache) DirectProcess(uid string) (ranker.Permissible, error) {
	r, e := p.provider.GetRawPermissible(uid)
	if e != nil {
		return ranker.Permissible{}, e
	}
	return p.process(r, p.provider)
}

func (p *ProcessorWithCache) GetCache(uid string) (ranker.Permissible, bool) {
	if p.provider.LastChanged() > p.lastChanged {
		p.ClearCache()
		p.lastChanged = p.provider.LastChanged()
		return ranker.Permissible{}, false
	}
	if c, ok := p.cache[uid]; ok {
		return c, true
	}
	return ranker.Permissible{}, false
}

func (p *ProcessorWithCache) StoreCache(uid string, permissible ranker.Permissible) {
	p.cache[uid] = permissible
}

func (p *ProcessorWithCache) ClearCache() {
	p.cache = map[string]ranker.Permissible{}
}
