package ranker

import "strings"

var _ Judge = (*ExplicitJudge)(nil)

type ExplicitJudge struct {
}

func (j ExplicitJudge) HasPermission(p Permissible, node string) bool {
	for _, n := range p.Permission {
		if n == node {
			return true
		}
	}
	return false
}

func (j ExplicitJudge) HasPermissionWithLevel(p Permissible, node string, level int) bool {
	if p.Level > level {
		return false
	}
	return j.HasPermission(p, node)
}

func (j ExplicitJudge) IsHigherLevel(source Permissible, subject Permissible) bool {
	return source.Level > subject.Level
}

var _ Judge = (*ImplicitJudge)(nil)

type ImplicitJudge struct {
	Deliminator string
	Terminator  string
}

func (j ImplicitJudge) HasPermission(p Permissible, node string) bool {
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

func (j ImplicitJudge) HasPermissionWithLevel(p Permissible, node string, level int) bool {
	if p.Level > level {
		return false
	}
	return j.HasPermission(p, node)
}

func (j ImplicitJudge) IsHigherLevel(source Permissible, subject Permissible) bool {
	return source.Level > subject.Level
}

func (j ImplicitJudge) generateVariant(exp []string) []string {
	var o []string

	for i := 0; i < len(exp); i++ {
		o = append(o, strings.Join(exp[:i], j.Deliminator)+j.Terminator)
	}
	o = append(o, strings.Join(exp, j.Deliminator))
	return o
}
