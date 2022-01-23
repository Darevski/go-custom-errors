package errors

import (
	errs "errors"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	referenceErrType       = NotFound
	referenceLevel         = ControllerLevel
	referenceErrFormat     = "It is %s based on %s case"
	referenceErrFormatArgs = []interface{}{"test err", "test case"}
	referenceErrorText     = fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs...)
	referenceSeverity      = Debug
)

var referenceBaggage = ErrorBaggage{
	"key1": "value1",
	"key2": []string{"value1_1", "value1_2"},
}

var errNativeReference = errors.New(referenceErrorText)

func referenceError() *customErr {
	return &customErr{
		errType:    referenceErrType,
		baggage:    referenceBaggage,
		level:      referenceLevel,
		severity:   referenceSeverity,
		wrappedErr: errNativeReference,
	}
}

func Test_customErr_Baggage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	emptyCheckedErr := referenceError()

	assertions.Equal(referenceBaggage, checkedErr.GetBaggage(), "Check getting baggage")

	_ = checkedErr.SetBaggage(nil)
	assertions.Equal(make(ErrorBaggage), checkedErr.GetBaggage(), "Check nil baggage passing")

	expectedBaggage := ErrorBaggage{
		"key1": "value2",
		"key3": []int{1, 2, 3, 4},
	}
	assertions.Equal(expectedBaggage, emptyCheckedErr.SetBaggage(expectedBaggage).GetBaggage(), "Check setting Baggage")

	err := checkedErr.AddBaggage(expectedBaggage)
	assertions.Equal(err, checkedErr, "Check that AddBaggage realize calling chain pattern")
	assertions.Equal(expectedBaggage, checkedErr.GetBaggage(), "Check adding baggage")

	err = checkedErr.AddBaggage(ErrorBaggage{
		"key2": "value3",
	})
	assertions.Equal(err, checkedErr, "Check that AddBaggage realize calling chain pattern")
	expectedBaggage["key2"] = "value3"

	assertions.Equal(expectedBaggage, checkedErr.GetBaggage(), "Check setting baggage")

}

func Test_customErr_Severity(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()

	assertions.Equal(referenceSeverity, checkedErr.GetSeverity(), "Check getting severity")

	err := checkedErr.SetSeverity(Critical)
	assertions.Equal(err, checkedErr, "Check that SetSeverity realize calling chain pattern")
	assertions.Equal(Critical, checkedErr.GetSeverity(), "Check setting severity")
}

func Test_customErr_Level(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()

	assertions.Equal(referenceLevel, checkedErr.GetLevel(), "Check getting level")

	err := checkedErr.SetLevel(ContainerLevel)
	assertions.Equal(err, checkedErr, "Check that GetLevel realize calling chain pattern")
	assertions.Equal(ContainerLevel, checkedErr.GetLevel(), "Check setting level")
}

func Test_customErr_GetType(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	assertions.Equal(referenceErrType, checkedErr.GetType(), "Check getting type")
}

func Test_customErr_GetMessage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	assertions.Equal(referenceErrorText, checkedErr.GetMessage().String(), "Check getting message")
}

func Test_customErr_GetErrPath(t *testing.T) {
	assertions := assert.New(t)
	// Create new Error using NewBaseF and then get line of that call
	createdErr := NewBaseF(referenceErrFormat, referenceErrFormatArgs...)
	_, fn, line, _ := runtime.Caller(0)
	// Check that we have line of error creating in ErrorPath
	assertions.Contains(createdErr.GetPath().String(), fmt.Sprintf("%s:%d", fn, line-1), "Check getting err path")
}

func Test_customErr_Error(t *testing.T) {
	assertions := assert.New(t)
	err := errs.New(referenceErrorText)
	firstLevel := "Wrapped Level 1"
	secondLevel := "Wrapped Level 2"

	wrappedErr := Wrap(err, ErrorMessage(firstLevel))
	wrappedErr = Wrap(wrappedErr, ErrorMessage(secondLevel))
	assertions.Equal(wrappedErr.Error(), fmt.Sprintf("%s: %s: %s", secondLevel, firstLevel, referenceErrorText),
		"Check that Error() returns all errors messages in chain")
}

func Test_customErr_GetTraceSlice(t *testing.T) {
	assertions := assert.New(t)
	firstLevel := "Wrapped Level 1"
	secondLevel := "Wrapped Level 2"
	_, fn, line, _ := runtime.Caller(0)
	err := NewBaseF(referenceErrFormat, referenceErrFormatArgs...)
	wrappedErr := Wrap(err, ErrorMessage(firstLevel))
	wrappedErr = Wrap(wrappedErr, ErrorMessage(secondLevel))

	result := wrappedErr.GetTraceSlice()
	resultReference := []string{
		fmt.Sprintf("Message:\\s%s,\\sPath:\\s.+\\s.+%s:%d", secondLevel, fn, line+3),
		fmt.Sprintf("Message:\\s%s,\\sPath:\\s.+\\s.+%s:%d", firstLevel, fn, line+2),
		fmt.Sprintf("Message:\\s%s,\\sPath:\\s.+\\s.+%s:%d", referenceErrorText, fn, line+1),
		fmt.Sprintf("Cause:\\s%s[\\s\\S]+%s:%d", referenceErrorText, fn, line+1),
	}

	assertions.Equal(len(result), len(resultReference), "check that trace slice have correct elements count in it")
	for k, v := range resultReference {
		re := regexp.MustCompile(v)
		assertions.True(re.MatchString(result[k]), "check that %d trace element - %s, is equal with reference regexp %s", k+1, result[k], v)
	}

}

