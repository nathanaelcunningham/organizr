package testutil

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestAssertNoError(t *testing.T) {
	t.Run("passes when error is nil", func(t *testing.T) {
		AssertNoError(t, nil)
	})
}

func TestAssertError(t *testing.T) {
	t.Run("passes when error is not nil", func(t *testing.T) {
		AssertError(t, errors.New("test error"), "")
	})

	t.Run("checks error message contains substring", func(t *testing.T) {
		err := errors.New("test error message")
		AssertError(t, err, "error message")
	})
}

func TestAssertEqual(t *testing.T) {
	t.Run("passes when values are equal", func(t *testing.T) {
		AssertEqual(t, 42, 42)
		AssertEqual(t, "test", "test")
		AssertEqual(t, []int{1, 2, 3}, []int{1, 2, 3})
	})
}

func TestAssertJSONEqual(t *testing.T) {
	t.Run("passes when JSON is structurally equal", func(t *testing.T) {
		got := `{"name":"test","value":42}`
		want := `{"value":42,"name":"test"}`
		AssertJSONEqual(t, got, want)
	})
}

func TestNewTestContext(t *testing.T) {
	t.Run("creates context with default timeout", func(t *testing.T) {
		ctx := NewTestContext(0)
		if ctx == nil {
			t.Fatal("expected context, got nil")
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("expected context to have deadline")
		}

		// Check that deadline is approximately 5 seconds from now
		expectedDeadline := time.Now().Add(5 * time.Second)
		if deadline.Before(expectedDeadline.Add(-100*time.Millisecond)) ||
			deadline.After(expectedDeadline.Add(100*time.Millisecond)) {
			t.Errorf("expected deadline around %v, got %v", expectedDeadline, deadline)
		}
	})

	t.Run("creates context with custom timeout", func(t *testing.T) {
		ctx := NewTestContext(10 * time.Second)
		if ctx == nil {
			t.Fatal("expected context, got nil")
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("expected context to have deadline")
		}

		// Check that deadline is approximately 10 seconds from now
		expectedDeadline := time.Now().Add(10 * time.Second)
		if deadline.Before(expectedDeadline.Add(-100*time.Millisecond)) ||
			deadline.After(expectedDeadline.Add(100*time.Millisecond)) {
			t.Errorf("expected deadline around %v, got %v", expectedDeadline, deadline)
		}
	})

	t.Run("context has no parent cancel", func(t *testing.T) {
		ctx := NewTestContext(1 * time.Second)
		select {
		case <-ctx.Done():
			t.Error("context should not be done immediately")
		default:
			// Expected
		}
	})
}

func TestNewTestHTTPRequest(t *testing.T) {
	t.Run("creates request without body", func(t *testing.T) {
		req := NewTestHTTPRequest(http.MethodGet, "/test", nil)
		if req == nil {
			t.Fatal("expected request, got nil")
		}

		if req.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", req.Method)
		}

		if req.URL.Path != "/test" {
			t.Errorf("expected path /test, got %s", req.URL.Path)
		}
	})

	t.Run("creates request with JSON body", func(t *testing.T) {
		body := map[string]interface{}{
			"name":  "test",
			"value": 42,
		}

		req := NewTestHTTPRequest(http.MethodPost, "/api/test", body)
		if req == nil {
			t.Fatal("expected request, got nil")
		}

		if req.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", req.Method)
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
		}

		if req.Body == nil {
			t.Error("expected request body to be set")
		}
	})

	t.Run("panics with invalid body", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()

		// channels cannot be marshaled to JSON
		invalidBody := make(chan int)
		NewTestHTTPRequest(http.MethodPost, "/test", invalidBody)
	})
}
