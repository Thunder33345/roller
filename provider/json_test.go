package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Thunder33345/roller"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestJSONSequence(t *testing.T) {
	r := require.New(t)
	sample := `[{"Name": "test","ID": "100"},{"Name": "test2","ID": "101"},{"Name": "test3","ID": "102"}]`
	file, err := os.CreateTemp("", "test*.json")
	defer func() {
		r.Nil(os.Remove(file.Name()))
	}()

	r.Nil(err)
	_, err = file.WriteString(sample)
	r.Nil(err)

	j, e := NewJSON(file)
	defer func() {
		r.Nil(j.Close())
	}()

	r.Nil(e)
	r.Equal(3, len(j.groups))
	r.Nil(j.AddGroup(roller.Group{Name: "test4", ID: "103"}))
	r.Equal(4, len(j.groups))
	r.Nil(j.Reload())
	r.Equal(3, len(j.groups))

	r.Nil(j.RemoveGroup("102"))
	r.Equal(2, len(j.groups))
	r.Nil(j.Save())
	r.Nil(j.Reload())
	r.Equal(2, len(j.groups))

	g, err := j.Group("100")
	r.Nil(err)
	r.Equal("test", g.Name)
}

func TestNewJSON(t *testing.T) {
	sample := `[{"Name": "test","ID": "100"},{"Name": "test2","ID": "101"}]`
	duplicatedSample := `[{"Name": "test","ID": "100"},{"Name": "test2","ID": "100"}]`
	syntax := `[{"Name": "test","ID": "100"},]`
	tests := []struct {
		name             string
		file             io.ReadWriter
		want             []roller.Group
		wantErr          bool
		wantDuplicateErr bool
		wantSyntaxErr    bool
	}{
		{
			name: "load empty array",
			file: bytes.NewBufferString("[]"),
			want: []roller.Group{},
		}, {
			name: "load empty",
			file: bytes.NewBufferString(""),
			want: []roller.Group(nil),
		}, {
			name: "load something",
			file: bytes.NewBufferString(sample),
			want: []roller.Group{{Name: "test", ID: "100"}, {Name: "test2", ID: "101"}},
		}, {
			name:             "load duplicated",
			file:             bytes.NewBufferString(duplicatedSample),
			wantDuplicateErr: true,
		}, {
			name:          "load bad syntax",
			file:          bytes.NewBufferString(syntax),
			wantSyntaxErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := NewJSON(tt.file)
			r := require.New(t)
			if tt.wantErr || tt.wantDuplicateErr || tt.wantSyntaxErr {
				r.Error(err)
				if tt.wantDuplicateErr {
					var dup DuplicateGroupIDError
					r.ErrorAs(err, &dup)
				}
				if tt.wantSyntaxErr {
					s := &json.SyntaxError{}
					r.ErrorAs(err, &s)
				}
				return
			}
			r.Nil(err)
			r.NotNil(j)
			r.Equal(tt.want, j.groups)
		})
	}
}

func TestJSON_Save(t *testing.T) {
	tests := []struct {
		name            string
		groups          []roller.Group
		readOnly        bool
		file            *bytes.Buffer
		wantErr         bool
		want            string
		wantDupeErr     bool
		wantReadOnlyErr bool
	}{
		{
			name:    "successful",
			groups:  []roller.Group{{Name: "test", ID: "100"}, {Name: "test2", ID: "101"}},
			file:    &bytes.Buffer{},
			wantErr: false,
			want: `[{"name":"test","ref_name":"","id":"100","weight":0,"permission":{}},{"name":"test2","ref_name":"","id":"101","weight":0,"permission":{}}]
`,
		}, {
			name:    "successful prefilled buffer",
			groups:  []roller.Group{{Name: "test", ID: "100"}},
			file:    bytes.NewBufferString("[1,2,3,4,5,6,7,8,9,0]"),
			wantErr: false,
			want: `[{"name":"test","ref_name":"","id":"100","weight":0,"permission":{}}]
`,
		}, {
			name:        "duplicated",
			groups:      []roller.Group{{Name: "test", ID: "100"}, {Name: "test2", ID: "100"}},
			file:        &bytes.Buffer{},
			wantDupeErr: true,
		}, {
			name:            "readonly",
			groups:          []roller.Group{{Name: "test", ID: "100"}, {Name: "test2", ID: "100"}},
			file:            &bytes.Buffer{},
			readOnly:        true,
			wantReadOnlyErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups:   tt.groups,
				file:     tt.file,
				readOnly: tt.readOnly,
			}
			r := require.New(t)
			err := j.Save()
			if tt.wantErr || tt.wantDupeErr || tt.wantReadOnlyErr {
				r.Error(err)
				if tt.wantDupeErr {
					var d DuplicateGroupIDError
					r.ErrorAs(err, &d)
				}
				if tt.wantReadOnlyErr {
					var e ReadOnlyError
					r.ErrorAs(err, &e)
				}
				return
			}
			r.Nil(err)
			r.Equal(tt.want, tt.file.String())
		})
	}
}

