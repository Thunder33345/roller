package provider

import (
	"encoding/json"
	"github.com/Thunder33345/roller"
	"io"
)

var _ GroupStorer = (*JSON)(nil)

type JSON struct {
	groups []roller.Group
	//file is where the data will be read and written to
	//io.Closer is supported and will be closed when JSON.Close is called
	file io.ReadWriter
	//allowUnknown suppresses un known fields
	allowUnknown bool
	//readOnly stops the configuration from being altered or being saved to disk
	readOnly bool
	//indent is the key to use when writing out
	indent string
	//unsafeSave suppresses duplicate uid check when saving
	//will still push the error down to next load
	unsafeSave bool
}

func NewJSON(file io.ReadWriter) (*JSON, error) {
	j := &JSON{file: file, indent: "\t"}
	if err := j.Load(); err != nil {
		return nil, err
	}
	return j, nil
}

func NewJSONWithOptions(file io.ReadWriter, allowUnknown bool, readOnly bool, indent string, unsafeSave bool) (*JSON, error) {
	j := &JSON{file: file, allowUnknown: allowUnknown, readOnly: readOnly, indent: indent, unsafeSave: unsafeSave}
	if err := j.Load(); err != nil {
		return nil, err
	}
	return j, nil
}

func (j *JSON) Group(id string) (roller.Group, error) {
	i, g := j.findGroup(id)
	if i >= 0 {
		return g, nil
	}
	return roller.Group{}, NewGroupNotFoundError(id)
}

func (j *JSON) AddGroup(group roller.Group) error {
	if j.readOnly {
		return ReadOnlyError{}
	}
	i, _ := j.findGroup(group.ID)
	if i >= 0 {
		j.groups[i] = group
		return nil
	}
	j.groups = append(j.groups, group)
	return nil
}

func (j *JSON) RemoveGroup(id string) error {
	if j.readOnly {
		return ReadOnlyError{}
	}
	i, _ := j.findGroup(id)
	if i >= 0 {
		j.groups = append(j.groups[:i], j.groups[i+1:]...)
		return nil
	}
	return NewGroupNotFoundError(id)
}

func (j *JSON) WalkGroup(f func(group roller.Group, last bool) (halt bool)) error {
	for i, g := range j.groups {
		halt := f(g, len(j.groups)-1 == i)
		if halt {
			return nil
		}
	}
	return nil
}

func (j *JSON) Load() error {
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

	var tg []roller.Group
	if err := dec.Decode(&tg); err != nil {
		return err
	}
	if err := j.duplicateCheck(tg); err != nil {
		return err
	}
	j.groups = tg
	return nil
}

func (j *JSON) Reload() error {
	c := j.groups
	j.groups = nil
	err := j.Load()
	if err != nil {
		j.groups = c
		return err
	}
	return nil
}

func (j *JSON) Save() error {
	if j.readOnly {
		return ReadOnlyError{}
	}
	if err := j.duplicateCheck(j.groups); err != nil && !j.unsafeSave {
		return err
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

	if err := enc.Encode(j.groups); err != nil {
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

func (j *JSON) findGroup(id string) (int, roller.Group) {
	for i, g := range j.groups {
		if g.ID == id {
			return i, g
		}
	}
	return -1, roller.Group{}
}

func (j *JSON) duplicateCheck(groups []roller.Group) error {
	found := make(map[string]int, len(groups))
	for i, g := range groups {
		di, exist := found[g.ID]
		if exist {
			og := groups[di]
			return NewDuplicateIDError(og, g)
		}
		found[g.ID] = i
	}
	return nil
}
