package errors

import (
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
	var errors []string
	for k := range c.errSlice {
		errors = append(errors, c.errSlice[k].Error())
	}
	return fmt.Sprintf("there are %d custom err in errSlice, errors: %v", len(c.errSlice), errors)
}

// IsEmpty return:
// true if errors slice is empty
// false if there are any errors in it
func (c *customErrs) IsEmpty() bool {
	return len(c.errSlice) == 0
}

// isErrorExist checks if there is an error with the specified parameters in the errors slice
func (c *customErrs) isErrorExist(code ErrorCode, level ErrorDataLevel) bool {
	for _, v := range c.errSlice {
		if v.GetDataLevel() == level && v.GetCode() == code {
			return true
		}
	}
	return false
}
