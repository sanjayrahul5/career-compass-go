package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFrame returns a formatted string representing the frame of the call
func GetFrame(function uintptr, file string, line int, _ bool) string {
	absPath, _ := filepath.Rel(strings.Split(file, "career-compass-go")[0]+"career-compass-go", file)

	arr := strings.Split(runtime.FuncForPC(function).Name(), ".")
	funcName := arr[len(arr)-1]
	if funcName == "0" {
		funcName = arr[len(arr)-1]
	}

	return fmt.Sprintf("[%s][%s][%d] ", absPath, funcName, line)
}
