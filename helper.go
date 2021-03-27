package errors

import (
	"fmt"
	"runtime"
)

const (
	SkipFunctionHelper = 1
	skipPackage        = 2
)

// DetectPath return function call line
func DetectPath(skip int) ErrorPath {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		return ErrorPath(fmt.Sprintf("Called from %s, line #%d", file, line))
	}
	return ""
}
