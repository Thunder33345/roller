package ranker

import (
	"time"
)

type GroupProvider interface {
	GetGroup(uid string) (Group, bool)
}

type GroupProviderFunc func(uid string) (Group, bool)

func (g GroupProviderFunc) GetGroup(uid string) (Group, bool) {
	return g(uid)
}

type Provider interface {
	GetRawPermissible(uid string) (RawPermissible, error)
	GetGroup(uid string) (Group, bool)
	LastChanged() int64
}

var _ Provider = (*MemoryProvider)(nil)

type MemoryProvider struct {
	groups      map[string]Group
	raw         map[string]RawPermissible
	defaultRaw  *RawPermissible
	lastChanged int64
}

func (p *MemoryProvider) GetRawPermissible(uid string) (RawPermissible, error) {
	if r, ok := p.raw[uid]; ok {
		return r, nil
	}
	if p.defaultRaw == nil {
		return RawPermissible{}, MissingPermissible{uid: uid}
	}
	return *p.defaultRaw, nil
}

func (p *MemoryProvider) GetGroup(uid string) (Group, bool) {
	if g, ok := p.groups[uid]; ok {
		return g, true
	}
	return Group{}, false
}

func (p *MemoryProvider) SetGroup(group Group) {
	p.groups[group.UID] = group
	p.lastChanged = time.Now().Unix()
}

func (p *MemoryProvider) LastChanged() int64 {
	return p.lastChanged
}
