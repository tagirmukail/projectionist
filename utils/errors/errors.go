package errors

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	httpErrMsgPtrn = "[code:%d][time:%s] - %s"
	errMsgPtrn     = "[httpCode:%d][code:%d][time:%s] - error: %s"
)

var (
	timeZone   *time.Location
	mx         = sync.RWMutex{}
	knownCodes = knownCodesWithErrMsg{}
)

type knownCodesWithErrMsg map[Code]IError

type IError interface {
	Code() Code
	Error() string
	HTTPCode() int
	HTTPError() string
	SetArgs(args ...string)
}

type Error struct {
	t             time.Time
	code          Code
	err           string
	httpCode      int
	httpError     string
	keyValuesArgs []string
}

type Code int

// New - new error
func New(code Code, httpCode int, httpError string, err string) IError {
	checkTimeZone()

	e := &Error{
		t:             time.Now().In(timeZone),
		code:          code,
		err:           err,
		httpCode:      httpCode,
		httpError:     httpError,
		keyValuesArgs: []string{},
	}

	mx.Lock()
	knownCodes[code] = e
	mx.Unlock()

	return e
}

func Newf(code Code, httpCode int, httpError string, errPtrn string, keyValuesArgs ...string) IError {
	checkTimeZone()

	e := &Error{
		t:             time.Now().In(timeZone),
		code:          code,
		err:           errPtrn,
		httpCode:      httpCode,
		httpError:     httpError,
		keyValuesArgs: keyValuesArgs,
	}

	mx.Lock()
	knownCodes[code] = e
	mx.Unlock()

	return e
}

func Get(code Code) IError {
	mx.RLock()
	defer mx.RUnlock()
	return knownCodes[code]
}

// Code - returned application code
func (e *Error) Code() Code {
	return e.code
}

// Error - returned error for log
func (e *Error) Error() string {
	b := strings.Builder{}
	b.WriteString(e.err)
	if len(e.keyValuesArgs) > 0 {
		b.WriteString(", ")
	}

	var kvCount = 0
	for _, item := range e.keyValuesArgs {
		switch kvCount {
		case 0:
			b.WriteString(item)
		case 1:
			b.WriteString(": ")
			b.WriteString(item)
		case 2:
			kvCount = 0
			b.WriteString(", ")
		default:
		}

		kvCount++
	}

	return fmt.Sprintf(errMsgPtrn, e.httpCode, e.code, e.t.Format(time.RFC3339Nano), b.String())
}

// HTTPCode - returned http error code
func (e *Error) HTTPCode() int {
	return e.httpCode
}

// HTTPError - error message for http client with code
func (e *Error) HTTPError() string {
	return fmt.Sprintf(httpErrMsgPtrn, e.code, e.t.Format(time.RFC1123), e.httpError)
}

func (e *Error) SetArgs(args ...string) {
	e.keyValuesArgs = append(e.keyValuesArgs, args...)
}

// SetTimeZone - set time zone for errors
// location - for example "America/New_York"
func SetTimeZone(name string) error {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return err
	}

	timeZone = loc

	return nil
}

// checkTimeZone - check installed time location
// default: UTC
func checkTimeZone() {
	if timeZone == nil {
		timeZone = time.UTC
	}
}
