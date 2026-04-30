package exampleutil

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"ai-gram/middleware"
	"ai-gram/telegram"
)

// ParseInt64ListEnv parses a comma-separated int64 environment variable.
func ParseInt64ListEnv(name string) ([]int64, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	values := make([]int64, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s contains invalid int64 value", name)
		}
		values = append(values, id)
	}

	return values, nil
}

// AccessConfigFromEnv builds access-control config for examples from env.
func AccessConfigFromEnv() (middleware.AccessConfig, error) {
	mode := middleware.AccessMode(strings.TrimSpace(OptionalEnv("AIGRAM_ACCESS_MODE", string(middleware.AccessModeAdmin))))
	switch mode {
	case middleware.AccessModeAdmin, middleware.AccessModePublic, middleware.AccessModeOff:
	default:
		return middleware.AccessConfig{}, fmt.Errorf("AIGRAM_ACCESS_MODE must be one of %q, %q, %q", middleware.AccessModeAdmin, middleware.AccessModePublic, middleware.AccessModeOff)
	}

	adminUserIDs, err := ParseInt64ListEnv("AIGRAM_ADMIN_USER_IDS")
	if err != nil {
		return middleware.AccessConfig{}, err
	}
	allowedUserIDs, err := ParseInt64ListEnv("AIGRAM_ALLOWED_USER_IDS")
	if err != nil {
		return middleware.AccessConfig{}, err
	}
	allowedChatIDs, err := ParseInt64ListEnv("AIGRAM_ALLOWED_CHAT_IDS")
	if err != nil {
		return middleware.AccessConfig{}, err
	}

	if len(adminUserIDs) == 0 {
		if fallback, ok := numericEnvID("AIGRAM_CHAT_ID"); ok {
			adminUserIDs = appendUniqueInt64(adminUserIDs, fallback)
			allowedChatIDs = appendUniqueInt64(allowedChatIDs, fallback)
		}
	}

	return middleware.AccessConfig{
		Mode:           mode,
		AdminUserIDs:   adminUserIDs,
		AllowedUserIDs: allowedUserIDs,
		AllowedChatIDs: allowedChatIDs,
	}, nil
}

// AccessController is a runtime-mutable access policy for examples.
type AccessController struct {
	mu             sync.RWMutex
	mode           middleware.AccessMode
	adminUserIDs   []int64
	allowedUserIDs []int64
	allowedChatIDs []int64
}

// NewAccessController creates a controller from static access config.
func NewAccessController(config middleware.AccessConfig) *AccessController {
	mode := config.Mode
	if mode == "" {
		mode = middleware.AccessModeAdmin
	}
	return &AccessController{
		mode:           mode,
		adminUserIDs:   cloneInt64s(config.AdminUserIDs),
		allowedUserIDs: cloneInt64s(config.AllowedUserIDs),
		allowedChatIDs: cloneInt64s(config.AllowedChatIDs),
	}
}

// IsAllowed reports whether update is allowed by the current runtime mode.
func (c *AccessController) IsAllowed(ctx context.Context, update telegram.Update) bool {
	if c == nil {
		return false
	}
	mode, adminUserIDs, allowedUserIDs, allowedChatIDs := c.snapshot()
	switch mode {
	case middleware.AccessModePublic, middleware.AccessModeOff:
		return true
	case middleware.AccessModeAdmin, "":
		return updateAllowed(update, adminUserIDs, allowedUserIDs, allowedChatIDs)
	default:
		return false
	}
}

// IsAdmin reports whether update comes from an admin user.
func (c *AccessController) IsAdmin(update telegram.Update) bool {
	if c == nil {
		return false
	}
	c.mu.RLock()
	adminUserIDs := cloneInt64s(c.adminUserIDs)
	c.mu.RUnlock()

	user := update.EffectiveUser()
	return user != nil && containsInt64(adminUserIDs, user.ID)
}

// Mode returns the current runtime access mode.
func (c *AccessController) Mode() middleware.AccessMode {
	if c == nil {
		return middleware.AccessModeAdmin
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode
}

// SetMode changes the current runtime access mode.
func (c *AccessController) SetMode(mode middleware.AccessMode) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mode = mode
}

func (c *AccessController) snapshot() (middleware.AccessMode, []int64, []int64, []int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode, cloneInt64s(c.adminUserIDs), cloneInt64s(c.allowedUserIDs), cloneInt64s(c.allowedChatIDs)
}

func updateAllowed(update telegram.Update, adminUserIDs []int64, allowedUserIDs []int64, allowedChatIDs []int64) bool {
	if user := update.EffectiveUser(); user != nil {
		if containsInt64(adminUserIDs, user.ID) || containsInt64(allowedUserIDs, user.ID) {
			return true
		}
	}
	if chat := update.EffectiveChat(); chat != nil && containsInt64(allowedChatIDs, chat.ID) {
		return true
	}
	return false
}

func numericEnvID(name string) (int64, bool) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	return id, err == nil
}

func appendUniqueInt64(values []int64, value int64) []int64 {
	if containsInt64(values, value) {
		return values
	}
	return append(values, value)
}

func containsInt64(values []int64, target int64) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func cloneInt64s(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]int64, len(values))
	copy(cloned, values)
	return cloned
}
