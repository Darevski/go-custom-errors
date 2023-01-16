package errors

import (
	"fmt"
	errs "github.com/pkg/errors"
	"regexp"
)

// New create custom error with the provided params && error message
func New(
	errType ErrorType, errLevel ErrorLevel, baggage ErrorBaggage,
	severity ErrorSeverity, message ErrorMessage,
) CustomError {
	// Initialize empty ErrorBaggage map to prevent panics
	if baggage == nil {
		baggage = make(ErrorBaggage)
	}
	return newCustomErr(errType, baggage, errLevel, severity, errs.New(message.String()))
}

// NewF create custom error with params && error message that formats according to a format specifier
func NewF(
	errType ErrorType, errLevel ErrorLevel, baggage ErrorBaggage,
	severity ErrorSeverity, format string, args ...interface{},
) CustomError {
	// Initialize empty ErrorBaggage map to prevent panics
	if baggage == nil {
		baggage = make(ErrorBaggage)
	}
	return newCustomErr(errType, baggage, errLevel, severity, errs.Errorf(format, args...))
}

// NewBase create custom error with specified message
// also all error attributes are set to default values
func NewBase(message ErrorMessage) CustomError {
	return newCustomErr(DefaultType, make(ErrorBaggage), DefaultLevel, DefaultSeverity, errs.New(message.String()))
}

// NewBaseF create custom error with error message that formats according to a format specifier
// also all error attributes are set to default values
func NewBaseF(format string, args ...interface{}) CustomError {
	return newCustomErr(DefaultType, make(ErrorBaggage), DefaultLevel, DefaultSeverity, errs.Errorf(format, args...))
}

// Wrap error with message. If wrapped error implements CustomError interface than all semantic data such
// a severity, error level, error type will be copied into result error
func Wrap(err error, message ErrorMessage) CustomError {
	wrappedErr := errs.Wrap(err, message.String())
	if customErr, ok := err.(CustomError); ok {
		return newCustomErr(customErr.GetType(), make(ErrorBaggage), customErr.GetLevel(), customErr.GetSeverity(), wrappedErr)
	}
	return newCustomErr(DefaultType, make(ErrorBaggage), DefaultLevel, DefaultSeverity, wrappedErr)
}

// WrapF is analogous to the Wrap method, except that instead of ErrorMessage there are formatting arguments for the message
func WrapF(err error, format string, args ...interface{}) CustomError {
	wrappedErr := errs.Wrapf(err, format, args...)
	if customErr, ok := err.(CustomError); ok {
		return newCustomErr(customErr.GetType(), make(ErrorBaggage), customErr.GetLevel(), customErr.GetSeverity(), wrappedErr)
	}
	return newCustomErr(DefaultType, make(ErrorBaggage), DefaultLevel, DefaultSeverity, wrappedErr)
}

const callerSkip = 1

type ErrorLevel int
type ErrorSeverity int
type ErrorBaggage map[string]interface{}

// customErr provides custom err struct type
type customErr struct {
	// Describes an error in the form of a code, analogous to http error
	errType ErrorType
	// Map of interface{} values that related to error
	baggage ErrorBaggage
	// Describes the error from the data level where the error occurred
	level ErrorLevel
	// It is used to indicate the severity of errors and can be used
	// For example, to determine whether a given error should be recorded in the debug log
	severity ErrorSeverity
	// Original/Wrapped error
	//wrappedErr Unwrapped
	wrappedErr error
}

type ErrorMessage string

func (m ErrorMessage) String() string {
	return string(m)
}

type ErrorPath string

func (p ErrorPath) String() string {
	return string(p)
}

// newCustomErr is constructor for customErr struct
func newCustomErr(errType ErrorType, baggage ErrorBaggage, dataLayer ErrorLevel, severity ErrorSeverity, originalErr error) *customErr {
	return &customErr{errType: errType, baggage: baggage, level: dataLayer, severity: severity, wrappedErr: originalErr}
}

