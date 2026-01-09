package testutil

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// AssertNoError fails the test if err is not nil with a clear message.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// AssertError fails the test if err is nil, optionally checks if error message contains substring.
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if msg != "" && !strings.Contains(err.Error(), msg) {
		t.Errorf("expected error to contain %q, got: %v", msg, err)
	}
}

// AssertEqual compares values using deep equality, fails with diff on mismatch.
func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %+v\nwant: %+v", got, want)
	}
}

// AssertJSONEqual parses and compares JSON strings, fails with structural diff.
func AssertJSONEqual(t *testing.T, got, want string) {
	t.Helper()

	var gotData, wantData interface{}

	if err := json.Unmarshal([]byte(got), &gotData); err != nil {
		t.Fatalf("failed to unmarshal got JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(want), &wantData); err != nil {
		t.Fatalf("failed to unmarshal want JSON: %v", err)
	}

	if !reflect.DeepEqual(gotData, wantData) {
		t.Errorf("\ngot JSON:  %s\nwant JSON: %s", got, want)
	}
}

// NewTestContext creates a context with timeout (default 5s if zero).
func NewTestContext(timeout time.Duration) context.Context {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

// NewTestHTTPRequest creates an httptest request with JSON body marshaling.
func NewTestHTTPRequest(method, path string, body interface{}) *http.Request {
	var req *http.Request

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			panic("failed to marshal request body: " + err.Error())
		}
		req = httptest.NewRequest(method, path, strings.NewReader(string(jsonBytes)))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	return req
}
