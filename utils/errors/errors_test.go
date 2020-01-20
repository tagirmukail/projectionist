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
		t.Errorf("known codes not contain err: %v", err)
	}

}
