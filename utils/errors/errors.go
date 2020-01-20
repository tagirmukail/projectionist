package errors

import (
	"fmt"
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
}

type Error struct {
	t         time.Time
	code      Code
	err       string
	httpCode  int
	httpError string
}

type Code int

// New - new error
func New(code Code, httpCode int, httpError string, err string) IError {
	checkTimeZone()

	e := &Error{
		t:         time.Now().In(timeZone),
		code:      code,
		err:       err,
		httpCode:  httpCode,
		httpError: httpError,
	}

	mx.Lock()
	knownCodes[code] = e
	mx.Unlock()

	return e
}

func Newf(code Code, httpCode int, httpError string, err string, args ...interface{}) IError {
	checkTimeZone()

	e := &Error{
		t:         time.Now().In(timeZone),
		code:      code,
		err:       fmt.Sprintf(err, args...),
		httpCode:  httpCode,
		httpError: httpError,
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
	return fmt.Sprintf(errMsgPtrn, e.httpCode, e.code, e.t.Format(time.RFC3339Nano), e.err)
}

// HTTPCode - returned http error code
func (e *Error) HTTPCode() int {
	return e.httpCode
}

// HTTPError - error message for http client with code
func (e *Error) HTTPError() string {
	return fmt.Sprintf(httpErrMsgPtrn, e.code, e.t.Format(time.RFC1123), e.httpError)
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
