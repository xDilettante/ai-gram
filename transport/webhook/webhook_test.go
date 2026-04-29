package webhook

import (
	"context"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ai-gram/dispatch"
	"ai-gram/telegram"
)

func TestNewValidation(t *testing.T) {
	_, err := New(nil, Config{})
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
	assertNoToken(t, err, "123:secret")

	for _, token := range []string{"bad token", "bad.token", strings.Repeat("a", 257)} {
		_, err := New(HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{SecretToken: token})
		if err == nil {
			t.Fatalf("expected error for token %q", token)
		}
		assertNoToken(t, err, token)
	}

	handler, err := New(HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{SecretToken: "valid_TOKEN-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if handler == nil {
		t.Fatal("expected handler")
	}
}

func TestMethodsAndContentType(t *testing.T) {
	handler := newTestHandler(t, nil, Config{})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/webhook", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET status = %d", recorder.Code)
	}

	tests := []struct {
		name        string
		contentType string
		wantStatus  int
	}{
		{name: "missing", wantStatus: http.StatusUnsupportedMediaType},
		{name: "text", contentType: "text/plain", wantStatus: http.StatusUnsupportedMediaType},
		{name: "malformed", contentType: "application/json; charset", wantStatus: http.StatusUnsupportedMediaType},
		{name: "json", contentType: "application/json", wantStatus: http.StatusOK},
		{name: "json charset", contentType: "application/json; charset=utf-8", wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(validUpdateJSON("hello")))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			handler.ServeHTTP(recorder, req)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d, body=%q", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
		})
	}
}

func TestSecretTokenChecks(t *testing.T) {
	const secret = "secret_TOKEN-123"
	var calls int
	handler := newTestHandler(t, HandlerFunc(func(context.Context, telegram.Update) error {
		calls++
		return nil
	}), Config{SecretToken: secret})

	for _, tt := range []struct {
		name   string
		header string
		want   int
	}{
		{name: "missing", want: http.StatusUnauthorized},
		{name: "wrong", header: "wrong", want: http.StatusUnauthorized},
		{name: "correct", header: secret, want: http.StatusOK},
	} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := newJSONRequest(validUpdateJSON("hello"))
			if tt.header != "" {
				req.Header.Set(secretTokenHeader, tt.header)
			}
			handler.ServeHTTP(recorder, req)
			if recorder.Code != tt.want {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.want)
			}
		})
	}
	if calls != 1 {
		t.Fatalf("expected exactly one handled request, got %d", calls)
	}
}

func TestSecretTokenMatches(t *testing.T) {
	if !secretTokenMatches("secret", "secret") {
		t.Fatal("expected matching secrets")
	}
	if secretTokenMatches("secret", "wrong") {
		t.Fatal("expected mismatch")
	}
	if secretTokenMatches("secret", "secret-long") {
		t.Fatal("expected length mismatch")
	}
	if !secretTokenMatches("", "anything") {
		t.Fatal("empty expected secret should disable matching")
	}
}

func TestBodyAndJSONHandling(t *testing.T) {
	handler := newTestHandler(t, nil, Config{MaxBodyBytes: 8})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(`{bad`))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("invalid JSON status = %d", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(validUpdateJSON("too large")))
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("large body status = %d", recorder.Code)
	}
}

func TestValidUpdateIsDecodedAndContextPassed(t *testing.T) {
	type contextKey string
	const key contextKey = "key"
	var gotUpdate telegram.Update
	var gotValue any
	var calls int
	handler := newTestHandler(t, HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		calls++
		gotUpdate = update
		gotValue = ctx.Value(key)
		return nil
	}), Config{})

	recorder := httptest.NewRecorder()
	req := newJSONRequest(validUpdateJSON("hello"))
	req = req.WithContext(context.WithValue(req.Context(), key, "value"))
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%q", recorder.Code, recorder.Body.String())
	}
	if calls != 1 {
		t.Fatalf("handler calls = %d", calls)
	}
	if gotValue != "value" {
		t.Fatalf("context value = %#v", gotValue)
	}
	if gotUpdate.UpdateID != 10 || gotUpdate.Message == nil || gotUpdate.Message.Text != "hello" {
		t.Fatalf("unexpected update: %+v", gotUpdate)
	}
}

