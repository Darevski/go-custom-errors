package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_customErr_AddErr(t *testing.T) {
	assertions := assert.New(t)
	var cErrs []CustomError

	errs := NewMultiply()

	for k := 0; k < 100; k++ {
		err := NewBase(ErrorMessage(fmt.Sprintf("%d error", k)))
		errs.AddErr(err)
		cErrs = append(cErrs, err)
	}

	for k, v := range errs.GetErrs() {
		assertions.Equal(cErrs[k], v, "Checking err slice")
	}
}


func Test_customErrs_IsEmpty(t *testing.T) {
	assertions := assert.New(t)
	errs := NewMultiply()
	assertions.True(errs.IsEmpty())
	for k := 0; k < 100; k++ {
		errs.AddErr(NewBase(ErrorMessage(fmt.Sprintf("%d error", k))))
	}
	assertions.False(errs.IsEmpty())
}

func Test_customErrs_isErrorExist(t *testing.T) {
	assertions := assert.New(t)
	errs := NewMultiply()
	errs.AddErr(NotFound.NewBase("Not Found Test"))
	errs.AddErr(InvalidArguments.NewBase("Not Found Test"))
	errs.AddErr(InternalError.NewBase("Not Found Test"))


	assertions.True(errs.IsErrorExist(NotFound.NewBase("Test Not Found")))
	assertions.False(errs.IsErrorExist(BadRequest.NewBase("Test Not Found")))
}
