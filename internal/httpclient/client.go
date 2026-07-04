// Package httpclient contains the internal HTTP transport used by higher-level Bot API clients.
package httpclient

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

// Doer sends an HTTP request and returns an HTTP response.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is a small wrapper around an HTTP transport.
//
// Client does not know about dispatching, polling, webhooks, or Telegram method semantics.
type Client struct {
	doer Doer
}

// RequestError wraps a transport error without exposing request URLs that may contain secrets.
type RequestError struct {
	err error
}

// Error returns a redacted transport error message.
func (e *RequestError) Error() string {
	return "telegram request failed"
}

// Unwrap returns the underlying transport error with sensitive URLs redacted when possible.
func (e *RequestError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

// StatusError reports a non-successful HTTP response from the remote server
// whose body did not look like a Telegram Bot API error payload (JSON error
// bodies are passed through to the caller so they can surface the Telegram
// description instead).
type StatusError struct {
	StatusCode int
	// Method is the Bot API method name (last URL path segment); it never
	// contains the bot token.
	Method string
}

// Error returns a redacted HTTP status error message.
func (e *StatusError) Error() string {
	if e == nil {
		return "telegram HTTP request failed"
	}
	if e.Method != "" {
		return fmt.Sprintf("telegram %s HTTP request failed with status %d", e.Method, e.StatusCode)
	}

	return fmt.Sprintf("telegram HTTP request failed with status %d", e.StatusCode)
}

// New creates an internal HTTP client around doer.
//
// If doer is nil, a new http.Client with a bounded timeout is used.
func New(doer Doer) *Client {
	if doer == nil {
		doer = &http.Client{Timeout: defaultTimeout}
	} else if client, ok := doer.(*http.Client); ok && client == nil {
		doer = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{doer: doer}
}

// Do sends req with ctx attached, reads the full response body, and closes it.
//
// Non-2xx responses whose body looks like a Telegram Bot API error payload
// ({"ok":false,...}) are returned to the caller as a body without an error:
// the Bot API reports errors with matching HTTP statuses, and the JSON
// description ("Bad Request: message is not modified", retry_after, ...) is
// far more useful than a bare status code. Non-JSON error bodies (proxies,
// gateways) still yield a StatusError.
func (c *Client) Do(ctx context.Context, req *http.Request) ([]byte, error) {
	resp, err := c.open(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read telegram response body: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if looksLikeTelegramPayload(body) {
			return body, nil
		}
		return nil, &StatusError{StatusCode: resp.StatusCode, Method: methodFromRequest(req)}
	}

	return body, nil
}

// looksLikeTelegramPayload reports whether body is plausibly a Bot API JSON
// response (as opposed to an HTML error page from a proxy or gateway).
func looksLikeTelegramPayload(body []byte) bool {
	trimmed := bytes.TrimSpace(body)
	return len(trimmed) > 0 && trimmed[0] == '{' && bytes.Contains(trimmed, []byte(`"ok"`))
}

// methodFromRequest extracts the Bot API method (last URL path segment).
// The segment never contains the bot token, so it is safe to expose in errors.
func methodFromRequest(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	path := strings.TrimRight(req.URL.Path, "/")
	if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
		path = path[idx+1:]
	}
	// Защита от прострела токена в нестандартных путях: сегмент вида
	// "bot<token>" не показываем.
	if strings.HasPrefix(path, "bot") && strings.ContainsAny(path, ":") {
		return ""
	}
	return path
}

// Copy sends req with ctx attached, streams the response body into w, and closes it.
func (c *Client) Copy(ctx context.Context, req *http.Request, w io.Writer) error {
	if w == nil {
		return stderrors.New("writer is required")
	}

	resp, err := c.open(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, resp.Body)
		return &StatusError{StatusCode: resp.StatusCode, Method: methodFromRequest(req)}
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("copy telegram response body: %w", err)
	}

	return nil
}

func (c *Client) open(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, stderrors.New("context is required")
	}
	if req == nil {
		return nil, stderrors.New("request is required")
	}

	resp, err := c.doer.Do(req.WithContext(ctx))
	if err != nil {
		return nil, &RequestError{err: sanitizeTransportError(err)}
	}
	if resp == nil {
		return nil, stderrors.New("telegram request returned nil response")
	}
	if resp.Body == nil {
		return nil, stderrors.New("telegram request returned nil response body")
	}
	return resp, nil
}

func sanitizeTransportError(err error) error {
	var urlErr *url.Error
	if stderrors.As(err, &urlErr) {
		redacted := *urlErr
		redacted.URL = "[redacted]"
		if urlErr.Err != nil {
			redacted.Err = sanitizeTransportError(urlErr.Err)
		}
		return &redacted
	}

	return err
}
