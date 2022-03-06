package roller

import (
	"fmt"
)

//providerGroupError is an error used by BasicProcessor when Provider cannot retrieve a certain group
type providerGroupError struct {
	group string
	cause error
}

func (e providerGroupError) Error() string {
	return fmt.Sprintf(`cant retrieve group "%s": %v`, e.group, e.cause)
}
func (e providerGroupError) Unwrap() error {
	return e.cause
}
