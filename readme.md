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
import cErrors "github.com/Darevski/go-custom-errors"

//....

func MyFunc() error{
    // Create new Custom Error that implements standard golang Error interface
	err := cErrors.New(ErrorType, ErrorLevel, ErrorBaggage, ErrorSeverity, ErrorMessage)
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
        err = cErrors.Wrap(err, "Error on JSON read")
        
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
        err := cErrors.WrapF(err, "Read error for %s", r)
        err.AddBaggage(cErrors.ErrorBaggage{
            "JsonReaderID" : someModel.ID
        })
    }
    
    //......
}
```

### Available Interfaces 

All methods of **CustomError** and **MultipleCustomErrs** can be viewed at [this file](interfaces.go)


### Error Codes && DataLayers|Severity Levels

The following error codes, dataLayer and severity levels are currently available:

``` go
const (
	DefaultLevel = ErrorLevel(iota)
	DataLevel
	UseCaseLevel
	ContainerLevel
	ControllerLevel
	TransportLevel
)

const (
	DefaultSeverity = ErrorSeverity(iota)
	Debug
	Info
	Warning
	Critical
	Fatal
	Panic
)

const (
	DefaultType = ErrorType(iota)
	NotFound
	InvalidArguments

	InternalError
	BadRequest

	AccessDenied
	Unauthorized
)

```

Fell free to use any INT code for own codes/levels