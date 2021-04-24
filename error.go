package errors

import (
	"fmt"
)

type ErrorCode int
type ErrorPath string
type ErrorMessage string
type ErrorDataLevel int
type ErrorSeverity int
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
	// It is used to indicate the severity of errors and can be used
	// For example, to determine whether a given error should be recorded in the debug log
	severity ErrorSeverity
}

// StackError struct used for StackTrace output
type StackError struct {
	Message   ErrorMessage
	Baggage   ErrorBaggage
	DataLayer ErrorDataLevel
	Code      ErrorCode
}

// GetDataLevel returns the error level based on the data level at which the error occurred
func (e *customErr) GetDataLevel() ErrorDataLevel {
	return e.dataLayer
}

// GetCode returns a code of customErr
func (e *customErr) GetCode() ErrorCode {
	return e.code
}

// GetSeverity return severity value
func (e *customErr) GetSeverity() ErrorSeverity {
	return e.severity
}

// SetSeverity set severity value
func (e *customErr) SetSeverity(severity ErrorSeverity) CustomError {
	e.severity = severity
	return e
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
func (e *customErr) AddOperation(message ErrorMessage, baggage ErrorBaggage, severity ErrorSeverity) CustomError {
	return New(e.code, message, e.dataLayer, DetectPath(skipPackage), baggage, severity, e)
}

// AddBaggage add fields with values to err
func (e *customErr) AddBaggage(baggage ErrorBaggage) CustomError {
	for k, v := range baggage {
		e.baggage[k] = v
	}
	return e
}

// SetCode set error code
func (e *customErr) SetCode(code ErrorCode) CustomError {
	e.code = code
	return e
}

// SetMessage set message of error
func (e *customErr) SetMessage(message ErrorMessage) CustomError {
	e.message = message
	return e
}

// SetPath set error path
func (e *customErr) SetPath(errPath ErrorPath) CustomError {
	e.path = errPath
	return e
}

// SetDataLevel set data layer level
func (e *customErr) SetDataLevel(dataLayer ErrorDataLevel) CustomError {
	e.dataLayer = dataLayer
	return e
}

// SetBaggage set baggage of error
func (e *customErr) SetBaggage(baggage ErrorBaggage) CustomError {
	// do not allow to use nil as storage
	if baggage == nil {
		e.baggage = make(ErrorBaggage)
		return nil
	}
	e.baggage = baggage
	return e
}

// Unwrap return wrapped error with standard error interface
func (e *customErr) Unwrap() error {
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
