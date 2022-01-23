package errors

import (
	"errors"
	"fmt"
)

// customErrs provides slice of customErr
type customErrs struct {
	errSlice []CustomError
}

// AddErr added error with CustomError interface in errors slice
func (c *customErrs) AddErr(errorInterface CustomError) {
	c.errSlice = append(c.errSlice, errorInterface)
}

// GetErrs return slice of errors with CustomError interface
func (c *customErrs) GetErrs() []CustomError {
	return c.errSlice
}

// Error represent string value of customErrs
func (c *customErrs) Error() string {
	var errs []string
	for k := range c.errSlice {
		errs = append(errs, c.errSlice[k].Error())
	}
	return fmt.Sprintf("there are %d custom err in errSlice, errs: %v", len(c.errSlice), errs)
}

// IsEmpty return:
// true if errors slice is empty
// false if there are any errors in it
func (c *customErrs) IsEmpty() bool {
	return len(c.errSlice) == 0
}

// IsErrorExist checks if there is an error with the specified parameters in the errors slice
func (c *customErrs) IsErrorExist(target error) bool {
	for _, v := range c.errSlice {
		if errors.Is(v, target) {
			return true
		}
	}
	return false
}

// NewMultiply create Multiple Errors representation struct that allowed to use MultipleCustomErrs interface
func NewMultiply() MultipleCustomErrs {
	return &customErrs{}
}
