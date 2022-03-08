package json

import (
	"encoding/json"
	"github.com/Thunder33345/roller"
	"github.com/Thunder33345/roller/provider"
	"io"
	"sync"
)

var _ provider.Provider = (*JSON)(nil)

type JSON struct {
	groups map[string]*groupData
	//file is where the data will be read and written to
	//io.Closer is supported and will be closed when JSON.Close is called
	file io.ReadWriter
	//allowUnknown suppresses un known fields
	allowUnknown bool
	//readOnly stops the configuration from being altered or being saved to disk
	readOnly bool
	//indent is the key to use when writing out
	indent string
	m      sync.RWMutex
}

func New(file io.ReadWriter, options ...Option) (*JSON, error) {
	j := &JSON{file: file, indent: "\t", groups: make(map[string]*groupData)}
	for _, option := range options {
		option(j)
	}
	if err := j.load(); err != nil {
		return nil, err
	}
	return j, nil
}

func (j *JSON) Group(groupID string) (roller.Group, error) {
	j.m.RLock()
	defer j.m.RUnlock()
	d, ok := j.groups[groupID]
	if ok {
		return d.Group, nil
	}
	return roller.Group{}, groupNotFoundError{id: groupID}
}

func (j *JSON) Flag(gid string, fid string) (roller.FlagEntry, bool, error) {
	j.m.RLock()
	defer j.m.Unlock()
	d, ok := j.groups[gid]
	if !ok {
		return roller.FlagEntry{}, false, groupNotFoundError{id: gid}
	}
	f, ok := d.Flags[fid]
	return f, ok, nil
}

func (j *JSON) SetGroup(groupID string, group roller.Group) error {
	j.m.Lock()
	defer j.m.Unlock()
	if j.readOnly {
		return readOnlyError{}
	}
	j.groups[groupID] = &groupData{
		Flags: nil,
		Group: group,
	}
	return nil
}

func (j *JSON) RemoveGroup(id string) error {
	j.m.Lock()
	defer j.m.Unlock()
	if j.readOnly {
		return readOnlyError{}
	}
	delete(j.groups, id)
	return nil
}

func (j *JSON) SetFlag(groupID string, flagID string, flag roller.FlagEntry) error {
	j.m.Lock()
	defer j.m.Unlock()
	if j.readOnly {
		return readOnlyError{}
	}
	d, ok := j.groups[groupID]
	if !ok {
		return groupNotFoundError{id: groupID}
	}
	if d.Flags == nil {
		d.Flags = make(map[string]roller.FlagEntry)
	}

	d.Flags[flagID] = flag
	return nil
}

func (j *JSON) RemoveFlag(groupID string, flagID string) error {
	j.m.Lock()
	defer j.m.Unlock()
	if j.readOnly {
		return readOnlyError{}
	}
	d, ok := j.groups[groupID]
	if !ok {
		return groupNotFoundError{id: groupID}
	}
	delete(d.Flags, flagID)
	return nil
}

func (j *JSON) WalkGroup(f func(group roller.Group, last bool) (halt bool)) error {
	j.m.RLock()
	defer j.m.RUnlock()
	i := 0
	for _, g := range j.groups {
		i++
		halt := f(g.Group, i >= len(j.groups))
		if halt {
			return nil
		}
	}
	return nil
}

func (j *JSON) WalkFlags(groupID string, f func(flag roller.FlagEntry, last bool) (halt bool)) error {
	j.m.RLock()
	defer j.m.RUnlock()
	d, ok := j.groups[groupID]
	if !ok {
		return groupNotFoundError{id: groupID}
	}
	i := 0
	for _, flag := range d.Flags {
		i++
		halt := f(flag, i >= len(d.Flags))
		if halt {
			return nil
		}
	}
	return nil
}

func (j *JSON) load() error {
	dec := json.NewDecoder(j.file)
	if !j.allowUnknown {
		dec.DisallowUnknownFields()
	}

	switch t := j.file.(type) {
	case io.Seeker:
		if _, err := t.Seek(0, 0); err != nil {
			return err
		}
	}

	if !dec.More() {
		return nil
	}

	var load groupDataSave
	if err := dec.Decode(&load); err != nil {
		return err
	}
	j.groups = load.Groups

	return nil
}

func (j *JSON) Save() error {
	j.m.RLock()
	defer j.m.RUnlock()
	if j.readOnly {
		return readOnlyError{}
	}

	switch t := j.file.(type) {
	case truncateSeeker:
		if err := t.Truncate(0); err != nil {
			return err
		}
		if _, err := t.Seek(0, 0); err != nil {
			return err
		}
	case reseter:
		t.Reset()
	}

	enc := json.NewEncoder(j.file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", j.indent)

	var save groupDataSave

	save.Groups = j.groups

	if err := enc.Encode(save); err != nil {
		return err
	}
	return nil
}

func (j *JSON) Close() error {
	j.groups = nil
	if c, ok := j.file.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

type groupData struct {
	Flags map[string]roller.FlagEntry `json:"flags,omitempty"`
	roller.Group
}

//groupDataSave is the json save structure
//todo make a proper ordered save system
type groupDataSave struct {
	Groups      map[string]*groupData `json:"groups,omitempty"`
	GroupsOrder []string              `json:"groups_order,omitempty"`
}
