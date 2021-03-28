# Go Custom Errors

Extended Errors for Go

## What? Why?

It is extremely important to be able to effectively debug problems in production. However, depending
on the simplicty of your language's error handling, managing errors and debugging can become tedious
and difficult. For instance, GoLang provides a basic error which contains a simple value - but it
lacks verbose information about the cause of the error or the state of the stack at error-time.

Custom Errors uses original GoLang built-in error types but stores additional extra data.

These are the next extra extended features that present in package:

1. **Wrapping Error** - the original error which happened in the system or error that has been created by this package
   could be wrapped using this package.
2. **Wrapping Stack Generation** - errors wrapping allow package to generate stack of errors
3. **Additional Information** - every custom error has information about place of error appearing, error code, 
   information of error data layer level, error severity
4. **Baggage Data** - a map which could be used to store custom data associated with an error.
5. **Multiple Error Storage** - package provides the ability to store multiple errors, with further operations with them,
like search etc.

## Quick Usage

### Basics

```go
import cErrors "gitlab.com/73five.com/go-custom-errors"

//....

func MyFunc() error{
    // Create new Custom Error that implements standard golang Error interface
	err := cErrors.NewErr(ErrorCode, ErrorMessage, DataLayerLevel, cErrors.DetectPath(cErrors.SkipFunctionHelper), Baggage, err)
	return err
    // .....
}

//....

if data, err := MyFunc(); err != nil {

    // Check Error of MyFunc function && if it`s Custom Error then check is it NotFound error 
	// If it is, then log that with custom message
    if custom, ok := err.(*cErrors.CustomError); ok && custom.IsErrorWithCodeExistInStack(cErrors.NotFound) {
    	logger.Log(custom.GetMessage())
    }
    //....
    
}
```

### Custom data in error

```go
func DoSomething(w http.ResponseWriter, r *http.Request) {
	
    //......
	
    if err = ReadJSONFromReader(r, someModel); err != nil {
        err = cErrors.NewErr(ErrorCode, ErrorMessage, DataLayerLevel, cErrors.DetectPath(cErrors.SkipFunctionHelper), nil, err)
        
        err.AddBaggage(cErrors.ErrorBaggage{
            "JsonReaderID" : someModel.ID
        })
    }
    
    //......
}
```

### Wrap error


```go
func DoSomething(w http.ResponseWriter, r *http.Request) {
	
    //......
	
    if err = ReadJSONFromReader(r, someModel); err != nil {
        err := cErrors.Wrap(err, "ReadJSONError")
        err.AddBaggage(cErrors.ErrorBaggage{
            "JsonReaderID" : someModel.ID
        })
    }
    
    //......
}
```

## Docs

#### func New

```go
func New(ErrorCode,ErrorMessage, ErrorDataLevel, ErrorPath, ErrorBaggage, error) CustomError
```

New create new custom error with provided params
Note it will also used for error wrapping with the possibility of creating and later getting an error stack


#### func Wrap

```go
func Wrap(err error, message ErrorMessage) CustomError
```

Wrap is a simplified version of the New function that will create custom error with empty error additional data

#### func NewMultiply

```go
func NewMultiply() MultipleCustomErrs 
```

NewMultiply create Multiple Errors representation struct that allowed to use MultipleCustomErrs interface


#### func AddOperation

```go
func AddOperation(message ErrorMessage, baggage ErrorBaggage) CustomError
```

AddOperation is a simplified version of the New function that allows to override an error by adding only a message
and baggage the error code and level are taken from the original error, ErrPath will be automatically written with
the place where the AddOperation function is called


### Available Interfaces 

All methods of **CustomError** and **MultipleCustomErrs** can be viewed at [this file](interfaces.go)


### Error Codes && DataLayers|Severity Levels

The following error codes, dataLayer and severity levels are currently available:

``` go
const (
	InternalError    ErrorCode = 1500
	InvalidArguments ErrorCode = 1412
	NotFound         ErrorCode = 1404
	AccessDenied     ErrorCode = 1403
	Unauthorized     ErrorCode = 1401
)

const (
	DataService ErrorDataLevel = 101
	UseCase     ErrorDataLevel = 102
	Container   ErrorDataLevel = 103
	Controller  ErrorDataLevel = 104
	Transport   ErrorDataLevel = 105
)

const (
	Debug    ErrorSeverity = 0
	Info     ErrorSeverity = 1
	Warning  ErrorSeverity = 2
	Critical ErrorSeverity = 3
	Fatal    ErrorSeverity = 4
	Panic    ErrorSeverity = 5
)

```

Fell free to use any INT code for own codes/levels