package rui

import (
	"fmt"
	"runtime"
)

// ProtocolInDebugLog If it is set to true, then the protocol of the exchange between
// clients and the server is displayed in the debug log
var ProtocolInDebugLog = false

var debugLogFunc func(string) = debugLog
var errorLogFunc func(string) = errorLog

// SetDebugLog sets a function for outputting debug info.
// The default value is nil (debug info is ignored)
func SetDebugLog(f func(string)) {
	debugLogFunc = f
}

// SetErrorLog sets a function for outputting error messages.
// The default value is log.Println(text)
func SetErrorLog(f func(string)) {
	errorLogFunc = f
}

// DebugLog print the text to the debug log
func DebugLog(text string) {
	if debugLogFunc != nil {
		debugLogFunc(text)
	}
}

// DebugLogF print the text to the debug log
func DebugLogF(format string, a ...any) {
	if debugLogFunc != nil {
		debugLogFunc(fmt.Sprintf(format, a...))
	}
}

var lastError = ""

// ErrorLog print the text to the error log
func ErrorLog(text string) {
	lastError = text
	if errorLogFunc != nil {
		errorLogFunc(text)
		errorStack()
	}
}

// ErrorLogF print the text to the error log
func ErrorLogF(format string, a ...any) {
	lastError = fmt.Sprintf(format, a...)
	if errorLogFunc != nil {
		errorLogFunc(lastError)
		errorStack()
	}
}

// LastError returns the last error text
func LastError() string {
	return lastError
}

func errorStack() {
	if errorLogFunc != nil {
		skip := 2
		_, file, line, ok := runtime.Caller(skip)
		for ok {
			errorLogFunc(fmt.Sprintf("\t%s: line %d", file, line))
			skip++
			_, file, line, ok = runtime.Caller(skip)
		}
	}
}
