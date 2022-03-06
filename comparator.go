package roller

import "strings"

//Comparator evaluates permissions and perform comparisons
type Comparator interface {
	//HasPermission returns true if the List given has said permission node
	HasPermission(p List, node string) bool
	//HasPermissionWithLevel returns true if List has given permission,
	//and their level is at least N or higher(>=)
	HasPermissionWithLevel(p List, node string, level int) bool
	//IsHigherLevel compares if source List's level is more than subject List
	//returns true if source is higher than subject(>)
	IsHigherLevel(source List, subject List) bool
}

//Insures that ExplicitComparator is Comparator
var _ Comparator = (*ExplicitComparator)(nil)

//ExplicitComparator evaluates permissions explicitly
//it will only match if exact permission node is present
//list having ["foo.bar"] will not match with a node of "foo.bar.baz"
type ExplicitComparator struct {
}

//HasPermission checks if List contains the exact node
func (j ExplicitComparator) HasPermission(p List, node string) bool {
	for _, n := range p.Permission {
		if n == node {
			return true
		}
	}
	return false
}

//HasPermissionWithLevel checks if List has the exact node, and their level is at least N
func (j ExplicitComparator) HasPermissionWithLevel(p List, node string, level int) bool {
	return p.Level >= level && j.HasPermission(p, node)
}

//IsHigherLevel compares if source List's level is more than subject List
func (j ExplicitComparator) IsHigherLevel(source List, subject List) bool {
	return source.Level > subject.Level
}

//Insure ImplicitComparator is Comparator
var _ Comparator = (*ImplicitComparator)(nil)

//ImplicitComparator evaluates permissions implicitly
//it can allow fuzzy matching and treat parent node as valid matches for children nodes
//it needs to know what's the Deliminator to generate all possible variants to match against
//having a Terminator set allow it to differentiate between granting foo.bar(exact) or foo.bar*(exact and all child)
//this functions by checking if List has all possible parent combinations of node
//for example looking up "foo.bar.baz" will check if List has "foo*" or "foo.bar*" or "foo.bar.baz*" or "foo.bar.baz"
type ImplicitComparator struct {
	//Deliminator is the character(s) tha will be separating permission nodes
	//For example "foo.bar.baz" means . is the deliminator in this example
	//deliminators should be consistent as it's required for the system to guess parent nodes
	Deliminator string
	//Terminator is the character used to terminate a recursive permission grant
	//Permissions ending with the terminator will be valid matches for children nodes
	//If the terminator is *, granting "foo.bar*" will allow recursive match on all children nodes
	//granting "foo.bar" would not match "foo.bar.baz"
	//setting terminator to "" will make all grants recursive
	Terminator string
	//IncludeTerminator dictates if the terminal alone is a valid match
	//Allows granting Terminator to match everything
	IncludeTerminator bool
}

//HasPermission fuzzy checks if List contains the node
//Checking for "foo.bar" will result in variations of parent node to be generated and checked against
//If node is "foo.bar", it would check if list has foo*(grant recursively) or foo.bar* or foo.bar(non-recursive grant)
func (j ImplicitComparator) HasPermission(p List, node string) bool {
	v := j.generateVariant(node)
	for _, sv := range v {
		for _, n := range p.Permission {
			if n == sv {
				return true
			}
		}
	}
	return false
}

//HasPermissionWithLevel fuzzy checks if List contains the node, and their level is at least N
func (j ImplicitComparator) HasPermissionWithLevel(p List, node string, level int) bool {
	if p.Level <= level {
		return false
	}
	return j.HasPermission(p, node)
}

//IsHigherLevel compares if source List's level is more than subject List
func (j ImplicitComparator) IsHigherLevel(source List, subject List) bool {
	return source.Level > subject.Level
}

//generateVariant takes in a permission node
//and returns a list of possible parent permission nodes
//for example foo.bar.baz will return:
//foo*, foo.bar*, foo.bar.baz*, foo.bar.baz
//if ImplicitComparator.IncludeTerminator is true *
func (j ImplicitComparator) generateVariant(str string) []string {
	exp := strings.Split(str, j.Deliminator)
	//create the slice of possible parent nodes with 2 extra space
	o := make([]string, 0, len(exp)+2)

	if j.IncludeTerminator {
		o = append(o, j.Terminator)
	}
	o = append(o, str)
	for i := 0; i < len(exp); i++ {
		//skip the last element if terminator is nothing to prevent duplicating
		if j.Terminator == "" && i+1 == len(exp) {
			continue
		}
		o = append(o, strings.Join(exp[:i+1], j.Deliminator)+j.Terminator)
	}
	return o
}
