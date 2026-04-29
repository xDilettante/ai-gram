// Package httpclient contains the internal HTTP transport used by higher-level Bot API clients.
package httpclient

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// StatusError reports a non-successful HTTP response from the remote server.
type StatusError struct {
	StatusCode int
}

// Error returns a redacted HTTP status error message.
func (e *StatusError) Error() string {
	if e == nil {
		return "telegram HTTP request failed"
	}

	return fmt.Sprintf("telegram HTTP request failed with status %d", e.StatusCode)
}

// New creates an internal HTTP client around doer.
//
// If doer is nil, a new http.Client with a bounded timeout is used.
func New(doer Doer) *Client {
	if doer == nil {
		doer = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{doer: doer}
}

// Do sends req with ctx attached, reads the full response body, and closes it.
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

	return body, nil
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
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, &StatusError{StatusCode: resp.StatusCode}
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
