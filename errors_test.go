package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	referenceErrCode   = NotFound
	referenceMessage   = "Not Found"
	referenceErrLayer  = UseCase
	referenceErrPath   = "Err Path"
	referenceErrorText = "test Error"
	referenceSeverity  = Debug
)

var referenceBaggage = ErrorBaggage{
	"key1": "value1",
	"key2": []string{"value1_1", "value1_2"},
}

var errNativeReference = errors.New(referenceErrorText)

func referenceError() customErr {
	return customErr{
		code:      referenceErrCode,
		message:   referenceMessage,
		baggage:   referenceBaggage,
		path:      referenceErrPath,
		nativeErr: errors.New(referenceErrorText),
		dataLayer: referenceErrLayer,
		severity:  referenceSeverity,
	}
}

func TestNewErr(t *testing.T) {
	assertions := assert.New(t)
	reference := referenceError()
	createdError := New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference)
	assertions.Equal(&reference, createdError, "Check creating mew error")

	reference.baggage = make(ErrorBaggage)
	createdError = New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, nil, referenceSeverity,
		errNativeReference)
	assertions.Equal(&reference, createdError, "Check creating new error with empty baggage")
}

func TestNewErrs(t *testing.T) {
	assertions := assert.New(t)
	assertions.Equal(&customErrs{}, NewMultiply())
}

func Test_customErr_AddBaggage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.baggage = make(ErrorBaggage)
	expectedBaggage := ErrorBaggage{
		"key1": "value2",
		"key3": []int{1, 2, 3, 4},
	}
	checkedErr.AddBaggage(expectedBaggage)
	assertions.Equal(expectedBaggage, checkedErr.GetBaggage())

	checkedErr.AddBaggage(ErrorBaggage{
		"key2": "value3",
	})

	expectedBaggage["key2"] = "value3"
	assertions.Equal(expectedBaggage, checkedErr.GetBaggage())
}

func Test_customErr_SetBaggage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.baggage = make(ErrorBaggage)
	baseBaggage := ErrorBaggage{
		"key1": "value2",
		"key3": []int{1, 2, 3, 4},
	}
	expectedBaggage := ErrorBaggage{
		"key2": "value3",
	}
	checkedErr.AddBaggage(baseBaggage)
	checkedErr.SetBaggage(expectedBaggage)
	assertions.Equal(expectedBaggage, checkedErr.GetBaggage())
}

func Test_customErr_SetSeverity(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.SetSeverity(Debug)
	assertions.Equal(Debug, checkedErr.GetSeverity())
}

func Test_customErr_GetSeverity(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.severity = Debug
	assertions.Equal(Debug, checkedErr.GetSeverity())
}

func Test_customErr_SetCode(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.SetCode(Unauthorized)
	assertions.Equal(Unauthorized, checkedErr.GetCode())
}

func Test_customErr_SetMessage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.SetMessage(ErrorMessage("Test"))
	assertions.Equal(ErrorMessage(referenceErrPath+", Test"), checkedErr.GetMessage())
}

func Test_customErr_SetPath(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.SetPath(ErrorPath("Test Path"))
	assertions.Equal(ErrorPath("Test Path"), checkedErr.GetPath())
}

func Test_customErr_SetDataLevel(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.SetDataLevel(Controller)
	assertions.Equal(Controller, checkedErr.GetDataLevel())
}

func Test_customErr_Wrap(t *testing.T) {
	assertions := assert.New(t)
	err := errors.New("test error")
	checkedErr := Wrap(err, "WrappedError")
	checkedErr.SetPath("Test Path")
	assertions.Equal(ErrorMessage("Test Path, WrappedError"), checkedErr.GetMessage())
	assertions.Equal(err, checkedErr.Unwrap())
}

func Test_customErr_GetBaggage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.baggage = make(ErrorBaggage)
	expectedBaggage := ErrorBaggage{
		"key1": "value2",
		"key3": []int{1, 2, 3, 4},
	}
	checkedErr.baggage = expectedBaggage
	assertions.Equal(expectedBaggage, checkedErr.GetBaggage())
}

func Test_customErr_GetCode(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.code = InternalError
	assertions.Equal(InternalError, checkedErr.GetCode())
}

func Test_customErr_GetErrPath(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	errPath := ErrorPath("Err Path")
	checkedErr.path = errPath
	assertions.Equal(errPath, checkedErr.GetPath())
}

func Test_customErr_GetErrorDataLevel(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.dataLayer = Controller
	assertions.Equal(Controller, checkedErr.GetDataLevel())
}

func Test_customErr_GetFullTraceErrorsSlice(t *testing.T) {
	assertions := assert.New(t)
	requires := require.New(t)
	topErr := referenceError()
	topErr.message = "Top Error"
	middleErr := referenceError()
	middleErr.message = "Middle Error"
	lastErr := referenceError()
	lastErr.message = "Last Error"
	middleErr.nativeErr = &lastErr
	topErr.nativeErr = &middleErr

	referenceStack := []StackError{
		{
			Message:   referenceErrPath + ", Top Error",
			Baggage:   referenceBaggage,
			DataLayer: referenceErrLayer,
			Code:      referenceErrCode,
		},
		{
			Message:   referenceErrPath + ", Middle Error",
			Baggage:   referenceBaggage,
			DataLayer: referenceErrLayer,
			Code:      referenceErrCode,
		},
		{
			Message:   referenceErrPath + ", Last Error",
			Baggage:   referenceBaggage,
			DataLayer: referenceErrLayer,
			Code:      referenceErrCode,
		},
		{
			Message: referenceErrorText,
		},
	}
	requires.Len(topErr.GetFullTraceSlice(), 4, "check error referenceStack len")
	requires.Len(middleErr.GetFullTraceSlice(), 3, "check error referenceStack len")
	requires.Len(lastErr.GetFullTraceSlice(), 2, "check error referenceStack len")
	assertions.Equal(referenceStack, topErr.GetFullTraceSlice())
}