func Test_customErr_Unwrap(t *testing.T) {
	assertions := assert.New(t)
	err := errs.New(referenceErrorText)
	wrappedErr := Wrap(err, "Test")
	assertions.Equal(err, wrappedErr.Unwrap(), "Check getting wrapped error (pkg) one level")
	wrapped2Err := Wrap(wrappedErr, "Test Level 2")
	assertions.Equal(wrappedErr, wrapped2Err.Unwrap(), "Check getting wrapped error (customErr) two level")

	_, err = regexp.Compile("")
	wrappedErr = Wrap(err, "Test 2")
	assertions.Equal(err, wrappedErr.Unwrap(), "Check getting wrapped error (native)")
}

func Test_customErr_Cause(t *testing.T) {
	assertions := assert.New(t)

	err := errs.New(referenceErrorText)
	wrappedErr := Wrap(err, "Test Level 1")
	wrappedErr = Wrap(wrappedErr, "Test Level 2")
	assertions.Equal(err, Cause(wrappedErr), "Check getting cause of Error")

	_, err = regexp.Compile("")
	wrappedErr = Wrap(err, "Test Level 1")
	assertions.Equal(err, wrappedErr.Unwrap(), "Check getting cause of Error (native)")
}

func Test_customErr_Is(t *testing.T) {
	assertions := assert.New(t)

	errNotFound := NotFound.NewBase(ErrorMessage(referenceErrorText))
	errNotFoundTwo := NotFound.NewBase("Not Found 2")
	errBadRequest := BadRequest.NewBase("Bad Request")

	assertions.False(errNotFound.Is(errBadRequest), "Check error type comparison")
	assertions.True(errNotFound.Is(errNotFoundTwo), "Check error type comparison")
}

func Test_customErr_IsMessageExistInStack(t *testing.T) {
	var ers []CustomError

	assertions := assert.New(t)
	firstLevel := ErrorMessage("Wrapped Level 1")
	secondLevel := ErrorMessage("Wrapped Level 2")

	// Create package error with two upper layers
	err := NewBase(ErrorMessage(referenceErrorText))
	wrappedErr := Wrap(err, firstLevel)
	wrappedErr = Wrap(wrappedErr, secondLevel)
	ers = append(ers, wrappedErr)

	// Create "errors" package error with two upper layers
	errNative := errNativeReference
	wrappedErr = Wrap(errNative, firstLevel)
	wrappedErr = Wrap(wrappedErr, secondLevel)
	ers = append(ers, wrappedErr)

	for k, v := range ers {
		assertions.True(v.IsMessageExistInStack(firstLevel), "Checking the stack for a message in it (%d)", k)
		assertions.True(v.IsMessageExistInStack(secondLevel), "Checking the stack for a message in it (%d)", k)
		assertions.True(v.IsMessageExistInStack(ErrorMessage(referenceErrorText)), "Checking the stack for a source message in it (%d)", k)
		assertions.False(v.IsMessageExistInStack("Test Message Not Exist in Stack"), "Checking the stack for a message that not exist in it (%d)", k)
	}

}

func Test_customErr_getStack(t *testing.T) {
	assertions := assert.New(t)

	var stack []CustomError
	err := NewBaseF(referenceErrFormat, referenceErrFormatArgs...)
	stack = append([]CustomError{err}, stack...)
	wrappedErr := Wrap(err, "Level 1")
	stack = append([]CustomError{wrappedErr}, stack...)
	wrappedErr = Wrap(wrappedErr, "Level 2")
	stack = append([]CustomError{wrappedErr}, stack...)

	resultStack := new([]CustomError)
	wrappedErr.getStack(resultStack)
	assertions.Equal(len(*resultStack), len(stack))

	for k, v := range stack {
		assertions.Equal(v, (*resultStack)[k])
	}
}

func TestNew(t *testing.T) {
	var ers []CustomError
	assertions := assert.New(t)
	ers = append(ers, New(referenceErrType, referenceLevel, referenceBaggage, referenceSeverity, ErrorMessage(referenceErrorText)))
	ers = append(ers, referenceErrType.New(referenceLevel, referenceBaggage, referenceSeverity, ErrorMessage(referenceErrorText)))

	for _, v := range ers {
		assertions.Equal(referenceErrType, v.GetType())
		assertions.Equal(referenceLevel, v.GetLevel())
		assertions.Equal(referenceBaggage, v.GetBaggage())
		assertions.Equal(referenceSeverity, v.GetSeverity())
		assertions.Equal(referenceErrorText, v.GetMessage().String())
	}
	err := New(referenceErrType, referenceLevel, nil, referenceSeverity, ErrorMessage(referenceErrorText))
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")

	err = referenceErrType.New(referenceLevel, nil, referenceSeverity, ErrorMessage(referenceErrorText))
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")
}

