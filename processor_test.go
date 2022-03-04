package roller

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

type testGroup struct {
	Group
	Flags map[string]FlagEntry
}

type dummyProvider struct {
	groups []Group
}

func (d *dummyProvider) Flag(gid string, fid string) (FlagEntry, error) {
	panic("Flag() unsupported!")
}

func (d *dummyProvider) Group(uid string) (Group, error) {
	for _, v := range d.groups {
		if v.ID == uid {
			return v, nil
		}
	}
	return Group{}, NewNotFoundError(errors.New(fmt.Sprintf("group \"%s\" is not defined", uid)))
}

type dummyProviderWithFlag struct {
	groups []testGroup
}

func (d *dummyProviderWithFlag) Flag(gid string, fid string) (FlagEntry, error) {
	g, e := d.group(gid)
	if e != nil {
		return FlagEntry{}, e
	}
	v, ok := g.Flags[fid]
	if !ok {
		return FlagEntry{}, NewNotFoundError(errors.New(fmt.Sprintf("flag \"%s\" in group \"%s\" is not defined", fid, gid)))
	}
	return v, nil
}

func (d *dummyProviderWithFlag) Group(uid string) (Group, error) {
	g, e := d.group(uid)
	if e == nil {
		return g.Group, e
	}
	return Group{}, e
}

func (d *dummyProviderWithFlag) group(uid string) (testGroup, error) {
	for _, v := range d.groups {
		if v.ID == uid {
			return v, nil
		}
	}
	return testGroup{}, errors.New(fmt.Sprintf("group \"%s\" is not defined", uid))
}

