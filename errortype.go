package errors

import (
	errs "github.com/pkg/errors"
)

type ErrorType uint

// New create new custom error with provided params and type based on ErrorType
func (i ErrorType) New(
	errDataLevel ErrorLevel, baggage ErrorBaggage,
	severity ErrorSeverity, message ErrorMessage,
) CustomError {
	// Initialize empty ErrorBaggage map to prevent panics
	if baggage == nil {
		baggage = make(ErrorBaggage)
	}
	return newCustomErr(i, baggage, errDataLevel, severity, errs.New(message.String()))
}

// NewF create custom error with params && error message that formats according to a format specifier and type based on ErrorType
func (i ErrorType) NewF(
	errDataLevel ErrorLevel, baggage ErrorBaggage,
	severity ErrorSeverity, format string, args ...interface{},
) CustomError {
	// Initialize empty ErrorBaggage map to prevent panics
	if baggage == nil {
		baggage = make(ErrorBaggage)
	}

	return newCustomErr(i, baggage, errDataLevel, severity, errs.Errorf(format, args...))
}

// NewBase create custom error with specified message
// also all error attributes are set to default values, but Error Type set`s up based on ErrorType
func (i ErrorType) NewBase(message ErrorMessage) CustomError {
	return newCustomErr(i, make(ErrorBaggage), DefaultLevel, DefaultSeverity, errs.WithStack(errs.New(message.String())))
}

// NewBaseF create custom error with specified message
// also all error attributes are set to default values, but Error Type set`s up based on ErrorType
func (i ErrorType) NewBaseF(format string, args ...interface{}) CustomError {
	return newCustomErr(i, make(ErrorBaggage), DefaultLevel, DefaultSeverity, errs.Errorf(format, args...))
}

// Wrap is a simplified version of the NewBase function that will create custom error with empty error additional data
func (i ErrorType) Wrap(err error, message ErrorMessage) CustomError {
	wrappedErr := errs.Wrap(err, message.String())
	if customErr, ok := err.(CustomError); ok {
		return newCustomErr(i, make(ErrorBaggage), customErr.GetLevel(), customErr.GetSeverity(), wrappedErr)
	}
	return newCustomErr(i, make(ErrorBaggage), DefaultLevel, DefaultSeverity, wrappedErr)
}

// WrapF returns an error annotating err with a stack trace
// at the point WrapF is called, and the format specifier.
func (i ErrorType) WrapF(err error, format string, args ...interface{}) CustomError {
	wrappedErr := errs.Wrapf(err, format, args...)
	if customErr, ok := err.(CustomError); ok {
		return newCustomErr(i, make(ErrorBaggage), customErr.GetLevel(), customErr.GetSeverity(), wrappedErr)
	}
	return newCustomErr(i, make(ErrorBaggage), DefaultLevel, DefaultSeverity, wrappedErr)
}
