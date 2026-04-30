// Package middleware contains reusable dispatch middleware helpers.
//
// The package is transport-agnostic and does not send Bot API requests by itself.
// It currently includes panic recovery, timeout contexts, observability hooks,
// and access-control middleware for admin/public/off update handling.
package middleware
