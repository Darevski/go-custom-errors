package customErrors

// customErr levels
const (
	Dataservice = 1
	UseCase     = 2
	Transport   = 3
)

// customErr types
const (
	InternalError    = 500
	NotFound         = 404
	InvalidArguments = 412
)

// MultipleCustomErrs for errSlice of custom errs
type MultipleCustomErrs interface {
	AddErr(errorInterface CustomError)
	IsEmpty() bool
	GetErrs() []CustomError
	Error() string
}

// CustomError for custom customErr type
type CustomError interface {
	// Error implement error interface support
	Error() string
	// GetType return type code of error, see customErr levels
	GetType() int
	// GetCode return code of error, see customErr types
	GetCode() int
	// GetMessage return error message value
	GetMessage() string
	// GetStack return stack of error with baggage on every level
	getStack(result *[]map[string]interface{})
	// GetErrPath return file path of error
	GetErrPath() string
	// GetFullTrace return stack by getStack with baggage on top level
	GetFullTrace() map[string]interface{}
	// AddBaggage add fields for error baggage
	AddBaggage(baggage map[string]interface{})
	// GetBaggage return error baggage
	GetBaggage() map[string]interface{}
}