func TestBasicProcessor_Process(t *testing.T) {
	type fields struct {
		Groups          []Group
		WeightAscending bool
	}
	tests := []struct {
		name    string
		fields  fields
		r       RawList
		want    List
		wantErr bool
	}{
		{
			name: "Simple 1",
			fields: fields{Groups: []Group{
				{
					ID: "1", Weight: 1000, Permission: Entry{
						Level:  10,
						Grant:  []string{"1.1", "1.2"},
						Revoke: []string{"3.3", "2.2", "o.2"},
					},
				},
				{
					ID: "2", Weight: 800, Permission: Entry{
						Level:    5,
						SetLevel: true,
						Grant:    []string{"2.1", "2.2"},
						Revoke:   []string{"3.2", "1.2"},
					},
				}, {
					ID: "3", Weight: 500, Permission: Entry{
						Level:  4,
						Grant:  []string{"3.1", "3.2", "3.3"},
						Revoke: []string{"2.1"},
					},
				},
			}},
			r: RawList{
				Overwrites: Entry{
					Level:  3,
					Grant:  []string{"o.1", "o.2"},
					Revoke: []string{"o.1"},
				},
				Groups: []string{"2", "3", "1"},
			},
			want: List{
				Level:      18,
				Permission: []string{"3.1", "2.1", "1.1", "1.2", "o.1", "o.2"},
			},
		}, {
			name: "Simple 2",
			fields: fields{Groups: []Group{
				{
					ID: "1", Weight: 2, Permission: Entry{
						Level:  1,
						Grant:  []string{"1", "1.2", "1.3", "self.revoke"},
						Revoke: []string{"2.4"},
					},
				}, {
					ID: "2", Weight: 1, Permission: Entry{
						Level:  2,
						Grant:  []string{"2", "2.2", "2.3", "2.4"},
						Revoke: []string{"1.3"},
					},
				},
			}},
			r: RawList{
				Overwrites: Entry{
					Level:  10,
					Grant:  []string{"self.grant", "self.order"},
					Revoke: []string{"self.revoke", "self.order", "self.404"},
				},
				Groups: []string{"1", "2"},
			},
			want: List{
				Level:      13,
				Permission: []string{"2", "2.2", "2.3", "1", "1.2", "1.3", "self.grant", "self.order"},
			},
			wantErr: false,
		}, {
			name: "3rd test",
			fields: fields{
				Groups: []Group{
					{
						ID:     "-1",
						Weight: -1,
						Permission: Entry{
							Level: 100,
							Grant: []string{"-1.test"},
						},
					}, {
						ID:     "0",
						Weight: 0,
						Permission: Entry{
							Level: -50,
							Grant: []string{"0.test"},
						},
					}, {
						ID:     "1",
						Weight: 2,
						Permission: Entry{
							EmptySet: true,
							SetLevel: true,
							Level:    -5,
							Grant:    []string{"1.1", "1.2"},
						},
					}, {
						ID:     "2",
						Weight: 3,
						Permission: Entry{
							Level:  2,
							Grant:  []string{"2.1", "2.2", "2.3"},
							Revoke: []string{"1.2"},
						},
					}, {
						ID:     "3",
						Weight: 4,
						Permission: Entry{
							Level:  -1,
							Grant:  []string{"3.1"},
							Revoke: []string{"2.3"},
						},
					},
				},
			},
			r: RawList{
				Overwrites: Entry{
					Level: 50,
					Grant: []string{"r.1"},
				},
				Groups: []string{"3", "1", "0", "2", "-1"},
			},
			want: List{
				Level:      46,
				Permission: []string{"1.1", "2.1", "2.2", "3.1", "r.1"},
			},
			wantErr: false,
		}, {
			name: "4th Reverse",
			fields: fields{
				WeightAscending: true, Groups: []Group{
					{
						ID:     "1",
						Weight: -1,
						Permission: Entry{
							Level:    5,
							SetLevel: true,
							Grant:    []string{"1.1", "1.2"},
							Revoke:   []string{"low.2"},
						},
					}, {
						ID:     "2",
						Weight: -2,
						Permission: Entry{
							Level:    2,
							SetLevel: false,
							Grant:    []string{"2.1", "2.2"},
							Revoke:   []string{"1.2"},
						},
					}, {
						ID:     "3",
						Weight: 10,
						Permission: Entry{
							Level:  -10,
							Grant:  []string{"low.1", "low.2"},
							Revoke: []string{"1.1"},
						},
					},
				},
			},
			r: RawList{
				Overwrites: Entry{
					Level: -1,
					Grant: []string{"self"},
				},
				Groups: []string{"3", "1", "2"},
			},
			want: List{
				Level:      6,
				Permission: []string{"low.1", "1.1", "2.1", "2.2", "self"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p := BasicProcessor{
				Provider:        &dummyProvider{groups: tt.fields.Groups},
				WeightAscending: tt.fields.WeightAscending,
			}
			got, err := p.Process(tt.r)
			if tt.wantErr {
				a.Error(err, "Error expected, missing error")
				return
			}
			a.NoError(err, "Unexpected error")
			a.Equal(tt.want.Level, got.Level, "Level should be equal")
			a.Equal(tt.want.Permission, got.Permission, "Permissions should be equal")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("Missing", func(t *testing.T) {
		a := assert.New(t)
		p := BasicProcessor{
			Provider: &dummyProvider{groups: []Group{{ID: "1"}}},
		}
		_, err := p.Process(RawList{Groups: []string{"1", "2"}})
		a.Error(err)
		me := MissingGroupError{}
		errors.As(err, &me)
		a.Equal("2", me.Group())
		a.Equal(fmt.Sprintf("group \"2\" is not defined"), me.Unwrap().Error())
		a.Equal("failed to access group \"2\": group \"2\" is not defined", me.Error())
	})
}

func TestBasicProcessor_ProcessFlags(t *testing.T) {
	type fields struct {
		Groups          []testGroup
		WeightAscending bool
	}
	tests := []struct {
		name    string
		fields  fields
		r       RawList
		flags   []string
		want    List
		wantErr bool
	}{
		{
			name: "Simple flags",
			fields: fields{Groups: []testGroup{
				{
					Group: Group{
						ID: "1", Weight: 100, Permission: Entry{
							Level:  10,
							Grant:  []string{"1.1", "1fs.1"},
							Revoke: []string{"2.2", "o.1"},
						},
					},
					Flags: map[string]FlagEntry{
						"f1": {
							Weight: 100,
							Entry: Entry{
								Level:  1,
								Grant:  []string{"1f1.1"},
								Revoke: []string{"1f2.1", "1fs.1"},
							},
						},
						"f2": {
							Weight: 20,
							Entry: Entry{
								Level:  1,
								Grant:  []string{"1f2.1"},
								Revoke: []string{"1f1.1"},
							},
						},
					},
				},
				{
					Group: Group{
						ID: "2", Weight: 50, Permission: Entry{
							Level:    5,
							SetLevel: true,
							Grant:    []string{"2.1", "2.2"},
							Revoke:   []string{"1.1"},
						},
					},
					Flags: map[string]FlagEntry{
						"f1": {
							Weight:     100,
							Preprocess: true,
							Entry: Entry{
								Grant:  []string{"2f1.1"},
								Revoke: []string{"2.1"},
							},
						},
						"f2": {
							Weight:     10,
							Preprocess: true,
							Entry: Entry{
								Grant:  []string{},
								Revoke: []string{"2f1.1"},
							},
						}},
				}, {
					Group: Group{
						ID:     "3",
						Weight: 10,
					},
					Flags: map[string]FlagEntry{
						"f3": {
							Weight:     100,
							Preprocess: true,
							Entry: Entry{
								Grant: []string{"3f3"},
							},
						}},
				},
			}},
			r: RawList{
				Overwrites: Entry{
					Level:  3,
					Grant:  []string{"o.1"},
					Revoke: []string{"o.1"},
				},
				Groups: []string{"1", "2", "3"},
			},
			flags: []string{"f1", "f2", "f3"},
			want: List{
				Level:      20,
				Permission: []string{"3f3", "2f1.1", "2.1", "1.1", "1f1.1", "o.1"},
			},
		}, {
			name: "Want error",
			r: RawList{
				Groups: []string{"1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p := BasicProcessor{
				Provider:        &dummyProviderWithFlag{groups: tt.fields.Groups},
				WeightAscending: tt.fields.WeightAscending,
			}
			got, err := p.ProcessFlags(tt.r, tt.flags...)
			if tt.wantErr {
				a.Error(err, "Error expected, missing error")
				return
			}
			a.NoError(err, "Unexpected error")
			a.Equal(tt.want.Level, got.Level, "Level should be equal")
			a.ElementsMatch(tt.want.Permission, got.Permission, "Permissions should match")
		})
	}
}

func TestBasicProcessor_MergeEntry(t *testing.T) {
	type args struct {
		l  List
		es []Entry
	}
	tests := []struct {
		name string
		args args
		want List
	}{
		{
			name: "Simple",
			args: args{
				l: List{
					Level:      5,
					Permission: []string{"o.1", "o.2"},
				},
				es: []Entry{
					{
						Level:  3,
						Grant:  []string{"1.1", "1.2"},
						Revoke: []string{"o.2", "o.3", "2.1"},
					}, {
						Level:  -1,
						Grant:  []string{"2.1"},
						Revoke: []string{"2.1"},
					},
				},
			},
			want: List{
				Level:      7,
				Permission: []string{"o.1", "1.1", "1.2", "2.1"},
			},
		}, {
			name: "Merge grant revoke",
			args: args{
				l: List{
					Permission: []string{"o.1"},
				},
				es: []Entry{
					{
						Grant:  []string{"2.1"},
						Revoke: []string{"2.1"},
					},
				},
			},
			want: List{
				Permission: []string{"o.1", "2.1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := BasicProcessor{}
			a := assert.New(t)
			got := p.MergeEntry(tt.args.l, tt.args.es...)
			a.Equal(tt.want.Level, got.Level)
			a.Equal(tt.want.Permission, got.Permission)
		})
	}
}
func TestBasicProcessor_processSet(t *testing.T) {
	p := BasicProcessor{}
	l := func() List {
		return List{
			Level:      10,
			Permission: []string{"foo", "bar"},
		}
	}

	tests := []struct {
		name      string
		base      List
		arg       Entry
		wantLvl   int
		wantPerms []string
	}{
		{
			name:      "Set Level",
			base:      l(),
			arg:       Entry{Level: 15, SetLevel: true},
			wantLvl:   15,
			wantPerms: l().Permission,
		}, {
			name:      "Add Level",
			base:      l(),
			arg:       Entry{Level: 15},
			wantLvl:   l().Level + 15,
			wantPerms: l().Permission,
		}, {
			name:      "Negative Level",
			base:      l(),
			arg:       Entry{Level: -100},
			wantLvl:   l().Level - 100,
			wantPerms: l().Permission,
		}, {
			name:      "Empty set",
			base:      l(),
			arg:       Entry{EmptySet: true},
			wantLvl:   l().Level,
			wantPerms: []string{},
		}, {
			name:      "Revoke",
			base:      l(),
			arg:       Entry{Revoke: []string{"foo"}},
			wantLvl:   l().Level,
			wantPerms: []string{"bar"},
		}, {
			name:      "Revoke All",
			base:      l(),
			arg:       Entry{Revoke: []string{"foo", "bar"}},
			wantLvl:   l().Level,
			wantPerms: []string(nil),
		}, {
			name:      "Grant",
			base:      l(),
			arg:       Entry{Grant: []string{"far"}},
			wantLvl:   l().Level,
			wantPerms: []string{"foo", "bar", "far"},
		}, {
			name:      "Revoke And Grant",
			base:      l(),
			arg:       Entry{Revoke: []string{"foo", "boo"}, Grant: []string{"far", "boo"}},
			wantLvl:   l().Level,
			wantPerms: []string{"bar", "far", "boo"},
		}, {
			name: "Mixed stuff",
			base: l(),
			arg: Entry{
				Level:  -3,
				Grant:  []string{"fao", "far"},
				Revoke: []string{"bar"},
			},
			wantLvl:   l().Level - 3,
			wantPerms: []string{"foo", "fao", "far"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			base := tc.base
			arg := tc.arg
			r := p.processSet(base, arg)
			a.Equal(tc.wantLvl, r.Level)
			if len(tc.wantPerms) == 0 {
				a.Empty(r.Permission)
			} else {
				a.Equal(tc.wantPerms, r.Permission)
			}
			a.Equal(base, tc.base)
			a.Equal(arg, tc.arg)
		})
	}

	t.Run("Dirty Test", func(t *testing.T) {
		a := assert.New(t)
		ls := List{
			Level:      10,
			Permission: []string{"foo", "bar"},
		}
		r := p.processSet(ls, Entry{
			Level:  2,
			Revoke: []string{"foo"},
		})
		a.Equal([]string{"bar"}, r.Permission)
		a.Equal(12, r.Level)
		a.Equal([]string{"foo", "bar"}, ls.Permission)
		a.Equal(10, ls.Level)
	})
}

func TestBasicProcessor_removeNodes(t *testing.T) {
	p := BasicProcessor{WeightAscending: false}
	tests := []struct {
		name string
		in   []string
		arg  []string
		want []string
	}{
		{
			name: "simple",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"2", "3"},
			want: []string{"1", "4"},
		}, {
			name: "repeated",
			in:   []string{"1", "2", "2", "2", "4"},
			arg:  []string{"2"},
			want: []string{"1", "4"},
		}, {
			name: "multiple repeated",
			in:   []string{"1", "2", "3", "4", "4"},
			arg:  []string{"2", "3", "3", "4"},
			want: []string{"1"},
		}, {
			name: "extract center",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"2", "3"},
			want: []string{"1", "4"},
		}, {
			name: "remove all",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"2", "3", "4", "1"},
			want: []string(nil),
		}, {
			name: "nil arg",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string(nil),
			want: []string{"1", "2", "3", "4"},
		}, {
			name: "none arg",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{},
			want: []string{"1", "2", "3", "4"},
		}, {
			name: "impossible",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"5", "6", "7", "8"},
			want: []string{"1", "2", "3", "4"},
		}, {
			name: "empty input",
			in:   []string{},
			arg:  []string{"1", "2"},
			want: []string(nil),
		}, {
			name: "nil input",
			in:   []string(nil),
			arg:  []string{"1", "2"},
			want: []string(nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			st := make([]string, len(tc.in))
			copy(st, tc.in)
			r := p.removeNodes(st, tc.arg)
			if len(tc.want) == 0 {
				a.Empty(r)
			} else {
				a.Equal(tc.want, r)
				a.Equal(tc.in, st)
			}
		})
	}
}

func TestBasicProcessor_compare(t *testing.T) {
	//notes: yes the ascending and expected results wants are inverted
	//this is due to the way the ordered results are iterated with range in a first->last order
	//therefore it should be read the other way
	tests := []struct {
		name            string
		WeightAscending bool
		arg             []int
		want            []int
	}{
		{
			name:            "Ascending",
			WeightAscending: true,
			arg:             []int{2, 4, 6, 3, 1, 5},
			want:            []int{6, 5, 4, 3, 2, 1},
		}, {
			name:            "Descending",
			WeightAscending: false,
			arg:             []int{2, 4, 6, 3, 1, 5},
			want:            []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p := BasicProcessor{
				WeightAscending: tt.WeightAscending,
			}
			s := tt.arg
			sort.Slice(s, func(i, j int) bool {
				return p.compare(s[i], s[j])
			})
			a.Equal(tt.want, s)
		})
	}
}

