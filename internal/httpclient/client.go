// Package httpclient contains the internal HTTP transport used by higher-level Bot API clients.
package httpclient

import (
	"context"
	stderrors "errors"
	"net/http"
)

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
	Err error
}

// Error returns a redacted transport error message.
func (e *RequestError) Error() string {
	return "telegram request failed"
}

// Unwrap returns the underlying transport error.
func (e *RequestError) Unwrap() error {
	return e.Err
}

// New creates an internal HTTP client around doer.
//
// If doer is nil, http.DefaultClient is used.
func New(doer Doer) *Client {
	if doer == nil {
		doer = http.DefaultClient
	}

	return &Client{doer: doer}
}

// Do sends req with ctx attached.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, stderrors.New("context is required")
	}
	if req == nil {
		return nil, stderrors.New("request is required")
	}

	resp, err := c.doer.Do(req.WithContext(ctx))
	if err != nil {
		return nil, &RequestError{Err: err}
	}

	return resp, nil
}
