package util

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

//--------------------
// HELPER
//--------------------

// github.com/tideland's code

// callInfo bundles the info about the call environment
// when a logging statement occured.
type CallInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

// shortFormat returns a string representation in a short variant.
func (ci *CallInfo) ShortFormat() string {
	return fmt.Sprintf("[%s]", ci.packageName)
}

// verboseFormat returns a string representation in a more verbose variant.
func (ci *CallInfo) VerboseFormat() string {
	return fmt.Sprintf("[%s] (%s:%s:%d)", ci.packageName, ci.fileName, ci.funcName, ci.line)
}

// retrieveCallInfo
func RetrieveCallInfos(ns int) (calls []*CallInfo) {
	calls = make([]*CallInfo, ns+3)
	for i := 3; i < len(calls); i++ {
		c := RetrieveCallInfo(i)
		if nil == c {
			break
		}
		calls[i] = c
	}
	return
}

func Lfmt(n int) string {
	ci := RetrieveCallInfo(n)
	return ci.VerboseFormat()
}

func RetrieveCallInfo(n int) *CallInfo {
	pc, file, line, _ := runtime.Caller(n)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	if len(parts) <= (pl-1) || 0 > (pl-1) {
		return nil
	}
	funcName := parts[pl-1]

	if len(parts) <= (pl-2) || 0 > (pl-2) {
		return nil
	}
	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &CallInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}
