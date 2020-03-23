package errors

import "fmt"

type DomainStatusError struct {
	Err    error
	Domain string
}

func (e DomainStatusError) Error() string {
	return fmt.Sprintf("%s: %s", e.Domain, e.Err.Error())
}
