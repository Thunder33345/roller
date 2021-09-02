package perms_manager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExplicitComparator_HasPermission(t *testing.T) {
	type args struct {
		permission []string
		target     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty",
			args: args{
				permission: nil,
				target:     "",
			},
			want: false,
		}, {
			name: "Fail",
			args: args{
				permission: []string{"foo.bar"},
				target:     "foo",
			},
			want: false,
		}, {
			name: "Fail 2",
			args: args{
				permission: []string{"foo.bar"},
				target:     "foo.",
			},
			want: false,
		}, {
			name: "Match",
			args: args{
				permission: []string{"foo.bar"},
				target:     "foo.bar",
			},
			want: true,
		}, {
			name: "Match 2",
			args: args{
				permission: []string{"foo.bar.baz"},
				target:     "foo.bar.baz",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := ExplicitComparator{}
			a := assert.New(t)

			got := j.HasPermission(List{Permission: tt.args.permission}, tt.args.target)
			if tt.want {
				a.True(got)
			} else {
				a.False(got)
			}
		})
	}
}

func TestExplicitComparator_HasPermissionWithLevel(t *testing.T) {
	type args struct {
		p     List
		node  string
		level int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Lower level",
			args: args{
				p: List{
					Level:      4,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: 5,
			},
			want: false,
		},
		{
			name: "Equal level",
			args: args{
				p: List{
					Level:      5,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: 5,
			},
			want: false,
		},
		{
			name: "Higher level",
			args: args{
				p: List{
					Level:      6,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: 5,
			},
			want: true,
		}, {
			name: "Higher with no perm",
			args: args{
				p: List{
					Level: 10,
				},
				node:  "foo.bar",
				level: 5,
			},
			want: false,
		}, {
			name: "Negative high level",
			args: args{
				p: List{
					Level:      -4,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: -5,
			},
			want: true,
		}, {
			name: "Negative equal level",
			args: args{
				p: List{
					Level:      -5,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: -5,
			},
			want: false,
		}, {
			name: "Negative low level",
			args: args{
				p: List{
					Level:      -20,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: -5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := ExplicitComparator{}
			a := assert.New(t)

			got := j.HasPermissionWithLevel(tt.args.p, tt.args.node, tt.args.level)
			if tt.want {
				a.True(got)
			} else {
				a.False(got)
			}
		})
	}
}

func TestExplicitComparator_IsHigherLevel(t *testing.T) {
	type args struct {
		source  List
		subject List
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Higher",
			args: args{
				source:  List{Level: 10},
				subject: List{Level: 1},
			},
			want: true,
		}, {
			name: "Equal",
			args: args{
				source:  List{Level: 10},
				subject: List{Level: 10},
			},
			want: false,
		}, {
			name: "Lower",
			args: args{
				source:  List{Level: 5},
				subject: List{Level: 5},
			},
			want: false,
		}, {
			name: "Negative Higher",
			args: args{
				source:  List{Level: -1},
				subject: List{Level: -10},
			},
			want: true,
		}, {
			name: "Negative Equal",
			args: args{
				source:  List{Level: -10},
				subject: List{Level: -10},
			},
			want: false,
		}, {
			name: "Negative low",
			args: args{
				source:  List{Level: -100},
				subject: List{Level: -10},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := ExplicitComparator{}
			a := assert.New(t)
			got := j.IsHigherLevel(tt.args.source, tt.args.subject)
			if tt.want {
				a.True(got)
			} else {
				a.False(got)
			}
		})
	}
}

func TestImplicitComparator_HasPermission(t *testing.T) {
	type fields struct {
		Deliminator       string
		Terminator        string
		IncludeTerminator bool
	}
	type args struct {
		permission []string
		target     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Empty",
			args: args{
				permission: nil,
				target:     "",
			},
			want: false,
		}, {
			name: "Fail",
			args: args{
				permission: []string{"foo.bar"},
				target:     "foo",
			},
			want: false,
		}, {
			name: "Fail 2",
			args: args{
				permission: []string{"foo.bar"},
				target:     "foo.",
			},
			want: false,
		}, {
			name: "Match",
			args: args{
				permission: []string{"foo.bar"},
				target:     "foo.bar",
			},
			want: true,
		}, {
			name: "Match 2",
			args: args{
				permission: []string{"foo.bar.baz"},
				target:     "foo.bar.baz",
			},
			want: true,
		}, {
			name: "child perm",
			fields: fields{
				Deliminator: ".",
				Terminator:  "*",
			},
			args: args{
				permission: []string{"foo*"},
				target:     "foo.bar.baz",
			},
			want: true,
		}, {
			name: "missing wildcard",
			fields: fields{
				Deliminator: ".",
				Terminator:  "*",
			},
			args: args{
				permission: []string{"foo"},
				target:     "foo.bar.baz",
			},
			want: false,
		}, {
			name: "implicit match",
			fields: fields{
				Deliminator: ".",
			},
			args: args{
				permission: []string{"foo"},
				target:     "foo.bar.baz",
			},
			want: true,
		}, {
			name: "no terminator unrecognized symbol",
			fields: fields{
				Deliminator: ".",
			},
			args: args{
				permission: []string{"foo*"},
				target:     "foo.bar.baz",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.Deliminator
			if len(d) == 0 {
				d = "."
			}
			j := ImplicitComparator{
				Deliminator:       d,
				Terminator:        tt.fields.Terminator,
				IncludeTerminator: tt.fields.IncludeTerminator,
			}
			a := assert.New(t)

			got := j.HasPermission(List{Permission: tt.args.permission}, tt.args.target)
			if tt.want {
				a.True(got)
			} else {
				a.False(got)
			}
		})
	}
}

func TestImplicitComparator_HasPermissionWithLevel(t *testing.T) {
	type args struct {
		p     List
		node  string
		level int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Lower level",
			args: args{
				p: List{
					Level:      4,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: 5,
			},
			want: false,
		},
		{
			name: "Equal level",
			args: args{
				p: List{
					Level:      5,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: 5,
			},
			want: false,
		},
		{
			name: "Higher level",
			args: args{
				p: List{
					Level:      6,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: 5,
			},
			want: true,
		}, {
			name: "Higher with no perm",
			args: args{
				p: List{
					Level: 10,
				},
				node:  "foo.bar",
				level: 5,
			},
			want: false,
		}, {
			name: "Negative high level",
			args: args{
				p: List{
					Level:      -4,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: -5,
			},
			want: true,
		}, {
			name: "Negative equal level",
			args: args{
				p: List{
					Level:      -5,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: -5,
			},
			want: false,
		}, {
			name: "Negative low level",
			args: args{
				p: List{
					Level:      -20,
					Permission: []string{"foo.bar"},
				},
				node:  "foo.bar",
				level: -5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := ImplicitComparator{}
			a := assert.New(t)

			got := j.HasPermissionWithLevel(tt.args.p, tt.args.node, tt.args.level)
			if tt.want {
				a.True(got)
			} else {
				a.False(got)
			}
		})
	}
}

func TestImplicitComparator_IsHigherLevel(t *testing.T) {
	type args struct {
		source  List
		subject List
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Higher",
			args: args{
				source:  List{Level: 10},
				subject: List{Level: 1},
			},
			want: true,
		}, {
			name: "Equal",
			args: args{
				source:  List{Level: 10},
				subject: List{Level: 10},
			},
			want: false,
		}, {
			name: "Lower",
			args: args{
				source:  List{Level: 5},
				subject: List{Level: 5},
			},
			want: false,
		}, {
			name: "Negative Higher",
			args: args{
				source:  List{Level: -1},
				subject: List{Level: -10},
			},
			want: true,
		}, {
			name: "Negative Equal",
			args: args{
				source:  List{Level: -10},
				subject: List{Level: -10},
			},
			want: false,
		}, {
			name: "Negative low",
			args: args{
				source:  List{Level: -100},
				subject: List{Level: -10},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := ImplicitComparator{}
			a := assert.New(t)
			got := j.IsHigherLevel(tt.args.source, tt.args.subject)
			if tt.want {
				a.True(got)
			} else {
				a.False(got)
			}
		})
	}
}

func TestImplicitComparator_generateVariant(t *testing.T) {
	type fields struct {
		Deliminator       string
		Terminator        string
		IncludeTerminator bool
	}
	tests := []struct {
		name   string
		fields fields
		str    string
		want   []string
	}{
		{
			name:   "Simple",
			fields: fields{Deliminator: ".", Terminator: "*", IncludeTerminator: false},
			str:    "foo.bar.baz",
			want:   []string{"foo*", "foo.bar*", "foo.bar.baz*", "foo.bar.baz"},
		}, {
			name:   "Simple Terminator",
			fields: fields{Deliminator: ".", Terminator: "*", IncludeTerminator: true},
			str:    "foo.bar.baz",
			want:   []string{"foo*", "foo.bar*", "foo.bar.baz*", "foo.bar.baz", "*"},
		}, {
			name:   "Double dot",
			fields: fields{Deliminator: ".", Terminator: "*", IncludeTerminator: false},
			str:    "foo..bar",
			want:   []string{"foo*", "foo.*", "foo..bar*", "foo..bar"},
		}, {
			name:   "trailing dot",
			fields: fields{Deliminator: ".", Terminator: "*", IncludeTerminator: false},
			str:    "foo.bar.",
			want:   []string{"foo*", "foo.bar*", "foo.bar.*", "foo.bar."},
		}, {
			name:   "wildcard",
			fields: fields{Deliminator: ".", Terminator: "*", IncludeTerminator: false},
			str:    "foo.*",
			want:   []string{"foo*", "foo.**", "foo.*"},
		}, {
			name:   "include blank term",
			fields: fields{Deliminator: ".", IncludeTerminator: true},
			str:    "foo.bar.baz",
			want:   []string{"", "foo", "foo.bar", "foo.bar.baz"},
		}, {
			name:   "no terminator",
			fields: fields{Deliminator: "."},
			str:    "foo.bar.baz",
			want:   []string{"foo", "foo.bar", "foo.bar.baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := ImplicitComparator{
				Deliminator:       tt.fields.Deliminator,
				Terminator:        tt.fields.Terminator,
				IncludeTerminator: tt.fields.IncludeTerminator,
			}
			a := assert.New(t)
			got := j.generateVariant(tt.str)
			a.ElementsMatch(tt.want, got)
		})
	}
}
