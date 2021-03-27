package errors

// CustomErr Codes
const (
	InternalError    ErrorCode = 1500
	InvalidArguments ErrorCode = 1412
	NotFound         ErrorCode = 1404
	AccessDenied     ErrorCode = 1403
	Unauthorized     ErrorCode = 1401
)

// CustomErr DataLayers
// Based on "clean architecture"
const (
	DataService ErrorDataLevel = 101
	UseCase     ErrorDataLevel = 102
	Container   ErrorDataLevel = 103
	Controller  ErrorDataLevel = 104
	Transport   ErrorDataLevel = 105
)

// MultipleCustomErrs for errSlice of custom errs
type MultipleCustomErrs interface {
	// Adds an error to the errors storage
	AddErr(errorInterface CustomError)
	// Checks if the errors storage is empty
	IsEmpty() bool
	// Return slice of added errors
	GetErrs() []CustomError
	// Returns an error in the string representation
	Error() string
	// Checks if the error with specified attributes exist in storage
	isErrorExist(code ErrorCode, level ErrorDataLevel) bool
}

// CustomError for custom customErr type
type CustomError interface {
	// GetNative return original error that has been wrapped
	GetNative() error
	// Error implement error interface support
	Error() string
	// GetType return type code of error, see customErr levels
	GetDataLevel() ErrorDataLevel
	// GetCode return code of error, see customErr types
	GetCode() ErrorCode
	// GetMessage return error message value
	GetMessage() ErrorMessage
	// GetErrPath return file path of error
	GetPath() ErrorPath
	// GetBaggage return error baggage
	GetBaggage() ErrorBaggage
	// GetFullTrace return StackError slice
	GetFullTraceSlice() (result []StackError)
	// IsErrorExistInStack checks if there is an error with the specified parameters in the error stack
	IsErrorExistInStack(code ErrorCode, level ErrorDataLevel) bool
	// IsErrorWithCodeExistInStack checks if there is an error with specified Code exist in error stack
	IsErrorWithCodeExistInStack(code ErrorCode) bool
	// AddBaggage add fields for error baggage
	AddBaggage(baggage ErrorBaggage)
	// SetCode set error code
	SetCode(code ErrorCode)
	// SetMessage set message of error
	SetMessage(message ErrorMessage)
	// SetPath set error path
	SetPath(errPath ErrorPath)
	// SetDataLevel set data layer level
	SetDataLevel(dataLayer ErrorDataLevel)
	// SetBaggage set baggage of error - fully rewrite exist baggage
	SetBaggage(baggage ErrorBaggage)
	// AddOperation is a simplified version of the New function that allows to override an error by adding only a message
	// and baggage the error code and level are taken from the original error, ErrPath will be automatically written with
	// the place where the AddOperation function is called
	AddOperation(message ErrorMessage, baggage ErrorBaggage) CustomError
	// GetStack return StackError of error with baggage on every level
	getStack(result *[]StackError)
}

// New create new custom error with provided params
// also used for error wrapping with the possibility of creating and later getting an error stack
func New(
	code ErrorCode, message ErrorMessage, errDataLevel ErrorDataLevel, errPath ErrorPath, baggage ErrorBaggage,
	err error,
) CustomError {
	// Initialize empty ErrorBaggage map to prevent panics
	if baggage == nil {
		baggage = make(ErrorBaggage)
	}

	return &customErr{
		code:      code,
		message:   message,
		path:      errPath,
		nativeErr: err,
		dataLayer: errDataLevel,
		baggage:   baggage,
	}
}

// NewMultiply create Multiple Errors representation struct that allowed to use MultipleCustomErrs interface
func NewMultiply() MultipleCustomErrs {
	return &customErrs{}
}

// Wrap is a simplified version of the New function that will create custom error with empty error additional data
func Wrap(err error, message ErrorMessage) CustomError {
	return &customErr{
		path: DetectPath(skipPackage),
		message:   message,
		baggage:   make(ErrorBaggage),
		nativeErr: err,
	}
}
