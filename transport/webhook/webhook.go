// Package webhook contains an HTTP handler for incoming Telegram webhook updates.
package webhook

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"mime"
	"net/http"

	"ai-gram/internal/telegramsecret"
	"ai-gram/telegram"
)

const (
	secretTokenHeader   = "X-Telegram-Bot-Api-Secret-Token"
	defaultMaxBodyBytes = 1 << 20
)

// Handler handles one Telegram update.
type Handler interface {
	HandleUpdate(context.Context, telegram.Update) error
}

// HandlerFunc adapts a function to Handler.
type HandlerFunc func(context.Context, telegram.Update) error

// HandleUpdate calls f(ctx, update).
func (f HandlerFunc) HandleUpdate(ctx context.Context, update telegram.Update) error {
	return f(ctx, update)
}

// Config configures a webhook HTTP handler.
type Config struct {
	// SecretToken is the optional Telegram webhook secret token.
	SecretToken string
	// MaxBodyBytes limits the accepted request body size.
	MaxBodyBytes int64
	// OnError handles errors returned by the update handler.
	OnError func(context.Context, *telegram.Update, error)
}

type receiver struct {
	handler      Handler
	secretToken  string
	maxBodyBytes int64
	onError      func(context.Context, *telegram.Update, error)
}

// New creates an HTTP handler for Telegram webhook requests.
func New(handler Handler, config Config) (http.Handler, error) {
	if handler == nil {
		return nil, stderrors.New("handler is required")
	}
	if err := telegramsecret.Validate(config.SecretToken); err != nil {
		return nil, err
	}
	maxBodyBytes := config.MaxBodyBytes
	if maxBodyBytes <= 0 {
		maxBodyBytes = defaultMaxBodyBytes
	}

	return &receiver{
		handler:      handler,
		secretToken:  config.SecretToken,
		maxBodyBytes: maxBodyBytes,
		onError:      config.OnError,
	}, nil
}

func (r *receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeSafeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !isJSONContentType(req.Header.Get("Content-Type")) {
		writeSafeError(w, http.StatusUnsupportedMediaType, "unsupported media type")
		return
	}
	if r.secretToken != "" && !secretTokenMatches(r.secretToken, req.Header.Get(secretTokenHeader)) {
		writeSafeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	body := http.MaxBytesReader(w, req.Body, r.maxBodyBytes)
	defer body.Close()

	var update telegram.Update
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&update); err != nil {
		writeSafeError(w, statusForDecodeError(err), "bad request")
		return
	}
	if decoder.More() {
		writeSafeError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := r.handler.HandleUpdate(req.Context(), update); err != nil {
		if r.onError != nil {
			r.onError(req.Context(), &update, err)
		}
		writeSafeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func isJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	return err == nil && mediaType == "application/json"
}

func secretTokenMatches(expected string, actual string) bool {
	if expected == "" {
		return true
	}
	if len(expected) != len(actual) {
		return false
	}

	return telegramsecret.ConstantTimeEqual(actual, expected)
}

func statusForDecodeError(err error) int {
	var maxBytesErr *http.MaxBytesError
	if stderrors.As(err, &maxBytesErr) {
		return http.StatusRequestEntityTooLarge
	}

	return http.StatusBadRequest
}

func writeSafeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message + "\n"))
}
