package customErrors

import (
	"fmt"
	"runtime"
)

// DetectPath return function call line
func DetectPath(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		return fmt.Sprintf("Called from %s, line #%d", file, line)
	}
	return ""
}
