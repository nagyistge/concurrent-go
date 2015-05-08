package concurrent

import "fmt"

type CombinedError interface {
	error
	Errors() []error
}

func NewCombinedError(errs []error) CombinedError {
	return newCombinedError(errs)
}

type combinedError struct {
	errs []error
}

func newCombinedError(errs []error) *combinedError {
	return &combinedError{errs}
}

func (this *combinedError) Error() string {
	return fmt.Sprintf("%v", this.errs)
}

func (this *combinedError) Errors() []error {
	return this.errs
}
