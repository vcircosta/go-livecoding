package checker

import "fmt"

type UnreachableError struct {
	URL string
	Err error
}

func (e *UnreachableError) Error() string {
	return fmt.Sprintf("URL unreachable : %s (%v)", e.URL, e.Err)
}

func (e *UnreachableError) Unwrap() error {
	return e.Err
}
