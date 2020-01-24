package errors

import (
	"net/http"
	"strings"
	"testing"
)

// TODO: add tests

func TestLogic(t *testing.T) {
	err := SetTimeZone("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	err = New(
		1,                              // code
		http.StatusInternalServerError, // http code
		"something when wrong",         // http client error message
		"database connection is lost",  // error message
	)

	if !strings.Contains(err.Error(), "[httpCode:500][code:1][time:") {
		t.Errorf("error desc not contain %s error: %v", "[httpCode:500][code:1][time:", err)
	}

	if !strings.Contains(err.Error(), "] - error: database connection is lost") {
		t.Errorf("error desc not contain %s error: %v", "[httpCode:500][code:1][time:", err)
	}

	if Get(1) != err {
		t.Errorf("known codes not contain errPtrn: %v", err)
	}
}

func TestNew(t *testing.T) {
	type args struct {
		code      Code
		httpCode  int
		httpError string
		err       string
	}
	tests := []struct {
		name           string
		args           args
		wantErrMsg     string
		wantHttpErrMsg string
	}{
		{
			name: "ok",
			args: args{
				code:      1,
				httpCode:  http.StatusAlreadyReported,
				httpError: "something when wrong",
				err:       "database connection is lost",
			},
			wantErrMsg:     "database connection is lost",
			wantHttpErrMsg: "something when wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.code, tt.args.httpCode, tt.args.httpError, tt.args.err)
			if !strings.Contains(got.Error(), tt.wantErrMsg) {
				t.Errorf("error: %v - not contain %v", got, tt.wantErrMsg)
			}

			if !strings.Contains(got.HTTPError(), tt.wantHttpErrMsg) {
				t.Errorf("error: %v - not contain %v", got.HTTPError(), tt.wantHttpErrMsg)
			}
		})
	}
}

func TestNewf(t *testing.T) {
	type args struct {
		code      Code
		httpCode  int
		httpError string
		errPtrn   string
		args      []interface{}
	}
	tests := []struct {
		name           string
		args           args
		wantErrMsg     string
		wantHttpErrMsg string
	}{
		{
			name: "ok",
			args: args{
				code:      1,
				httpCode:  404,
				httpError: "not found",
				errPtrn:   "%v %d not exist in db",
				args:      []interface{}{"book", 1},
			},
			wantErrMsg:     "book 1 not exist in db",
			wantHttpErrMsg: "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Newf(tt.args.code, tt.args.httpCode, tt.args.httpError, tt.args.errPtrn, tt.args.args...)
			if !strings.Contains(got.Error(), tt.wantErrMsg) {
				t.Errorf("error: %v - not contain %v", got, tt.wantErrMsg)
			}

			if !strings.Contains(got.HTTPError(), tt.wantHttpErrMsg) {
				t.Errorf("error: %v - not contain %v", got.HTTPError(), tt.wantHttpErrMsg)
			}
		})
	}
}
