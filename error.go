package customErrors

import (
	"fmt"
)

const StackMessage = "message"
const StackBaggage = "baggage"
const StackKey = "stack"
const BaggageKey = "baggage"

// customErr provides custom err struct type
type customErr struct {
	code      int
	message   string
	baggage   map[string]interface{}
	errType   int
	errPath   string
	nativeErr error
}

// GetType returns a type of customErr
func (e *customErr) GetType() int {
	return e.errType
}

// customErr represent string value of customErr
func (e *customErr) Error() string {
	return fmt.Sprintf("message: %s", e.GetMessage())
}

func (e *customErr) getStack(result *[]map[string]interface{}) {
	stack := map[string]interface{}{
		StackMessage: e.GetMessage(),
		StackBaggage: e.GetBaggage(),
	}

	*result = append(*result, stack)
	if val, ok := e.nativeErr.(CustomError); ok {
		val.getStack(result)
	}
}

func (e *customErr) GetFullTrace() map[string]interface{} {
	result := make(map[string]interface{})
	stack := make([]map[string]interface{}, 0)
	e.getStack(&stack)

	if len(stack) > 0 {
		result[StackKey] = stack
	}

	if val := e.GetBaggage(); len(val) > 0 {
		result[BaggageKey] = val
	}

	return result
}

// GetCode returns a code of customErr
func (e *customErr) GetCode() int {
	return e.code
}

// GetMessage returns a message of customErr
func (e *customErr) GetMessage() string {
	return fmt.Sprintf("%s, %s", e.errPath, e.message)
}

// NewErr create new custom customErr with provided args
// Save prev customErr, so it could output trace based on customErr method
func NewErr(
	code int, message string, errType int, errPath string, err error, baggage map[string]interface{},
) CustomError {
	if baggage == nil {
		baggage = make(map[string]interface{})
	}

	return &customErr{
		code:      code,
		message:   message,
		errPath:   errPath,
		nativeErr: err,
		errType:   errType,
		baggage:   baggage,
	}
}

// AddBaggage add fields with values to err
func (e *customErr) AddBaggage(baggage map[string]interface{}) {
	for k, v := range baggage {
		e.baggage[k] = v
	}
}

// AddBaggage get fields with values
func (e *customErr) GetBaggage() map[string]interface{} {
	return e.baggage
}

func (e *customErr) GetErrPath() string {
	return e.errPath
}

func NewErrs() MultipleCustomErrs {
	return &customErrs{}
}
