package roller

import (
	"fmt"
)

//providerGroupError is an internal error used by BasicProcessor
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