func TestHandlerErrorReturns500AndCallsOnError(t *testing.T) {
	const tokenLike = "123:secret"
	want := stderrors.New("internal failure " + tokenLike)
	var gotUpdate *telegram.Update
	var gotErr error
	handler := newTestHandler(t, HandlerFunc(func(context.Context, telegram.Update) error {
		return want
	}), Config{OnError: func(ctx context.Context, update *telegram.Update, err error) {
		gotUpdate = update
		gotErr = err
	}})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(validUpdateJSON(tokenLike)))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", recorder.Code)
	}
	if gotUpdate == nil || gotUpdate.Message == nil || gotUpdate.Message.Text != tokenLike {
		t.Fatalf("unexpected OnError update: %+v", gotUpdate)
	}
	if !stderrors.Is(gotErr, want) {
		t.Fatalf("OnError err = %v, want %v", gotErr, want)
	}
	body := recorder.Body.String()
	if strings.Contains(body, "internal failure") || strings.Contains(body, tokenLike) {
		t.Fatalf("response leaked internals: %q", body)
	}
}

func TestInvalidRequestsDoNotCallOnError(t *testing.T) {
	var onErrorCalls int
	handler := newTestHandler(t, nil, Config{OnError: func(context.Context, *telegram.Update, error) {
		onErrorCalls++
	}})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(`{bad`))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", recorder.Code)
	}
	if onErrorCalls != 0 {
		t.Fatalf("OnError should not be called for decode errors, got %d", onErrorCalls)
	}
}

func TestSuccessReturns200AndCallsHandlerOnce(t *testing.T) {
	var calls int
	handler := newTestHandler(t, HandlerFunc(func(context.Context, telegram.Update) error {
		calls++
		return nil
	}), Config{})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(validUpdateJSON("hello")))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if body := recorder.Body.String(); body != "ok\n" {
		t.Fatalf("unexpected body: %q", body)
	}
	if calls != 1 {
		t.Fatalf("handler calls = %d", calls)
	}
}

func TestCompatibilityWithDispatchDispatcher(t *testing.T) {
	dispatcher := dispatch.New()
	var handledText string
	if err := dispatcher.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		handledText = update.Message.Text
		return nil
	}); err != nil {
		t.Fatalf("unexpected dispatch route error: %v", err)
	}
	var _ Handler = dispatcher

	handler := newTestHandler(t, dispatcher, Config{})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(validUpdateJSON("hello")))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if handledText != "hello" {
		t.Fatalf("unexpected handled text: %q", handledText)
	}
}

func TestHandlerFuncCompatibility(t *testing.T) {
	var _ Handler = HandlerFunc(func(context.Context, telegram.Update) error { return nil })

	var calls int
	handler := newTestHandler(t, HandlerFunc(func(context.Context, telegram.Update) error {
		calls++
		return nil
	}), Config{})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, newJSONRequest(validUpdateJSON("hello")))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if calls != 1 {
		t.Fatalf("handler calls = %d", calls)
	}
}

func newTestHandler(t *testing.T, handler Handler, config Config) http.Handler {
	t.Helper()
	if handler == nil {
		handler = HandlerFunc(func(context.Context, telegram.Update) error { return nil })
	}

	httpHandler, err := New(handler, config)
	if err != nil {
		t.Fatalf("unexpected New error: %v", err)
	}

	return httpHandler
}

func newJSONRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func validUpdateJSON(text string) string {
	return `{"update_id":10,"message":{"message_id":20,"chat":{"id":30,"type":"private"},"date":40,"text":"` + text + `"}}`
}

func assertNoToken(t *testing.T, err error, token string) {
	t.Helper()
	if strings.Contains(err.Error(), token) {
		t.Fatalf("error leaked token: %q", err.Error())
	}
}
