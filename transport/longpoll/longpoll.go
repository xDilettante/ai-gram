// Package longpoll contains the future long polling update transport.
package longpoll

// Config describes future long polling options.
type Config struct {
	// Limit is the maximum number of updates to request per polling call.
	Limit int
	// TimeoutSeconds is the long polling timeout sent to Telegram.
	TimeoutSeconds int
	// AllowedUpdates limits update types requested from Telegram.
	AllowedUpdates []string
}

// Poller is a placeholder for a future long polling update source.
//
// The polling loop will accept context.Context when network I/O is implemented.
type Poller struct {
	config Config
}

// New creates a Poller scaffold with config.
func New(config Config) *Poller {
	return &Poller{config: config}
}

// Config returns the Poller configuration.
func (p *Poller) Config() Config {
	if p == nil {
		return Config{}
	}

	return p.config
}