func TestJSON_Reload(t *testing.T) {

	tests := []struct {
		name       string
		groups     []roller.Group
		file       io.ReadWriter
		wantErr    bool
		wantGroups []roller.Group
	}{
		{
			name:       "rollback test",
			groups:     []roller.Group{{Name: "1", ID: "1"}, {Name: "2", ID: "2"}},
			file:       bytes.NewBufferString("[!!"),
			wantErr:    true,
			wantGroups: []roller.Group{{Name: "1", ID: "1"}, {Name: "2", ID: "2"}},
		}, {
			name:       "reload test",
			groups:     []roller.Group{{Name: "1", ID: "1"}, {Name: "2", ID: "2"}},
			file:       bytes.NewBufferString(`[{"Name": "test","ID": "100"},{"Name": "test2","ID": "101"}]`),
			wantErr:    false,
			wantGroups: []roller.Group{{Name: "test", ID: "100"}, {Name: "test2", ID: "101"}},
		}, {
			name:       "reload to nil",
			groups:     []roller.Group{{Name: "1", ID: "1"}, {Name: "2", ID: "2"}},
			file:       bytes.NewBufferString(`[]`),
			wantErr:    false,
			wantGroups: []roller.Group{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups: tt.groups,
				file:   tt.file,
			}

			r := require.New(t)
			err := j.Reload()
			r.Equal(tt.wantGroups, j.groups)

			if tt.wantErr {
				r.Error(err)
				return
			}
			r.Nil(err)
		})
	}
}

func TestJSON_Group(t *testing.T) {
	tests := []struct {
		name    string
		groups  []roller.Group
		id      string
		want    roller.Group
		wantErr bool
	}{
		{
			name:    "simple get",
			groups:  []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			id:      "5",
			want:    roller.Group{Name: "faz", ID: "5"},
			wantErr: false,
		}, {
			name:    "impossible get",
			groups:  []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			id:      "10",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups: tt.groups,
			}
			r := require.New(t)
			got, err := j.Group(tt.id)
			if tt.wantErr {
				r.NotNil(err)
				var e GroupNotFoundError
				r.True(errors.As(err, &e))
				r.Equal(tt.id, e.ID())
				r.Equal(e.Error(), fmt.Sprintf("group ID \"%s\" cant be found", tt.id))
				r.Zero(tt.want, "Expected should be zero when error is expected")
				return
			}
			r.Equal(tt.want, got)
			r.Nil(err)
		})
	}
}

func TestJSON_AddGroup(t *testing.T) {
	tests := []struct {
		name      string
		groups    []roller.Group
		argGroup  roller.Group
		wantState []roller.Group
		readOnly  bool
	}{
		{
			name:      "simple add",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "faz", ID: "5"}},
			argGroup:  roller.Group{Name: "far", ID: "4"},
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "faz", ID: "5"}, {Name: "far", ID: "4"}},
		}, {
			name:      "add update",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}},
			argGroup:  roller.Group{Name: "far", ID: "3"},
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "far", ID: "3"}},
		}, {
			name:      "readonly locked",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}},
			argGroup:  roller.Group{Name: "far", ID: "3"},
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}},
			readOnly:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups:   tt.groups,
				readOnly: tt.readOnly,
			}
			r := require.New(t)
			err := j.AddGroup(tt.argGroup)
			if tt.readOnly {
				r.Equal(tt.wantState, j.groups)
				r.Equal(tt.wantState, tt.groups, "Expected state and provided initial group should be same in readonly mode")
				r.NotNil(err)
				var e ReadOnlyError
				r.True(errors.As(err, &e))
				r.Equal("provider is set to readonly mode", e.Error())
				return
			}
			r.Nil(err)
			r.Equal(tt.wantState, j.groups)
		})
	}
}

func TestJSON_RemoveGroup(t *testing.T) {
	tests := []struct {
		name      string
		groups    []roller.Group
		id        string
		wantErr   bool
		wantState []roller.Group
		readOnly  bool
	}{
		{
			name:      "simple remove",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}},
			id:        "4",
			wantErr:   false,
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}},
		}, {
			name:      "impossible remove",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}},
			id:        "10",
			wantErr:   true,
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}},
		}, {
			name:      "glitched remove",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "2"}, {Name: "faz", ID: "2"}},
			id:        "2",
			wantErr:   false,
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "baz", ID: "2"}, {Name: "faz", ID: "2"}},
		}, {
			name:      "nil remove",
			groups:    []roller.Group{},
			id:        "2",
			wantErr:   true,
			wantState: []roller.Group{},
		}, {
			name:      "simple remove",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}},
			id:        "5",
			wantErr:   true,
			wantState: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}},
			readOnly:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups:   tt.groups,
				readOnly: tt.readOnly,
			}
			r := require.New(t)
			err := j.RemoveGroup(tt.id)
			r.Equal(tt.wantState, j.groups)

			if tt.wantErr {
				r.NotNil(err)
				if tt.readOnly {
					r.True(errors.Is(err, ReadOnlyError{}))
					return
				}
				var e GroupNotFoundError
				r.True(errors.As(err, &e))
				r.Equal(tt.id, e.ID())
				return
			}
			r.Nil(err)
		})
	}
}

