package util

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
)

type LogFunc func(fmt string, params ...interface{})

type LogFuncType int

const (
	LogErrName LogFuncType = iota
	LogInfoName
)

func WithErrLogger(ctx context.Context, fn LogFunc) context.Context {
	return context.WithValue(ctx, LogErrName, fn)
}

func WithInfoLogger(ctx context.Context, fn LogFunc) context.Context {
	return context.WithValue(ctx, LogInfoName, fn)
}

func GetErrLog(ctx context.Context) LogFunc {
	return GetLogFunc(ctx, LogErrName)
}

func GetInfoLog(ctx context.Context) LogFunc {
	return GetLogFunc(ctx, LogInfoName)
}

func GetLogFunc(ctx context.Context, name LogFuncType) (fn LogFunc) {
	if nil != ctx.Value(name) {
		fn = ctx.Value(name).(LogFunc)
	} else {
		fn = log.Printf
	}
	return
}

func LogErr(ctx context.Context, fmt string, params ...interface{}) {
	GetErrLog(ctx)(fmt, params...)
}

func LogInfo(ctx context.Context, fmt string, params ...interface{}) {
	GetInfoLog(ctx)(fmt, params...)
}

func StackTrace(n int) string {
	// purposely create _stackTrace, we wanted to pass 'j'
	return _stackTrace(n, 4)
}

func (this LogFunc) PrintStackTrace(n int) {
	this("%s", _stackTrace(n, 4))
}

func (this LogFunc) Ln(as ...interface{}) {
	this("%s", fmt.Sprintln(as...))
}

func _stackTrace(n, j int) string {
	calls := retrieveCallInfos(n, j)
	sb := bytes.NewBufferString("")
	for _, v := range calls {
		if nil != v {
			sb.WriteString(v.verboseFormat())
			sb.WriteString("\n\t")
		}
	}
	return sb.String()
}

//--------------------
// HELPER
//--------------------

// github.com/tideland's code

// callInfo bundles the info about the call environment
// when a logging statement occured.
type _callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

// shortFormat returns a string representation in a short variant.
func (ci *_callInfo) shortFormat() string {
	return fmt.Sprintf("[%s]", ci.packageName)
}

// verboseFormat returns a string representation in a more verbose variant.
func (ci *_callInfo) verboseFormat() string {
	return fmt.Sprintf("[%s] (%s:%s:%d)", ci.packageName, ci.fileName, ci.funcName, ci.line)
}

// retrieveCallInfo
func retrieveCallInfos(ns, j int) (calls []*_callInfo) {
	calls = make([]*_callInfo, ns)
	for i := 0; i < len(calls); i++ {
		c := retrieveCallInfo(j + i)
		if nil == c {
			break
		}
		calls[i] = c
	}
	return
}

func lfmt(n int) string {
	ci := retrieveCallInfo(n)
	return ci.verboseFormat()
}

func retrieveCallInfo(n int) *_callInfo {
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

	return &_callInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}
