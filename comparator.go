package perms_manager

import "strings"

//Comparator is an interface for something that is able to compare permission
type Comparator interface {
	//HasPermission returns true if the List given has said permission node
	HasPermission(p List, node string) bool
	//HasPermissionWithLevel returns true if List has a permission, and also met the level requirement
	//comparison method using >= or > is up to implementor
	HasPermissionWithLevel(p List, node string, level int) bool
	//IsHigherLevel compares if source List is higher then subject List
	//returns true if source is higher
	IsHigherLevel(source List, subject List) bool
}

//SelfComparator an interface for something that's capable of comparing itself
//by holding it's own List and Comparator
//this is not used anywhere in the lib except to serve as an generic interface that can be used in other libraries
type SelfComparator interface {
	HasPermission(node string) bool
	HasPermissionWithLevel(node string, level int) bool
	IsHigherLevel(subject List) bool
}

//Insures that ExplicitComparator is Comparator
var _ Comparator = (*ExplicitComparator)(nil)

//ExplicitComparator is an explicit comparator
//it will only match if exact permission node is present
type ExplicitComparator struct {
}

//HasPermission checks if List has the exact node
func (j ExplicitComparator) HasPermission(p List, node string) bool {
	for _, n := range p.Permission {
		if n == node {
			return true
		}
	}
	return false
}

func (j ExplicitComparator) HasPermissionWithLevel(p List, node string, level int) bool {
	if level > p.Level {
		return false
	}
	return j.HasPermission(p, node)
}

func (j ExplicitComparator) IsHigherLevel(source List, subject List) bool {
	return source.Level > subject.Level
}

//Insure ImplicitComparator is Comparator
var _ Comparator = (*ImplicitComparator)(nil)

//ImplicitComparator is an implicit comparator
//it's able to implicitly compare having parents of a nodes
//it needs to know what's the Deliminator to generate all possible variants
//having a Terminator set allow it to differentiate between granting foo.bar(exact) or foo.bar*(exact and all childs)
type ImplicitComparator struct {
	//Deliminator is the character(s) tha will be separating permission nodes
	//For example foo.bar.baz means . is the deliminator in this example
	Deliminator string
	//Terminator is the character used to terminate a wildcard permission grant
	//For example terminator is *
	//foo.bar would only match foo.bar, not foo.bar.baz
	//but foo.bar* with terminator would match foo.bar.baz.buz
	//setting terminator to "" will allow for even more implicit granting
	Terminator string
}

//HasPermission checks if a list has a certain permission
//Checking for foo.bar will result in variations of parent node to be generated and checked against
//It would check if the list have foo*, foo.bar* or foo.bar
func (j ImplicitComparator) HasPermission(p List, node string) bool {
	v := j.generateVariant(node)
	for _, n := range p.Permission {
		for _, sv := range v {
			if n == sv {
				return true
			}
		}
	}
	return false
}

func (j ImplicitComparator) HasPermissionWithLevel(p List, node string, level int) bool {
	if level > p.Level {
		return false
	}
	return j.HasPermission(p, node)
}

func (j ImplicitComparator) IsHigherLevel(source List, subject List) bool {
	return source.Level > subject.Level
}

//generateVariant takes in a permission node
//and returns a list of possible parent permission nodes
//for example foo.bar.baz will return:
//foo*, foo.bar*, foo.bar.baz*, foo.bar.baz
func (j ImplicitComparator) generateVariant(str string) []string {
	exp := strings.Split(str, j.Deliminator)
	var o []string

	for i := 0; i < len(exp); i++ {
		o = append(o, strings.Join(exp[:i], j.Deliminator)+j.Terminator)
	}
	o = append(o, strings.Join(exp, j.Deliminator))
	return o
}