func TestJSON_WalkGroup(t *testing.T) {
	tests := []struct {
		name       string
		groups     []roller.Group
		i          int
		wantGroups []roller.Group
	}{
		{
			name:       "short selection",
			groups:     []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			i:          4,
			wantGroups: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}},
		}, {
			name:       "exact selection",
			groups:     []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			i:          5,
			wantGroups: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
		}, {
			name:       "over selection",
			groups:     []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}},
			i:          5,
			wantGroups: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups: tt.groups,
			}
			i := tt.i
			r := require.New(t)
			rx := make(chan roller.Group, 10)
			gi := len(j.groups)
			err := j.WalkGroup(func(group roller.Group, last bool) (halt bool) {
				gi--
				i--
				rx <- group
				if gi <= 0 {
					r.Equal(0, gi)
					r.True(last)
				} else {
					r.False(last)
				}
				if i <= 0 {
					r.Equal(0, i)
					close(rx)
					return true
				}
				if last {
					close(rx)
				}
				return false
			})
			r.Nil(err)

			gs := make([]roller.Group, 0, len(j.groups)+2)
			for group := range rx {
				gs = append(gs, group)
			}
			r.Equal(tt.wantGroups, gs)
		})
	}
}

func TestJSON_duplicateCheck(t *testing.T) {
	tests := []struct {
		name      string
		groups    []roller.Group
		wantErr   bool
		original  roller.Group
		duplicate roller.Group
	}{
		{
			name:    "no errors",
			groups:  []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			wantErr: false,
		}, {
			name:      "simple duplicate",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "1"}, {Name: "far", ID: "3"}},
			wantErr:   true,
			original:  roller.Group{Name: "foo", ID: "1"},
			duplicate: roller.Group{Name: "baz", ID: "1"},
		}, {
			name:      "multiple duplicate",
			groups:    []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "1"}, {Name: "far", ID: "2"}, {Name: "faz", ID: "2"}},
			wantErr:   true,
			original:  roller.Group{Name: "foo", ID: "1"},
			duplicate: roller.Group{Name: "baz", ID: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{}
			r := require.New(t)
			err := j.duplicateCheck(tt.groups)
			if tt.wantErr {
				r.NotNil(err)
				var dup DuplicateGroupIDError
				r.True(errors.As(err, &dup))
				r.Equal(tt.original, dup.Original())
				r.Equal(tt.duplicate, dup.Duplicate())
			} else {
				r.Nil(err)
			}
		})
	}
	t.Run("Error Test", func(t *testing.T) {
		j := &JSON{}
		r := require.New(t)
		groups := []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "1"}}
		err := j.duplicateCheck(groups)
		r.NotNil(err)
		var dup DuplicateGroupIDError
		r.True(errors.As(err, &dup))
		r.Equal(roller.Group{Name: "foo", ID: "1"}, dup.Original())
		r.Equal(roller.Group{Name: "baz", ID: "1"}, dup.Duplicate())
		r.Equal("group ID not unique: ID \"1\"(foo[#]) already exist, "+
			"cant be shared with ID \"1\"(baz[#])", dup.Error())
	})
}

func TestJSON_findGroup(t *testing.T) {
	tests := []struct {
		name   string
		groups []roller.Group
		id     string
		wantI  int
		wantG  roller.Group
	}{
		{
			name:   "simple find",
			groups: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			id:     "5",
			wantI:  4,
			wantG:  roller.Group{Name: "faz", ID: "5"},
		}, {
			name:   "impossible find",
			groups: []roller.Group{{Name: "foo", ID: "1"}, {Name: "bar", ID: "2"}, {Name: "baz", ID: "3"}, {Name: "far", ID: "4"}, {Name: "faz", ID: "5"}},
			id:     "100",
			wantI:  -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				groups: tt.groups,
			}
			r := require.New(t)
			i, g := j.findGroup(tt.id)
			if tt.wantI <= -1 {
				r.Equal(tt.wantI, -1)
				r.Zero(g)
				r.Zero(tt.wantG, "Expected group should be zero when index is expected to be -1")
				return
			}
			r.Equal(tt.wantI, i)
			r.Equal(tt.wantG, g)
		})
	}
}
