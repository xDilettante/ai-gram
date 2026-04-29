// Package webhook contains the future webhook update transport.
package webhook

// Config describes future webhook receiver options.
type Config struct {
	// Path is the HTTP path that will receive Telegram webhook updates.
	Path string
	// SecretToken is the Telegram webhook secret token used to validate incoming requests.
	SecretToken string
	// MaxBodyBytes limits the accepted request body size when webhook I/O is implemented.
	MaxBodyBytes int64
}

// Receiver is a placeholder for a future webhook update receiver.
//
// Request handling and shutdown operations will accept context.Context when implemented.
type Receiver struct {
	config Config
}

// New creates a Receiver scaffold with config.
func New(config Config) *Receiver {
	return &Receiver{config: config}
}

// Config returns the Receiver configuration.
func (r *Receiver) Config() Config {
	if r == nil {
		return Config{}
	}

	return r.config
}