// GetLevel returns the error level based on the data level at which the error occurred
func (e *customErr) GetLevel() ErrorLevel {
	return e.level
}

// SetLevel set data layer level
func (e *customErr) SetLevel(dataLayer ErrorLevel) CustomError {
	e.level = dataLayer
	return e
}

// GetType returns a code of customErr
func (e *customErr) GetType() ErrorType {
	return e.errType
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
	re := regexp.MustCompile(`(?U)(.+):`)
	errMessage := e.wrappedErr.Error()
	if result := re.FindStringSubmatch(errMessage); len(result) > 0 {
		return ErrorMessage(result[1])
	}
	return ErrorMessage(errMessage)
}

// GetBaggage return baggage of error
func (e *customErr) GetBaggage() ErrorBaggage {
	return e.baggage
}

// AddBaggage add fields with values to err
func (e *customErr) AddBaggage(baggage ErrorBaggage) CustomError {
	for k, v := range baggage {
		e.baggage[k] = v
	}
	return e
}

// SetBaggage set baggage of error
func (e *customErr) SetBaggage(baggage ErrorBaggage) CustomError {
	// do not allow using nil pointer as a storage
	if baggage == nil {
		e.baggage = make(ErrorBaggage)
		return e
	}
	e.baggage = baggage
	return e
}

// GetPath return path of error
func (e *customErr) GetPath() ErrorPath {
	if err, ok := e.wrappedErr.(stackTracer); ok {
		st := err.StackTrace()
		if len(st) == 0 {
			return ""
		}
		if len(st) <= callerSkip {
			return ErrorPath(fmt.Sprintf("%+s:%d", st[0], st[0]))
		} else {
			return ErrorPath(fmt.Sprintf("%+s:%d", st[callerSkip], st[callerSkip]))
		}
	}
	return ""
}

// Error represent string value of customErr
func (e *customErr) Error() string {
	return e.wrappedErr.Error()
}

func (e *customErr) GetTraceSlice() (trace []string) {
	stack := make([]CustomError, 0)
	e.getStack(&stack)
	for _, v := range stack {
		trace = append(trace, fmt.Sprintf("Message: %s, Path: %s", v.GetMessage().String(), v.GetPath().String()))
	}
	trace = append(trace, fmt.Sprintf("Cause: %+v", Cause(e)))
	return trace
}

// Unwrap return wrapped error with standard error interface
func (e *customErr) Unwrap() error {
	if pkgErr, ok := e.wrappedErr.(Unwrapped); ok {
		// Unwrap wrapped Err and get StackErr
		if pkgWithStack, ok := pkgErr.Unwrap().(Unwrapped); ok {
			// Unwrap pkg StackErr
			return pkgWithStack.Unwrap()
		}
	}
	return e.wrappedErr
}

func Cause(e CustomError) error {
	if val, ok := e.Unwrap().(CustomError); ok {
		return Cause(val)
	} else {
		return e.Unwrap()
	}
}

// Is - is the function that is used to compare errors by ErrorType
func (e *customErr) Is(target error) bool {
	if err, ok := target.(CustomError); ok && err.GetType() == e.errType {
		return true
	}
	return false
}

// IsMessageExistInStack checks if there is an error with the specified parameters in the error stack
func (e *customErr) IsMessageExistInStack(message ErrorMessage) bool {
	stack := make([]CustomError, 0)
	e.getStack(&stack)
	for k, v := range stack {
		if v.GetMessage() == message {
			return true
		} else if k == (len(stack)-1) && v.Unwrap().Error() == message.String() {
			return true
		}
	}
	return false
}

// getStack return slice of errors
func (e *customErr) getStack(result *[]CustomError) {
	*result = append(*result, e)
	if val, ok := e.Unwrap().(CustomError); ok {
		val.getStack(result)
	}
}
