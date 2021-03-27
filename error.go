package errors

import (
	"fmt"
)

type ErrorCode int
type ErrorPath string
type ErrorMessage string
type ErrorDataLevel int
type ErrorBaggage map[string]interface{}

// customErr provides custom err struct type
type customErr struct {
	// Describes an error in the form of a code, analogous to http error
	code ErrorCode
	// Message of Error
	message ErrorMessage
	// Map of interface{} values that related to error
	baggage ErrorBaggage
	// Path of file && error line
	path ErrorPath
	// Original/Wrapped error
	nativeErr error
	// Describes the error from the data level where the error occurred
	dataLayer ErrorDataLevel
}

// StackError struct used for StackTrace output
type StackError struct {
	Message   ErrorMessage
	Baggage   ErrorBaggage
	DataLayer ErrorDataLevel
	Code      ErrorCode
}

// GetErrorDataLevel returns the error level based on the data level at which the error occurred
func (e *customErr) GetDataLevel() ErrorDataLevel {
	return e.dataLayer
}

// GetCode returns a code of customErr
func (e *customErr) GetCode() ErrorCode {
	return e.code
}

// GetMessage returns a message of error with path && errorMessage
func (e *customErr) GetMessage() ErrorMessage {
	return ErrorMessage(fmt.Sprintf("%s, %s", e.GetPath(), e.message))
}

// GetBaggage return baggage of error
func (e *customErr) GetBaggage() ErrorBaggage {
	return e.baggage
}

// GetErrPath return path of error
func (e *customErr) GetPath() ErrorPath {
	return e.path
}

// Error represent string value of customErr
func (e *customErr) Error() string {
	return fmt.Sprintf("message: %s", e.GetMessage())
}

func (e *customErr) GetFullTraceSlice() (result []StackError) {
	e.getStack(&result)
	return result
}

// AddOperation is a simplified version of the New function that allows to override an error by adding only a message
// and baggage the error code and level are taken from the original error, ErrPath will be automatically written with
// the place where the AddOperation function is called
func (e *customErr) AddOperation(message ErrorMessage, baggage ErrorBaggage) CustomError {
	return New(e.code, message, e.dataLayer, DetectPath(skipPackage), baggage, e)
}

// AddBaggage add fields with values to err
func (e *customErr) AddBaggage(baggage ErrorBaggage) {
	for k, v := range baggage {
		e.baggage[k] = v
	}
}

// SetCode set error code
func (e *customErr) SetCode(code ErrorCode) {
	e.code = code
}

// SetMessage set message of error
func (e *customErr) SetMessage(message ErrorMessage) {
	e.message = message
}

// SetPath set error path
func (e *customErr) SetPath(errPath ErrorPath) {
	e.path = errPath
}

// SetDataLevel set data layer level
func (e *customErr) SetDataLevel(dataLayer ErrorDataLevel) {
	e.dataLayer = dataLayer
}

// SetBaggage set baggage of error
func (e *customErr) SetBaggage(baggage ErrorBaggage) {
	// do not allow to use nil as storage
	if baggage == nil {
		e.baggage = make(ErrorBaggage)
		return
	}
	e.baggage = baggage
}

// AddBaggage add fields with values to err
func (e *customErr) GetNative() error {
	return e.nativeErr
}

// IsErrorExistInStack checks if there is an error with the specified parameters in the error stack
func (e *customErr) IsErrorExistInStack(code ErrorCode, level ErrorDataLevel) bool {
	stack := make([]StackError, 0)
	e.getStack(&stack)

	for _, v := range stack {
		if v.DataLayer == level && v.Code == code {
			return true
		}
	}
	return false
}

// IsErrorWithCodeExistInStack checks if there is an error with specified Code exist in error stack
func (e *customErr) IsErrorWithCodeExistInStack(code ErrorCode) bool {
	stack := make([]StackError, 0)
	e.getStack(&stack)

	for _, v := range stack {
		if v.Code == code {
			return true
		}
	}
	return false
}

// getStack return slice of errors
func (e *customErr) getStack(result *[]StackError) {
	*result = append(*result, StackError{
		Message:   e.GetMessage(),
		Baggage:   e.GetBaggage(),
		DataLayer: e.GetDataLevel(),
		Code:      e.GetCode(),
	})

	if val, ok := e.nativeErr.(CustomError); ok {
		val.getStack(result)
	} else {
		*result = append(*result, StackError{
			Message: ErrorMessage(e.nativeErr.Error()),
		})
	}
}
