package errors

import errs "github.com/pkg/errors"

// ErrorLevel describes the data level of error, based on "clean architecture"
const (
	DefaultLevel = ErrorLevel(iota)
	DataLevel
	UseCaseLevel
	ContainerLevel
	ControllerLevel
	TransportLevel
)

// ErrorSeverity describes the severity of error
const (
	DefaultSeverity = ErrorSeverity(iota)
	Debug
	Info
	Warning
	Critical
	Fatal
	Panic
)

// ErrorType describes the type of error
const (
	DefaultType = ErrorType(iota)
	NotFound
	InvalidArguments

	InternalError
	BadRequest

	AccessDenied
	Unauthorized
)

//go:generate stringer -type=ErrorType,ErrorLevel,ErrorSeverity -output=const.go

type stackTracer interface {
	StackTrace() errs.StackTrace
}

// MultipleCustomErrs for errSlice of custom errs
type MultipleCustomErrs interface {
	// Error returns an error in the string representation
	Error() string
	// AddErr adds an error to the errors storage
	AddErr(errorInterface CustomError)
	// IsEmpty checks if the errors storage is empty
	IsEmpty() bool
	// GetErrs return slice of added errors
	GetErrs() []CustomError
	// IsErrorExist returns true if errs struct has err with target error type
	IsErrorExist(target error) bool
}

type Unwrapped interface {
	// Unwrap return original error that has been wrapped
	Unwrap() error
	// Error implement error interface support
	Error() string
}

// CustomError for custom customErr type
type CustomError interface {
	Unwrapped
	// GetLevel return type code of error, see customErr levels
	GetLevel() ErrorLevel
	// GetType return code of error, see customErr types
	GetType() ErrorType
	// GetMessage return error message value
	GetMessage() ErrorMessage
	// GetPath return file path of error
	GetPath() ErrorPath
	// GetSeverity return severity value
	GetSeverity() ErrorSeverity
	// SetSeverity set severity value
	SetSeverity(ErrorSeverity) CustomError
	// GetBaggage return error baggage
	GetBaggage() ErrorBaggage
	// GetTraceSlice return error stack trace in string slice
	// Include error message && error path, also cause error displayed with full stackTrace
	GetTraceSlice() (trace []string)
	// AddBaggage add fields for error baggage
	AddBaggage(baggage ErrorBaggage) CustomError
	// SetLevel set data level
	SetLevel(dataLayer ErrorLevel) CustomError
	// SetBaggage set baggage of error - fully rewrite exist baggage
	SetBaggage(baggage ErrorBaggage) CustomError
	// Is method is for errors.Is comparison supporting
	// Compare errors by ErrorType
	Is(target error) bool
	// getStack return slice of CustomError that represent error stacktrace
	getStack(result *[]CustomError)
	// IsMessageExistInStack represent is error stack contains error with specified message or not
	IsMessageExistInStack(message ErrorMessage) bool
}