func Test_customErr_GetMessage(t *testing.T) {
	assertions := assert.New(t)
	checkedErr := referenceError()
	checkedErr.message = "CRITICAL Error"
	assertions.Equal(ErrorMessage(referenceErrPath+", CRITICAL Error"), checkedErr.GetMessage())
}

func Test_customErr_IsErrorExistInStack(t *testing.T) {
	assertions := assert.New(t)

	topErr := referenceError()
	topErr.message = "Top Error"
	topErr.dataLayer = DataService
	topErr.code = AccessDenied
	middleErr := referenceError()
	middleErr.message = "Middle Error"
	middleErr.dataLayer = UseCase
	middleErr.code = Unauthorized
	lastErr := referenceError()
	lastErr.message = "Last Error"
	lastErr.dataLayer = Container
	lastErr.code = InvalidArguments

	middleErr.nativeErr = &lastErr
	topErr.nativeErr = &middleErr

	assertions.Equal(true, topErr.IsErrorExistInStack(AccessDenied, DataService))
	assertions.Equal(false, middleErr.IsErrorExistInStack(AccessDenied, DataService))
	assertions.Equal(true, lastErr.IsErrorExistInStack(InvalidArguments, Container))
}

func Test_customErr_IsErrorWithCodeExistInStack(t *testing.T) {
	assertions := assert.New(t)

	topErr := referenceError()
	topErr.message = "Top Error"
	topErr.code = AccessDenied
	middleErr := referenceError()
	middleErr.message = "Middle Error"
	middleErr.code = Unauthorized
	lastErr := referenceError()
	lastErr.message = "Last Error"
	lastErr.code = Unauthorized

	middleErr.nativeErr = &lastErr
	topErr.nativeErr = &middleErr

	assertions.Equal(true, topErr.IsErrorWithCodeExistInStack(AccessDenied))
	assertions.Equal(true, middleErr.IsErrorWithCodeExistInStack(Unauthorized))
	assertions.Equal(false, lastErr.IsErrorWithCodeExistInStack(InvalidArguments))
}

func Test_customErr_AddOperation(t *testing.T) {
	assertions := assert.New(t)

	refErr := referenceError()
	refErr.message = "Top Error"
	refErr.dataLayer = DataService
	refErr.code = AccessDenied

	referenceBaggageClone := referenceBaggage
	referenceBaggageClone["test"] = "test"

	wrappedErr := refErr.AddOperation("Wrapped Error", referenceBaggageClone, referenceSeverity)
	assertions.Len(refErr.GetFullTraceSlice(), 2, "Check original Err stack length")
	assertions.Len(wrappedErr.GetFullTraceSlice(), 3, "Check wrapped Err stack length")
	assertions.Equal(refErr.GetCode(), wrappedErr.GetCode(), "Check that ref code was cloned to wrapped Err code")
	assertions.Equal(refErr.GetDataLevel(), wrappedErr.GetDataLevel(),
		"Check that ref dataLayer was cloned to wrapped Err dataLayer")
	assertions.Equal(&refErr, wrappedErr.Unwrap(), "Check Relation between errors")
	assertions.Equal(referenceBaggageClone, wrappedErr.GetBaggage(),
		"Check that wrapped error baggage contains all changes")
	assertions.Equal(referenceBaggage, refErr.GetBaggage(), "Check that refErr baggage doesnt change")
}

func Test_customErrs_GetErr(t *testing.T) {
	assertions := assert.New(t)
	errs := NewMultiply()
	assertions.Len(errs.GetErrs(), 0)
	errs.AddErr(New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	errs.AddErr(New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	errs.AddErr(New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	assertions.Len(errs.GetErrs(), 3)
}

func Test_customErrs_AddErr(t *testing.T) {
	assertions := assert.New(t)
	requires := require.New(t)
	errs := NewMultiply()
	errs.AddErr(New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	requires.Len(errs.GetErrs(), 1)
	assertions.Equal(errs.GetErrs()[0].GetCode(), referenceErrCode)
}

func Test_customErrs_IsEmpty(t *testing.T) {
	assertions := assert.New(t)
	errs := NewMultiply()
	assertions.Len(errs.GetErrs(), 0)
	assertions.Equal(true, errs.IsEmpty())
	errs.AddErr(New(referenceErrCode, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	assertions.Equal(false, errs.IsEmpty())
}

func Test_customErrs_isErrorExist(t *testing.T) {
	assertions := assert.New(t)
	errs := NewMultiply()
	errs.AddErr(New(InternalError, referenceMessage, UseCase, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	errs.AddErr(New(InvalidArguments, referenceMessage, Controller, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))
	errs.AddErr(New(Unauthorized, referenceMessage, Transport, referenceErrPath, referenceBaggage,
		referenceSeverity, errNativeReference))

	assertions.Equal(false, errs.isErrorExist(InternalError, Controller))
	assertions.Equal(true, errs.isErrorExist(InternalError, UseCase))
	assertions.Equal(true, errs.isErrorExist(Unauthorized, Transport))
	assertions.Equal(true, errs.isErrorExist(InvalidArguments, Controller))
	assertions.Equal(false, errs.isErrorExist(Unauthorized, Controller))
}