func TestNewF(t *testing.T) {
	var ers []CustomError
	assertions := assert.New(t)
	ers = append(ers, NewF(referenceErrType, referenceLevel, referenceBaggage, referenceSeverity, referenceErrFormat, referenceErrFormatArgs))
	ers = append(ers, referenceErrType.NewF(referenceLevel, referenceBaggage, referenceSeverity, referenceErrFormat, referenceErrFormatArgs))

	for k, v := range ers {
		assertions.Equal(referenceErrType, v.GetType())
		assertions.Equal(referenceLevel, v.GetLevel())
		assertions.Equal(referenceBaggage, v.GetBaggage())
		assertions.Equal(referenceSeverity, v.GetSeverity())
		assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), v.GetMessage().String(), "Check message formatting (%d)", k)

		v = NewF(referenceErrType, referenceLevel, nil, referenceSeverity, referenceErrFormat, referenceErrFormatArgs)
		assertions.Equal(ErrorBaggage{}, v.GetBaggage(), "Check nil baggage (%d)", k)
	}

	err := NewF(referenceErrType, referenceLevel, nil, referenceSeverity, referenceErrFormat, referenceErrFormatArgs)
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")

	err = referenceErrType.NewF(referenceLevel, nil, referenceSeverity, referenceErrFormat, referenceErrFormatArgs)
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")
}

func TestNewBase(t *testing.T) {
	assertions := assert.New(t)
	err := NewBase(ErrorMessage(referenceErrorText))

	assertions.Equal(DefaultType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(referenceErrorText, err.GetMessage().String(), "Check message")
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")
}

func TestNewTypeBase(t *testing.T) {
	assertions := assert.New(t)
	err := referenceErrType.NewBase(ErrorMessage(referenceErrorText))

	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(referenceErrorText, err.GetMessage().String(), "Check message")
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")
}

func TestNewBaseF(t *testing.T) {
	assertions := assert.New(t)
	err := NewBaseF(referenceErrFormat, referenceErrFormatArgs)

	assertions.Equal(DefaultType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), err.GetMessage().String(), "Check message formatting")
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")
}

func TestNewTypeBaseF(t *testing.T) {
	assertions := assert.New(t)
	err := referenceErrType.NewBaseF(referenceErrFormat, referenceErrFormatArgs)

	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), err.GetMessage().String(), "Check message formatting")
	assertions.Equal(ErrorBaggage{}, err.GetBaggage(), "Check nil baggage")
}

func TestTypeWrap(t *testing.T) {
	assertions := assert.New(t)

	err := referenceErrType.Wrap(errNativeReference, ErrorMessage(referenceErrorText))
	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(referenceErrorText, err.GetMessage().String())

	err = referenceErrType.Wrap(referenceError(), ErrorMessage(referenceErrorText))
	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(referenceLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(referenceSeverity, err.GetSeverity())
	assertions.Equal(referenceErrorText, err.GetMessage().String())
}

func TestWrap(t *testing.T) {
	assertions := assert.New(t)

	err := Wrap(errNativeReference, ErrorMessage(referenceErrorText))
	assertions.Equal(DefaultType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(referenceErrorText, err.GetMessage().String())

	err = Wrap(referenceError(), ErrorMessage(referenceErrorText))
	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(referenceLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(referenceSeverity, err.GetSeverity())
	assertions.Equal(referenceErrorText, err.GetMessage().String())
}

func TestTypeWrapF(t *testing.T) {
	assertions := assert.New(t)

	err := referenceErrType.WrapF(errNativeReference, referenceErrFormat, referenceErrFormatArgs)
	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), err.GetMessage().String())

	err = referenceErrType.WrapF(referenceError(), referenceErrFormat, referenceErrFormatArgs)
	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(referenceLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(referenceSeverity, err.GetSeverity())
	assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), err.GetMessage().String())
}

func TestWrapF(t *testing.T) {
	assertions := assert.New(t)

	err := WrapF(errNativeReference, referenceErrFormat, referenceErrFormatArgs)
	assertions.Equal(DefaultType, err.GetType())
	assertions.Equal(DefaultLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(DefaultSeverity, err.GetSeverity())
	assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), err.GetMessage().String(), "Check message formatting")

	err = WrapF(referenceError(), referenceErrFormat, referenceErrFormatArgs)
	assertions.Equal(referenceErrType, err.GetType())
	assertions.Equal(referenceLevel, err.GetLevel())
	assertions.Equal(ErrorBaggage{}, err.GetBaggage())
	assertions.Equal(referenceSeverity, err.GetSeverity())
	assertions.Equal(fmt.Sprintf(referenceErrFormat, referenceErrFormatArgs), err.GetMessage().String(), "Check message formatting")
}
