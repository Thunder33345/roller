package wrapper

import (
	"ranker"
	"time"
)

type Provider interface {
	GetRawPermissible(uid string) (ranker.RawPermissible, error)
	GetGroup(uid string) (ranker.Group, bool)
	LastChanged() int64
}

var _ Provider = (*MemoryProvider)(nil)

type MemoryProvider struct {
	groups      map[string]ranker.Group
	raw         map[string]ranker.RawPermissible
	defaultRaw  *ranker.RawPermissible
	lastChanged int64
}

func (p *MemoryProvider) GetRawPermissible(uid string) (ranker.RawPermissible, error) {
	if r, ok := p.raw[uid]; ok {
		return r, nil
	}
	if p.defaultRaw == nil {
		return ranker.RawPermissible{}, MissingPermissible{uid: uid}
	}
	return *p.defaultRaw, nil
}

func (p *MemoryProvider) GetGroup(uid string) (ranker.Group, bool) {
	if g, ok := p.groups[uid]; ok {
		return g, true
	}
	return ranker.Group{}, false
}

func (p *MemoryProvider) SetGroup(group ranker.Group) {
	p.groups[group.UID] = group
	p.lastChanged = time.Now().Unix()
}

func (p *MemoryProvider) LastChanged() int64 {
	return p.lastChanged
}
