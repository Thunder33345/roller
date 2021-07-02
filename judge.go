package ranker

import "strings"

type Comparator interface {
	HasPermission(p PermissionList, node string) bool
	HasPermissionWithLevel(p PermissionList, node string, level int) bool
	IsHigherLevel(source PermissionList, subject PermissionList) bool
}

var _ Comparator = (*ExplicitComparator)(nil)

type ExplicitComparator struct {
}

func (j ExplicitComparator) HasPermission(p PermissionList, node string) bool {
	for _, n := range p.Permission {
		if n == node {
			return true
		}
	}
	return false
}

func (j ExplicitComparator) HasPermissionWithLevel(p PermissionList, node string, level int) bool {
	if p.Level > level {
		return false
	}
	return j.HasPermission(p, node)
}

func (j ExplicitComparator) IsHigherLevel(source PermissionList, subject PermissionList) bool {
	return source.Level > subject.Level
}

var _ Comparator = (*ImplicitComparator)(nil)

type ImplicitComparator struct {
	Deliminator string
	Terminator  string
}

func (j ImplicitComparator) HasPermission(p PermissionList, node string) bool {
	v := j.generateVariant(strings.Split(node, j.Deliminator))
	for _, n := range p.Permission {
		for _, sv := range v {
			if n == sv {
				return true
			}
		}
	}
	return false
}

func (j ImplicitComparator) HasPermissionWithLevel(p PermissionList, node string, level int) bool {
	if p.Level > level {
		return false
	}
	return j.HasPermission(p, node)
}

func (j ImplicitComparator) IsHigherLevel(source PermissionList, subject PermissionList) bool {
	return source.Level > subject.Level
}

func (j ImplicitComparator) generateVariant(exp []string) []string {
	var o []string

	for i := 0; i < len(exp); i++ {
		o = append(o, strings.Join(exp[:i], j.Deliminator)+j.Terminator)
	}
	o = append(o, strings.Join(exp, j.Deliminator))
	return o
}