func BenchmarkBasicProcessor_removeNodes(b *testing.B) {
	p := BasicProcessor{}
	data := genBenchData(sliceArgs)
	b.ResetTimer()
	for _, c := range data {
		b.Run(c.args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.removeNodes(c.haystack, c.needles)
			}
		})
	}
}

func BenchmarkBasicProcessor_Process(b *testing.B) {
	i1 := randSlice(75, 25)
	i2 := randSlice(100, 25)
	i3 := randSlice(100, 25)
	i4 := randSlice(150, 25)
	p := BasicProcessor{
		Provider: &dummyProvider{groups: []Group{
			{
				ID:     "-1",
				Weight: -1,
				Permission: Entry{
					Level: 100,
					Grant: i1,
				},
			}, {
				ID:     "0",
				Weight: 0,
				Permission: Entry{
					Level:  -50,
					Grant:  i2,
					Revoke: randNeedles(i1, len(i1)/3),
				},
			}, {
				ID:     "1",
				Weight: 2,
				Permission: Entry{
					EmptySet: true,
					SetLevel: true,
					Level:    -5,
					Grant:    append([]string{"1.1", "1.2"}, i3...),
				},
			}, {
				ID:     "2",
				Weight: 3,
				Permission: Entry{
					Level:  2,
					Grant:  i4,
					Revoke: randNeedles(i3, len(i3)/3),
				},
			}, {
				ID:     "3",
				Weight: 4,
				Permission: Entry{
					Level:  -1,
					Grant:  []string{"3.1"},
					Revoke: randNeedles(i4, len(i4)/2),
				},
			},
		}},
	}
	r := RawList{
		Overwrites: Entry{
			Level: 50,
			Grant: []string{"r.1"},
		},
		Groups: []string{"3", "1", "0", "2", "-1"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Process(r)
	}
}

func BenchmarkBasicProcessor_ProcessMulti(b *testing.B) {
	var entry []Entry

	for i := 0; i < 10; i++ {
		entry = append(entry, genEntry(40, 150, entry))
	}

	group, ids := genGroups(entry)
	p := BasicProcessor{Provider: &dummyProvider{groups: group}}
	b.ResetTimer()
	b.Run("something", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = p.Process(RawList{Groups: ids})
		}
	})
}

func genGroups(entry []Entry) ([]Group, []string) {
	var gs []Group
	var ids []string

	for _, e := range entry {
		str := randStringBytes(5)
		ids = append(ids, str)
		gs = append(gs, Group{
			ID:         str,
			Weight:     rand.Intn(1000),
			Permission: e,
		})
	}
	return gs, ids
}
func genEntry(len int, c float64, prev []Entry) Entry {
	rb := func() bool {
		return rand.Intn(2) == 1
	}
	var rev []string

	for _, e := range prev {
		for _, g := range e.Grant {
			if float64(rand.Intn(1000)) <= c {
				rev = append(rev, g)
			}
		}
	}

	e := Entry{
		EmptySet: rb(),
		Level:    rand.Intn(1000),
		SetLevel: rb(),
		Grant:    randSlice(len, 10),
		Revoke:   rev,
	}
	return e
}
